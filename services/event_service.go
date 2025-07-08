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

	switch event.Type {
	case "ip":
		ep.processIPEvent(event)
	case "email":
		ep.processEmailEvent(event)
	case "user_agent":
		ep.processUserAgentEvent(event)
	case "country":
		ep.processCountryEvent(event)
	default:
		log.Printf("Unknown event type: %s", event.Type)
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

// Stop gracefully stops the event processor
func (ep *EventProcessor) Stop() {
	ep.cancel()
	ep.wg.Wait()
	log.Println("Event processor stopped")
}
