// Copyright 2016 LINE Corporation
//
// LINE Corporation licenses this file to you under the Apache License,
// version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at:
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// +build appengine

package main

import (
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
)

func init() {
	handler, err := httphandler.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	handler.HandleMessage(func(ctx context.Context, bot *linebot.Client, event *linebot.Event) {
		switch message := event.Message.(type) {
		case *linebot.TextMessage:
			if event.Source.Type == linebot.EventSourceTypeUser {
				_, err := bot.PushMessage(event.Source.UserID, linebot.NewTextMessage(message.Text)).Do()
				if err != nil {
					log.Print(err)
				}
			}
		}
	})

	http.Handle("/callback", handler)
}
