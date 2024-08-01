package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_CHANNEL_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/callback", callbackHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

func callbackHandler(c echo.Context) error {
	events, err := bot.ParseRequest(c.Request())
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			return c.String(http.StatusBadRequest, "Invalid Signature")
		}
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch event.Message.(type) {
			case *linebot.TextMessage:
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Hello, world")).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
	return c.String(http.StatusOK, "OK")
}
