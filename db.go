package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	//"github.com/k0kubun/pp"
	//"strconv"
)

const (
	dbFile = "db.json"
)

// Users struct which contains
// an array of users
type Users struct {
	Users []User `json:"users"`
}

// User struct which contains a name
// a type and a list of social links
type User struct {
	Name   string `json:"name"`
	Social Social `json:"social"`
	Date Date `json:"date"`
	Repost Repost `json:"repost"`
}

// Social struct which contains a
// list of links
type Social struct {
	Instagram               string `json:"instagram"`
	InstagramID             int    `json:"instagram_id"`
	VkPageID                int    `json:"vk_page_id"`
	VkPublicID              int    `json:"vk_public_id"`
}

// Date struct which contains a
// list of timestamps of the latest posts
type Date struct {
	InstagramPost  int    `json:"instagram_post"`
	InstagramStory int    `json:"instagram_story"`
	VkPage         int    `json:"vk_page"`
	VkPublic       int    `json:"vk_public"`
}

type Repost struct {
	InstagramPost  bool    `json:"instagram_post"`
	InstagramStory bool    `json:"instagram_story"`
	VkPage         bool    `json:"vk_page"`
	VkPublic       bool    `json:"vk_public"`
}

type Payload struct {
	Timestamp               int64    `json:"timestamp"`
	InstagramPostTimestamp  int64    `json:"instagram_post_timestamp"`
	InstagramStoryTimestamp int64    `json:"instagram_story_timestamp"`
	VkPageTimestamp         int64    `json:"vk_page_timestamp"`
	VkPublicTimestamp       int64    `json:"vk_public_timestamp"`
	VkStatusTimestamp       int64    `json:"vk_status_timestamp"`
	Person                  string   `json:"person"`             //
	InstagramUsername       string   `json:"instagram_username"` //
	InstagramID             int64    `json:"instagram_id"`       //
	Type                    string   `json:"type"`
	From                    string   `json:"from"`
	Source                  string   `json:"source"`
	TelegramChanID          int64    `json:"telegram_chan_id"` //
	RepostMakabaEnabled     bool     `json:"repost_makaba_enabled"`
	RepostTelegramEnabled   bool     `json:"repost_telegram_enabled"`
	RepostVkStatusEnabled   bool     `json:"repost_vk_status_enabled"`
	RepostVkPageEnabled     bool     `json:"repost_vk_page_enabled"`
	RepostVkPublicEnabled   bool     `json:"repost_vk_public_enabled"`
	RepostTelegramChanID    int64    `json:"repost_telegram_chan_id"` //
	VkPageID                int64    `json:"vk_page_id"`              //
	VkPublicID              int64    `json:"vk_public_id"`            //
	DvachBoard              string   `json:"2ch_board"`
	Files                   []string `json:"files"`
	Caption                 string   `json:"caption"`
}

// we initialize our Users array
var users Users


// LoadDBJSON is trying to load our DB into structs from json file
func LoadDBJSON() {
	log.Println("Getting list of users.")
	jsonFile, err := os.Open(dbFile)
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Printf("%v file not found. It looks like we are running first time.\nTrying to create empty DB...", dbFile)
		panic(err)
	}
	log.Printf("%v 📂 opened.", dbFile)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined globally
	err = json.Unmarshal(byteValue, &users)
	if err != nil {
		panic(err)
	}
}

// SaveDBJSON is saving structs back to json file
func SaveDBJSON() {
	byteValue, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(dbFile, byteValue, 0644)
	if err != nil {
		panic(err)
	}
	log.Printf("%v 📁 saved.", dbFile)
}
