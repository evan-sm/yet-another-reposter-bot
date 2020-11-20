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
}

// Social struct which contains a
// list of links
type Social struct {
	Instagram               string `json:"instagram"`
	InstagramID             int    `json:"instagram_id"`
	InstagramPostTimestamp  int    `json:"instagram_post_timestamp"`
	InstagramStoryTimestamp int    `json:"instagram_story_timestamp"`
	VkPageID                int    `json:"vk_page_id"`
	VkPageTimestamp         int    `json:"vk_page_timestamp"`
	VkPublicID              int    `json:"vk_public_id"`
	VkPublicTimestamp       int    `json:"vk_public_timestamp"`
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
