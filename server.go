package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nxadm/tail"
)

var sse = NewSSE[*tail.Line]()

func startServer(listenAddress, filepath string) (err error) {
	t, err := tail.TailFile(filepath, tail.Config{Follow: true, ReOpen: true, Location: &tail.SeekInfo{
		Whence: io.SeekStart,
	}})
	if err != nil {
		return err
	}

	go func() {
		for {
			for line := range t.Lines {
				sse.Message <- line
			}
		}
	}()

	var root = echo.New()
	root.HidePort = true
	root.HideBanner = true

	root.Use(middleware.Recover())

	root.GET(
		"/api/sse",
		sseHandler(filepath),
		sse.Middleware(3000),
	)

	root.GET("/*", echo.NotFoundHandler, middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "webapp",
		Index:      "index.html",
		Browse:     false,
		Filesystem: http.FS(webapp),
	}))

	go func() {
		log.Printf("start listen http://%s\n", listenAddress)
		if err := root.Start(listenAddress); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	var quit = make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Stop()
	t.Cleanup()

	if err := root.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	return nil
}
