package main

import (
	"log"
	"strings"
	//"reflect"
	"time"
	//"github.com/k0kubun/pp"
	tb "gopkg.in/tucnak/telebot.v2"
)

var tg *tb.Bot

func sendRepostTG(p Payload) bool {
	log.Println("Sending payload to telegram.")

	var err error
	var album tb.Album

	tg, err = tb.NewBot(tb.Settings{
		Token:  cfg.TelegramToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	// Create InputMedia for SendAlbum method
	for _, v := range p.Files {
		if strings.Contains(v, ".jpg") {
			if len(album) == 0 {
				album = append(album, &tb.Photo{File: tb.FromURL(v), Caption: p.Caption})
			} else {
				album = append(album, &tb.Photo{File: tb.FromURL(v)})
			}
		}
		if strings.Contains(v, ".mp4") {
			if len(album) == 0 {
				album = append(album, &tb.Video{File: tb.FromURL(v), Caption: p.Caption})
			} else {
				album = append(album, &tb.Video{File: tb.FromURL(v)})
			}
		}
	}

	menu := &tb.ReplyMarkup{}
	menu.Inline(
		menu.Row(menu.URL("URL", p.Source)),)

	if len(album) == 0 {
		_, err = tg.Send(tb.ChatID(p.TelegramChanID), p.Caption, menu)
		if err != nil {
			log.Printf("SendAlbum failed: %v", err)
			return false
		}
		return true
	}
	_, err = tg.SendAlbum(tb.ChatID(p.TelegramChanID), album, menu)
	if err != nil {
		log.Printf("SendAlbum failed: %v", err)
		return false
	}
	return true
}
