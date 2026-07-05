package handlers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"medical-webhook/internal/domain/event"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// sseHeartbeatInterval is how often a keepalive comment is written so an idle
// stream can detect a disconnected client (a flush then fails) and release its
// goroutine and Redis subscription, instead of blocking forever on the next
// event that may never come.
const sseHeartbeatInterval = 15 * time.Second

// SSEHandler handles Server-Sent Events connections for real-time streaming
type SSEHandler struct {
	eventBus event.EventBus
}

// NewSSEHandler creates a new SSE handler
func NewSSEHandler(eventBus event.EventBus) *SSEHandler {
	return &SSEHandler{
		eventBus: eventBus,
	}
}

// Stream handles the SSE endpoint
// GET /api/events/stream?types=equipment.updated,ticket.created
func (h *SSEHandler) Stream(c *fiber.Ctx) error {
	// Check if event bus is available (Redis connected)
	if h.eventBus == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"success": false,
			"error":   "Real-time events not available (Redis not connected)",
		})
	}

	// Parse optional event type filter from query params
	var eventTypes []string
	if types := c.Query("types"); types != "" {
		eventTypes = strings.Split(types, ",")
		for i := range eventTypes {
			eventTypes[i] = strings.TrimSpace(eventTypes[i])
		}
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "Cache-Control")

	log.Printf("📡 SSE client connected (filter: %v)", eventTypes)

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		// ⚠️ CRITICAL FIX: Fiber/fasthttp uses object pooling for request contexts.
		// Passing c.Context() directly to a long-running Redis Subscribe goroutine
		// causes a panic when fasthttp recycles the context.
		// Solution: Create a detached standard context for Redis, and cancel it manually
		// when the SSE stream ends or the client disconnects.
		bgCtx, cancel := context.WithCancel(context.Background())
		defer cancel() // Ensure Redis subscription is cleaned up when stream ends

		// Subscribe to events from the event bus using the detached context
		events, err := h.eventBus.Subscribe(bgCtx, eventTypes...)
		if err != nil {
			log.Printf("❌ SSE subscribe error: %v", err)
			fmt.Fprintf(w, "event: error\ndata: {\"error\":\"subscribe failed\"}\n\n")
			w.Flush()
			return
		}

		// Send initial connection confirmation
		fmt.Fprintf(w, "event: connected\ndata: {\"message\":\"Connected to event stream\"}\n\n")
		w.Flush()

		// Heartbeat so an idle client's disconnect is detected (via a failing
		// flush) even when no matching event ever arrives — otherwise the loop
		// blocks on <-events forever, leaking this goroutine and the Redis
		// subscription.
		heartbeat := time.NewTicker(sseHeartbeatInterval)
		defer heartbeat.Stop()

		for {
			select {
			case <-bgCtx.Done():
				// This won't actually trigger unless we explicitly call cancel(),
				// which happens when the handler exits. But we also need to know
				// when the client actually disconnects.
				log.Println("📡 SSE client disconnected (bgCtx)")
				return
			case <-heartbeat.C:
				// SSE comment line (ignored by EventSource clients); a flush
				// error means the client has gone away.
				fmt.Fprintf(w, ": ping\n\n")
				if err := w.Flush(); err != nil {
					log.Printf("📡 SSE heartbeat flush failed (client disconnected): %v", err)
					return
				}
			case evt, ok := <-events:
				if !ok {
					log.Println("📡 Event channel closed")
					return
				}

				data, err := json.Marshal(evt)
				if err != nil {
					log.Printf("⚠️ SSE marshal error: %v", err)
					continue
				}

				// Write SSE format: "event: <type>\ndata: <json>\n\n"
				fmt.Fprintf(w, "event: %s\ndata: %s\n\n", evt.Type, string(data))
				if err := w.Flush(); err != nil {
					log.Printf("📡 SSE flush error (client disconnected): %v", err)
					return
				}
			}
		}
	})

	return nil
}
