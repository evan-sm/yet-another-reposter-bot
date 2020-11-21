package main

import (
	//"encoding/json"
	//"fmt"
	//"io/ioutil"
	"log"
	"reflect"
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
	LoadDBJSON()

	pp.Println(users)

	// Iterate through every user within our users array
	for k, v := range users.Users {
		log.Printf("key: \"%v\" | value: \"%v\"", k, v)

		// Check socials
		log.Println(reflect.TypeOf(users.Users[k]))
		checkVK(users.Users[k])
		checkIG()
		checkTT()

		//log.Printf("%v", post)
		//log.Printf("%v", post.Text)


		break
		time.Sleep(2000)
		//checkIGStories(users.Users[k].Social.InstagramID)
		//checkIGPost()
		// Do some changes to struct
		users.Users[k].Social.InstagramID = 123
	}

	//SaveDBJSON()
}

func checkIG() {
	log.Println("Checking Instagram")
}


func checkTT() {
	log.Println("Checking TikTok")
}

func sendRepost() {
	log.Println("Preparing for reposting")
}