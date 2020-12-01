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
	Name    string  `json:"name"`
	Social  Social  `json:"social"`
	Date    Date    `json:"date"`
	Setting Setting `json:"setting"`
	Repost  Repost  `json:"repost"`
}

// Social struct which contains a
// list of links
type Social struct {
	Instagram   string `json:"instagram"`
	InstagramID int    `json:"instagram_id"`
	VkPageID    int    `json:"vk_page_id"`
	VkPublicID  int    `json:"vk_public_id"`
	TikTok      string `json:"tiktok"`
}

// Date struct which contains a
// list of timestamps of the latest posts
type Date struct {
	InstagramPost  int `json:"instagram_post"`
	InstagramStory int `json:"instagram_story"`
	VkPage         int `json:"vk_page"`
	VkPublic       int `json:"vk_public"`
	TikTok         int `json:"tiktok"`
}

type Setting struct {
	InstagramPost  bool `json:"instagram_post"`
	InstagramStory bool `json:"instagram_story"`
	VkPage         bool `json:"vk_page"`
	VkPublic       bool `json:"vk_public"`
	TikTok         bool `json:"tiktok"`
	Makaba         bool `json:"makaba"`
}

type Repost struct {
	TelegramChanID int    `json:"telegram_channel_id"`
	Board          string `json:"board"`
	Thread         string `json:"thread"`
}

type Payload struct {
	// Body
	Person    string
	Timestamp int
	Caption   string
	From      string
	Type      string
	Source    string
	Files     []string

	// Destination
	TelegramChanID int
	Board          string
	Thread         string
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
