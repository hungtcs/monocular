package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Message struct {
	ID    any
	Data  string
	Event string
}

type SSE[T any] struct {
	Message       chan T
	AllClients    map[chan T]bool
	ClientClose   chan chan T
	ClientConnect chan chan T
}

func NewSSE[T any]() (sse *SSE[T]) {
	sse = &SSE[T]{
		Message:       make(chan T),
		AllClients:    make(map[chan T]bool),
		ClientClose:   make(chan chan T),
		ClientConnect: make(chan chan T),
	}

	go sse.start()

	return sse
}

func (sse *SSE[T]) start() {
	for {
		select {

		case client := <-sse.ClientConnect:
			sse.AllClients[client] = true
			log.Printf("client connected clients=%d\n", len(sse.AllClients))

		case client := <-sse.ClientClose:
			delete(sse.AllClients, client)
			close(client)
			log.Printf("client disconnected clients=%d\n", len(sse.AllClients))

		case eventMsg := <-sse.Message:
			for clientMessageChan := range sse.AllClients {
				clientMessageChan <- eventMsg
			}
		}
	}
}

func (sse *SSE[T]) Middleware(retry uint) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Content-Type", "text/event-stream")
			c.Response().Header().Set("Cache-Control", "no-cache")
			c.Response().Header().Set("Connection", "keep-alive")
			c.Response().Header().Set("Transfer-Encoding", "chunked")
			c.Response().WriteHeader(http.StatusOK)

			c.Response().Write([]byte(fmt.Sprintf("retry: %d\n\n", retry)))
			c.Response().Flush()

			clientChan := make(chan T)
			sse.ClientConnect <- clientChan

			defer func() {
				sse.ClientClose <- clientChan
			}()

			c.Set("CLIENT_CHAN", clientChan)

			return next(c)
		}
	}
}

func (sse *SSE[T]) SendMessage(c echo.Context, message Message) {
	var buf = bytes.NewBuffer(nil)

	if message.ID != nil {
		buf.WriteString(fmt.Sprintf("id: %v\n", message.ID))
	}
	buf.WriteString(fmt.Sprintf("event: %s\n", message.Event))
	buf.WriteString(fmt.Sprintf("data: %s\n", message.Data))
	buf.WriteString("\n")

	c.Response().Write(buf.Bytes())
	c.Response().Flush()
}
