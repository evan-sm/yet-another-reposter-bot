package main

import (
	"encoding/json"
	"errors"
	"fmt"
	//"io/ioutil"
	"log"
	//"os"
	"github.com/k0kubun/pp"
	"github.com/parnurzeal/gorequest"
	"github.com/tidwall/gjson"
	"strconv"
)

var IGSession bool
var IGSessionID string

type Items []struct {
	//	Typename               string                `json:"__typename"`
	DisplayURL string `json:"display_url"`
	IsVideo    bool   `json:"is_video"`
	VideoURL   string `json:"video_url,omitempty"`
}

var items Items

func checkIG(k int, u User) {
	// Switch session ID to avoid rate-limits
	if IGSession {
		IGSessionID = cfg.IGSessionID
	} else {
		IGSessionID = cfg.IGSessionID2
	}
	log.Printf("Checking Instagram\nIGSessionID is: %v\n", IGSessionID)

	if post, new := checkNewIGPost(u); new {
		log.Printf("New post found: %v\n", post)
		payload := preparePayloadFromIG(u, post)
		pp.Println(payload)

		// Send to Telegram
		if ok := sendRepostTG(payload); ok {
			users.Users[k].Date.InstagramPost = payload.Timestamp
			SaveDBJSON()
		}
		// Send to Makaba
		if ok := sendRepostMakaba(u, payload); ok {
			users.Users[k].Date.InstagramPost = payload.Timestamp
			SaveDBJSON()
		}
	}

	if items, new := checkNewIGStories(u); new {
		log.Printf("New stories found: %v\n", items)
		payload := preparePayloadFromIGStories(u, items)
		pp.Println(payload)

		// Send to Telegram
		if ok := sendRepostTG(payload); ok {
			users.Users[k].Date.InstagramStory = payload.Timestamp
			SaveDBJSON()
		}
		// Send to Makaba
		if ok := sendRepostMakaba(u, payload); ok {
			users.Users[k].Date.InstagramStory = payload.Timestamp
			SaveDBJSON()
		}
	}
	//checkNewIGStories(u)
}

func checkNewIGPost(u User) (string, bool) {
	// Skip if repost disabled
	if !u.Setting.InstagramPost {
		return "", false
	}
	log.Printf("Checking IG Post()\n")

	// Get all posts from instagram profile
	js, err := extractJSONFromProfile(u)
	if err != nil {
		log.Printf("%v", err)
		return "", false
	}

	// Drop all posts except first one
	post, err := getLastIGPost(js, u)
	if err != nil {
		log.Printf("%v", err) // All posts are old
		return "", false
	}

	//log.Printf("\n\n\nJSON: %v", post)

	// Ok, we got new posts
	return post, true
}

func checkNewIGStories(u User) (string, bool) {
	// Skip if repost disabled
	if !u.Setting.InstagramStory {
		return "", false
	}
	log.Printf("Checking IG stories()\n")

	// Get all stories from instagram profile
	js, err := extractJSONFromStories(u)
	if err != nil {
		log.Printf("%v", err)
		return "", false
	}

	//log.Printf("body: %v", js[0:100])

	// Drop all old stories
	items, err := getNewIGStories(js, u)
	if err != nil {
		log.Printf("%v", err) // All posts are old
		return "", false
	}

	if items == "" {
		return "", false
	}
	//log.Printf("\n\n\nJSON: %v", post)

	// Ok, we got new posts
	return items, true
}

// extractJSONFromProfile retrieves json from profile. Example: instagram.com/wmw/?__a=1
func extractJSONFromProfile(u User) (string, error) {
	url := fmt.Sprintf(`https://www.instagram.com/%v/?__a=1`, u.Social.Instagram)
	log.Printf("url: %v\n", url)
	req := gorequest.New()
	resp, body, errs := req.Get(url).
		Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Safari/537.36").
		Set("cookie", IGSessionID).End()
	//Retry(4, 1200*time.Second, http.StatusBadRequest, http.StatusInternalServerError, http.StatusTooManyRequests).End()
	log.Printf("resp: %v", resp.Status)
	if errs != nil {
		log.Fatalf("%v\n%v", errs, resp)
	}
	return body, nil
}

// extractJSONFromStories retrieves json from stories. Example: instagram.com/stories/wmw/
func extractJSONFromStories(u User) (string, error) {
	url := fmt.Sprintf(`https://i.instagram.com/api/v1/feed/user/%v/story/`, u.Social.InstagramID)
	log.Printf("url: %v\n", url)
	req := gorequest.New()
	resp, body, errs := req.Get(url).
		Set("user-agent", "Instagram 10.26.0 (iPhone7,2; iOS 10_1_1; en_US; en-US; scale=2.00; gamut=normal; 750x1334) AppleWebKit/420+").
		Set("cookie", IGSessionID).End()
	//Retry(4, 1200*time.Second, http.StatusBadRequest, http.StatusInternalServerError, http.StatusTooManyRequests).End()
	log.Printf("resp: %v\n", resp.Status)
	if errs != nil {
		log.Fatalf("%v\n%v", errs, resp)
	}
	return body, nil
}

func getLastIGPost(js string, u User) (string, error) {
	path := fmt.Sprintf(`graphql.user.edge_owner_to_timeline_media.edges.#(node.taken_at_timestamp>%v)`, u.Date.InstagramPost)
	post := gjson.Get(js, path).String()
	if post == "" {
		return "", errors.New("IG profile posts are all old.")
	}
	return post, nil
}

func getNewIGStories(js string, u User) (string, error) {
	path := fmt.Sprintf(`reel.items.#(taken_at>%v)#`, u.Date.InstagramStory)
	items := gjson.Get(js, path).String()
	if items == "[]" {
		return "", errors.New("All IG stories are all old.")
	}
	return items, nil
}

func preparePayloadFromIG(u User, post string) Payload {
	var files []string
	p := Payload{}

	switch gjson.Get(post, "node.__typename").String() {
	case "GraphSidecar":
		log.Printf("This post got multiple items")

		p.Timestamp = getTimestampIGPost(post)
		p.Caption = getCaptionIGPost(post)
		p.Source = getSourceIGPost(post)

		// Convert JSON to Go structs
		path := "node.edge_sidecar_to_children.edges.#.node"
		js := gjson.Get(post, path).String()
		if err := json.Unmarshal([]byte(js), &items); err != nil {
			panic(err)
		}

		// Get all items (.jpg, .mp4) from post
		for _, v := range items {
			if v.IsVideo {
				files = append(files, v.VideoURL) // add .mp4
			} else {
				files = append(files, v.DisplayURL) // add .jpg
			}
		}
		pp.Println(files)
	case "GraphVideo":
		log.Printf("This post is a video")

		p.Timestamp = getTimestampIGPost(post)
		p.Caption = getCaptionIGPost(post)
		p.Source = getSourceIGPost(post)

		url := gjson.Get(post, "node.video_url").String()
		files = append(files, url)

		pp.Println(files)
	case "GraphImage":
		log.Printf("This post is a photo post.")

		p.Timestamp = getTimestampIGPost(post)
		p.Caption = getCaptionIGPost(post)
		p.Source = getSourceIGPost(post)
		url := gjson.Get(post, "node.display_url").String()
		files = append(files, url)
	}

	p.Person = u.Name
	p.From = "instagram"
	p.Type = "post"
	p.TelegramChanID = u.Repost.TelegramChanID
	p.Board = u.Repost.Board
	p.Thread = u.Repost.Thread
	p.Files = files
	return p
}

func getTimestampIGPost(post string) int {
	time64 := gjson.Get(post, "node.taken_at_timestamp").String()
	time, err := strconv.Atoi(time64)
	if err != nil {
		log.Printf("%v", err)
	}
	return time
}

func getTimestampIGStory(item string) int {
	time64 := gjson.Get(item, "taken_at").String()
	time, err := strconv.Atoi(time64)
	if err != nil {
		log.Printf("%v", err)
	}
	return time
}

func getCaptionIGPost(post string) string {
	return gjson.Get(post, "node.edge_media_to_caption.edges.0.node.text").String()
}

func getSourceIGPost(post string) string {
	source := gjson.Get(post, "node.shortcode").String()
	return fmt.Sprintf("https://www.instagram.com/p/%v/", source)
}

func preparePayloadFromIGStories(u User, items string) Payload {
	var files []string
	var count int
	p := Payload{}

	result := gjson.Get(items, "@valid")
	result.ForEach(func(key, value gjson.Result) bool {
		log.Printf("\n\n%v", value.String()[0:20])

		switch gjson.Get(value.String(), "media_type").String() {
		case "2": // video
			if count == 4 {
				return false
			}
			count = count + 1
			log.Printf("count: %v", count)
			url := gjson.Get(value.String(), "video_versions.0.url").String()
			log.Printf("Got video story: %v", url)
			files = append(files, url) // add .mp4
			p.Timestamp = getTimestampIGStory(value.String())
		case "1": // image
			if count == 4 {
				return false
			}
			count = count + 1
			log.Printf("count: %v", count)
			url := gjson.Get(value.String(), "image_versions2.candidates.0.url").String()
			log.Printf("Got image story: %v", url)
			files = append(files, url) // add .jpg
			p.Timestamp = getTimestampIGStory(value.String())
		}
		return true // keep iterating
	})

	p.Person = u.Name
	p.From = "instagram"
	p.Type = "story"
	p.TelegramChanID = u.Repost.TelegramChanID
	p.Board = u.Repost.Board
	p.Thread = u.Repost.Thread
	p.Files = files

	return p
}
