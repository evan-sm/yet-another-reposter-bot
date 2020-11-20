package main

import (
	//"encoding/json"
	//"fmt"
	//"io/ioutil"
	"github.com/SevereCloud/vksdk/v2/api"
	//"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/k0kubun/pp"
	"log"
	//"os"
	//"strconv"
)

func checkVKWallGet(id int) {
	log.Printf("Checking VK wall. %d\n", id)
	vk := api.NewVK(cfg.VKToken)

	wall, err := vk.WallGet(api.Params{
		"owner_id": id,
		"count":    4,
		"filter":   "owner",
	})
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range wall.Items {
		if len(v.CopyHistory) != 0 {
			continue
		}
		pp.Println(k, v.Text)
		pp.Println(k, len(v.CopyHistory), len(v.Attachments))
	}
}

func checkVKStatus(id int) {
	log.Printf("Checking VK status %v\n", id)

}
