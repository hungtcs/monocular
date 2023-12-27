package main

import (
	"github.com/labstack/echo/v4"
	"github.com/nxadm/tail"
)

func sseHandler(filepath string) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		var clientChan = c.Get("CLIENT_CHAN").(chan *tail.Line)
		for {
			select {
			case line := <-clientChan:
				sse.SendMessage(c, Message{
					ID:    line.Num,
					Data:  line.Text,
					Event: "log",
				})
			case <-c.Request().Context().Done():
				return nil
			}
		}
	}
}
