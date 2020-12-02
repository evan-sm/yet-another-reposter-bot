package main

import (
	//"encoding/json"
	//"fmt"
	//"io/ioutil"
	//"github.com/SevereCloud/vksdk/v2/api/params"
	"errors"
	"fmt"
	"log"
	"reflect"
	//"os"
	//"strconv"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/k0kubun/pp"
)

func checkVK(k int, u User) {
	if post, new := checkNewPostVKPublic(u); new {
		// Prepare payload
		payload := preparePayloadFromVK(u, post)
		pp.Println(payload)

		// Send to Telegram
		if ok := sendRepostTG(u, payload); ok {
			users.Users[k].Date.VkPublic = payload.Timestamp
			SaveDBJSON()
		}
		// Send to Makaba
		if ok := sendRepostMakaba(u, payload); ok {
			users.Users[k].Date.VkPublic = payload.Timestamp
			SaveDBJSON()
		}
	}

	if post, new := checkNewPostVKPage(u); new {
		// Prepare payload
		payload := preparePayloadFromVK(u, post)
		pp.Println(payload)

		// Send to Telegram
		if ok := sendRepostTG(u, payload); ok {
			users.Users[k].Date.VkPage = payload.Timestamp
			SaveDBJSON()
		}

		// Send to Makaba
		if ok := sendRepostMakaba(u, payload); ok {
			users.Users[k].Date.VkPublic = payload.Timestamp
			SaveDBJSON()
		}
	}
}

func checkNewPostVKPublic(u User) (object.WallWallpost, bool) {
	// Skip if repost disabled
	if !u.Setting.VkPublic {
		return object.WallWallpost{}, false
	}
	log.Printf("Checking VK public. https://vk.com/wall%v\n", u.Social.VkPublicID)

	post, err := getLastVKPost(u.Social.VkPublicID)
	if err != nil {
		log.Printf("%v", err)
		return post, false
	}

	if reflect.ValueOf(post).IsZero() {
		log.Println("empty")
		return post, false
	}

	if post.Date > u.Date.VkPublic {
		log.Printf("New VK public post found! \"%v\"", post.Text)
		return post, true
	}
	return post, false
}

func checkNewPostVKPage(u User) (object.WallWallpost, bool) {
	// Skip if repost disabled
	if !u.Setting.VkPage {
		return object.WallWallpost{}, false
	}
	log.Printf("Checking VK page. https://vk.com/id%v\n", u.Social.VkPageID)

	post, err := getLastVKPost(u.Social.VkPageID)
	if err != nil {
		log.Printf("%v", err)
		return post, false
	}

	if reflect.ValueOf(post).IsZero() {
		log.Println("empty")
		return post, false
	}

	if post.Date > u.Date.VkPage {
		log.Printf("New VK page post found! \"%v\"", post.Text)
		return post, true
	}
	return post, false
}

func checkVKPage() {
	log.Println("Checking VK page.")
}

func getLastVKPost(id int) (object.WallWallpost, error) {
	var date int
	var i int
	var post object.WallWallpost

	vk := api.NewVK(cfg.VKToken)

	wall, err := vk.WallGet(api.Params{
		"owner_id": id,
		"count":    2,
		"filter":   "owner",
	})
	if err != nil {
		log.Fatal(err)
	}

	if len(wall.Items) == 0 {
		return post, errors.New("Wall have no posts.")
	}

	// Return post if it's the only one
	if len(wall.Items) == 1 {
		return wall.Items[0], nil
	}

	// Drop
	for k, v := range wall.Items {

		// Ignore reposts
		if len(v.CopyHistory) != 0 {
			continue
		}

		if v.Date > date {
			date = v.Date
			i = k
		}

		pp.Println(k, v.Text, len(v.CopyHistory), len(v.Attachments), len(v.Attachments), v.Date)
	}
	//log.Printf("date: %d; index: %d", date, i)
	//log.Println(reflect.TypeOf(wall.Items[i]))
	return wall.Items[i], nil
}

func checkVKStatus(id int) {
	log.Printf("Checking VK status %v\n", id)

}

func preparePayloadFromVK(u User, post object.WallWallpost) Payload {
	var files []string

	log.Printf("Preparing payload from VK post.")

	for _, v := range post.Attachments {
		if v.Type == "photo" {
			url := getVKOriginalSize(v.Photo)
			files = append(files, url) // add .jpg url to slice
		}
	}
	//pp.Printf("%v\n", files)
	p := Payload{}
	p.Person = u.Name
	p.Timestamp = post.Date
	p.From = "vk"
	p.Caption = post.Text
	p.Type = "post"
	p.TelegramChanID = u.Repost.TelegramChanID
	p.Board = u.Repost.Board
	p.Thread = u.Repost.Thread
	p.Source = fmt.Sprintf("https://vk.com/wall%v_%v", post.OwnerID, post.ID)
	p.Files = files
	return p
}

func getVKOriginalSize(p object.PhotosPhoto) string {
	var width float64
	var url string
	for _, v := range p.Sizes {
		if v.Width > width {
			width = v.Width
			url = v.URL
		}
	}
	log.Printf("Original .jpg 🖼 %vp @ %v\n", width, url)
	return url
}
