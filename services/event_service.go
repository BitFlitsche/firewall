package services

import (
	"context"
	"firewall/models"
	"log"
	"sync"
	"time"
)

// Event represents a data change event
type Event struct {
	Type      string      `json:"type"`
	Action    string      `json:"action"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// EventProcessor handles event processing
type EventProcessor struct {
	events chan Event
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

var (
	eventProcessor *EventProcessor
	once           sync.Once
)

// GetEventProcessor returns the singleton event processor
func GetEventProcessor() *EventProcessor {
	once.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		eventProcessor = &EventProcessor{
			events: make(chan Event, 1000), // Buffer for 1000 events
			ctx:    ctx,
			cancel: cancel,
		}
		eventProcessor.start()
	})
	return eventProcessor
}

// PublishEvent publishes an event to the event stream
func PublishEvent(eventType, action string, data interface{}) {
	processor := GetEventProcessor()
	event := Event{
		Type:      eventType,
		Action:    action,
		Data:      data,
		Timestamp: time.Now(),
	}

	select {
	case processor.events <- event:
		log.Printf("Event published: %s.%s", eventType, action)
	default:
		log.Printf("Warning: Event buffer full, dropping event: %s.%s", eventType, action)
	}
}

// start begins processing events
func (ep *EventProcessor) start() {
	ep.wg.Add(1)
	go func() {
		defer ep.wg.Done()
		for {
			select {
			case event := <-ep.events:
				ep.processEvent(event)
			case <-ep.ctx.Done():
				return
			}
		}
	}()
}

// processEvent handles individual events
func (ep *EventProcessor) processEvent(event Event) {
	log.Printf("Processing event: %s.%s", event.Type, event.Action)

	// Get cache instance for invalidation
	cache := GetCacheFactory()

	switch event.Type {
	case "ip":
		ep.processIPEvent(event)
		// Invalidate cache for IP-related data
		if event.Action == "created" || event.Action == "updated" || event.Action == "deleted" {
			cache.InvalidateAll("ip")
			cache.InvalidateFilter("ip")
		}
	case "email":
		ep.processEmailEvent(event)
		// Invalidate cache for email-related data
		if event.Action == "created" || event.Action == "updated" || event.Action == "deleted" {
			cache.InvalidateAll("email")
			cache.InvalidateFilter("email")
		}
	case "user_agent":
		ep.processUserAgentEvent(event)
		// Invalidate cache for user agent-related data
		if event.Action == "created" || event.Action == "updated" || event.Action == "deleted" {
			cache.InvalidateAll("user_agent")
			cache.InvalidateFilter("user_agent")
		}
	case "country":
		ep.processCountryEvent(event)
		// Invalidate cache for country-related data
		if event.Action == "created" || event.Action == "updated" || event.Action == "deleted" {
			cache.InvalidateAll("country")
			cache.InvalidateFilter("country")
		}
	case "charset":
		ep.processCharsetEvent(event)
		// Invalidate cache for charset-related data
		if event.Action == "created" || event.Action == "updated" || event.Action == "deleted" {
			cache.InvalidateAll("charset")
		}
	case "username":
		ep.processUsernameEvent(event)
		// Invalidate cache for username-related data
		if event.Action == "created" || event.Action == "updated" || event.Action == "deleted" {
			cache.InvalidateAll("username")
			cache.InvalidateFilter("username")
		}
	default:
		log.Printf("Unknown event type: %s", event.Type)
	}
}

// Charset-Event-Handler
func (ep *EventProcessor) processCharsetEvent(event Event) {
	switch event.Action {
	case "created", "updated":
		if charsetData, ok := event.Data.(models.CharsetRule); ok {
			if err := SyncCharsetToES(charsetData); err != nil {
				log.Printf("Error indexing charset: %v", err)
			}
		}
	case "deleted":
		if charsetData, ok := event.Data.(models.CharsetRule); ok {
			if err := DeleteCharsetFromES(charsetData.ID); err != nil {
				log.Printf("Error deleting charset from ES: %v", err)
			}
		}
	}
}

// processIPEvent handles IP-related events
func (ep *EventProcessor) processIPEvent(event Event) {
	switch event.Action {
	case "created", "updated":
		if ipData, ok := event.Data.(models.IP); ok {
			if err := IndexIPAddress(ipData); err != nil {
				log.Printf("Error indexing IP: %v", err)
				// Queue for retry
				QueueForRetry("sync_ip", ipData)
			}
		}
	case "deleted":
		// Handle deletion if needed
		log.Printf("IP deletion event received")
	}
}

// processEmailEvent handles email-related events
func (ep *EventProcessor) processEmailEvent(event Event) {
	switch event.Action {
	case "created", "updated":
		if emailData, ok := event.Data.(models.Email); ok {
			if err := IndexEmail(emailData); err != nil {
				log.Printf("Error indexing email: %v", err)
				QueueForRetry("sync_email", emailData)
			}
		}
	case "deleted":
		log.Printf("Email deletion event received")
	}
}

// processUserAgentEvent handles user agent-related events
func (ep *EventProcessor) processUserAgentEvent(event Event) {
	switch event.Action {
	case "created", "updated":
		if userAgentData, ok := event.Data.(models.UserAgent); ok {
			if err := IndexUserAgent(userAgentData); err != nil {
				log.Printf("Error indexing user agent: %v", err)
				QueueForRetry("sync_user_agent", userAgentData)
			}
		}
	case "deleted":
		log.Printf("User agent deletion event received")
	}
}

// processCountryEvent handles country-related events
func (ep *EventProcessor) processCountryEvent(event Event) {
	switch event.Action {
	case "created", "updated":
		if countryData, ok := event.Data.(models.Country); ok {
			if err := IndexCountry(countryData); err != nil {
				log.Printf("Error indexing country: %v", err)
				QueueForRetry("sync_country", countryData)
			}
		}
	case "deleted":
		log.Printf("Country deletion event received")
	}
}

// Username-Event-Handler
func (ep *EventProcessor) processUsernameEvent(event Event) {
	switch event.Action {
	case "created", "updated":
		if usernameData, ok := event.Data.(models.UsernameRule); ok {
			if err := SyncUsernameToES(usernameData); err != nil {
				log.Printf("Error indexing username: %v", err)
			}
		}
	case "deleted":
		if usernameData, ok := event.Data.(models.UsernameRule); ok {
			if err := DeleteUsernameFromES(usernameData.ID); err != nil {
				log.Printf("Error deleting username from ES: %v", err)
			}
		}
	}
}

// Stop gracefully stops the event processor
func (ep *EventProcessor) Stop() {
	ep.cancel()
	ep.wg.Wait()
	log.Println("Event processor stopped")
}
