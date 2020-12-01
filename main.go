package main

import (
	//"encoding/json"
	//"fmt"
	//"io/ioutil"
	"log"
	//"reflect"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/k0kubun/pp"
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

	for {
		LoadDBJSON()
		pp.Println(users)

		// Iterate through every user within our users array
		for k, _ := range users.Users {
			//log.Printf("key: \"%v\" | value: \"%v\"", k, v)

			// Check socials
			//log.Println(reflect.TypeOf(users.Users[k]))
			checkVK(k, users.Users[k])
			checkIG(k, users.Users[k])
			//checkTT(users.Users[k])

			//log.Printf("%v", post)
			//log.Printf("%v", post.Text)

			time.Sleep(2000)
			//checkIGStories(users.Users[k].Social.InstagramID)
			//checkIGPost()
			// Do some changes to struct
		}
		time.Sleep(5 * time.Minute)
	}

}

func checkTT(u User) {
	log.Printf("Checking TikTok %v\n", u.Social.TikTok)
	videos, err := getLikedVideos(u, 3)
	if err != nil {
		log.Printf("%v", err)
	}
	pp.Println(videos)
}

func sendRepost() {
	log.Println("Preparing for reposting")
}
