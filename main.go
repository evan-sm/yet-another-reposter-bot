package main

import (
	//"encoding/json"
	//"fmt"
	//"io/ioutil"
	"github.com/caarlos0/env/v6"
	"github.com/k0kubun/pp"
	"log"
	"time"
	//"os"
	//"strconv"
)

func main() {
	log.Println("Yet Another Reposter Bot started.")
	if err := env.Parse(&cfg); err != nil {
		log.Fatalln("Config", err)
	}
	checkNewPosts()

}

func checkNewPosts() {
	LoadDBJSON()

	pp.Println(users)

	// Iterate through every user within our users array
	for k, v := range users.Users {
		log.Printf("key: \"%v\" | value: \"%v\"", k, v)

		// Check socials
		checkVKWallGet(users.Users[k].Social.VkPublicID)
		break
		time.Sleep(2000)
		//checkIGStories(users.Users[k].Social.InstagramID)
		//checkIGPost()
		// Do some changes to struct
		users.Users[k].Social.InstagramID = 123
	}

	//SaveDBJSON()
}
