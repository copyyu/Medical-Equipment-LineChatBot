package handlers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"medical-webhook/internal/domain/event"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

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

		for {
			select {
			case <-bgCtx.Done():
				// This won't actually trigger unless we explicitly call cancel(),
				// which happens when the handler exits. But we also need to know
				// when the client actually disconnects.
				log.Println("📡 SSE client disconnected (bgCtx)")
				return
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

// StreamWithFastHTTP is an alternative implementation using fasthttp hijack
// for better compatibility with Fiber's streaming
func (h *SSEHandler) StreamWithFastHTTP(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/event-stream")
	ctx.Response.Header.Set("Cache-Control", "no-cache")
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")

	fmt.Fprintf(ctx, "event: connected\ndata: {\"message\":\"Connected\"}\n\n")
	ctx.SetBodyStreamWriter(func(w *bufio.Writer) {
		events, err := h.eventBus.Subscribe(ctx)
		if err != nil {
			return
		}

		for evt := range events {
			data, _ := json.Marshal(evt)
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", evt.Type, string(data))
			w.Flush()
		}
	})
}
