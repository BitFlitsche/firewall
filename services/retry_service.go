package services

import (
	"context"
	"firewall/models"
	"log"
	"sync"
	"time"
)

// RetryItem represents an item to be retried
type RetryItem struct {
	Type      string      `json:"type"`
	Action    string      `json:"action"`
	Data      interface{} `json:"data"`
	Attempts  int         `json:"attempts"`
	NextRetry time.Time   `json:"next_retry"`
}

// RetryQueue handles retrying failed operations
type RetryQueue struct {
	items chan RetryItem
	wg    sync.WaitGroup
	ctx   context.Context
	stop  context.CancelFunc
}

var (
	retryQueue *RetryQueue
	retryOnce  sync.Once
)

// GetRetryQueue returns the singleton retry queue
func GetRetryQueue() *RetryQueue {
	retryOnce.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		retryQueue = &RetryQueue{
			items: make(chan RetryItem, 1000),
			ctx:   ctx,
			stop:  cancel,
		}
		retryQueue.start()
	})
	return retryQueue
}

// QueueForRetry adds an item to the retry queue
func QueueForRetry(retryType string, data interface{}) {
	queue := GetRetryQueue()
	item := RetryItem{
		Type:      retryType,
		Data:      data,
		Attempts:  0,
		NextRetry: time.Now().Add(time.Second), // Retry after 1 second
	}

	select {
	case queue.items <- item:
		log.Printf("Item queued for retry: %s", retryType)
	default:
		log.Printf("Warning: Retry queue full, dropping item: %s", retryType)
	}
}

// start begins processing retry items
func (rq *RetryQueue) start() {
	rq.wg.Add(1)
	go func() {
		defer rq.wg.Done()
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case item := <-rq.items:
				rq.processRetryItem(item)
			case <-ticker.C:
				// Check for items ready to retry
			case <-rq.ctx.Done():
				return
			}
		}
	}()
}

// processRetryItem handles retrying an item
func (rq *RetryQueue) processRetryItem(item RetryItem) {
	if time.Now().Before(item.NextRetry) {
		// Not ready to retry yet, requeue
		select {
		case rq.items <- item:
		default:
			log.Printf("Warning: Could not requeue retry item: %s", item.Type)
		}
		return
	}

	item.Attempts++
	maxAttempts := 5

	if item.Attempts > maxAttempts {
		log.Printf("Max retry attempts reached for %s, giving up", item.Type)
		return
	}

	// Calculate next retry time with exponential backoff
	backoff := time.Duration(item.Attempts*item.Attempts) * time.Second
	item.NextRetry = time.Now().Add(backoff)

	// Attempt the operation
	var err error
	switch item.Type {
	case "sync_ip":
		if ipData, ok := item.Data.(models.IP); ok {
			err = IndexIPAddress(ipData)
		}
	case "sync_email":
		if emailData, ok := item.Data.(models.Email); ok {
			err = IndexEmail(emailData)
		}
	case "sync_user_agent":
		if userAgentData, ok := item.Data.(models.UserAgent); ok {
			err = IndexUserAgent(userAgentData)
		}
	case "sync_country":
		if countryData, ok := item.Data.(models.Country); ok {
			err = IndexCountry(countryData)
		}
	}

	if err != nil {
		log.Printf("Retry attempt %d failed for %s: %v", item.Attempts, item.Type, err)
		// Requeue for next retry
		select {
		case rq.items <- item:
		default:
			log.Printf("Warning: Could not requeue retry item: %s", item.Type)
		}
	} else {
		log.Printf("Retry successful for %s after %d attempts", item.Type, item.Attempts)
	}
}

// Stop gracefully stops the retry queue
func (rq *RetryQueue) Stop() {
	rq.stop()
	rq.wg.Wait()
	log.Println("Retry queue stopped")
}
