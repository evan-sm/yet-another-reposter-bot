package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"time"
)

const Aid = 1233
const UserAgent = "com.zhiliaoapp.musically"

var secUIDReg = regexp.MustCompile(`(?m)secUid":"(.*?)"`)
var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func getSecUserID(username string) (string, error) {
	req := &fasthttp.Request{}
	res := &fasthttp.Response{}

	req.Header.SetMethod(http.MethodGet)
	req.SetRequestURI("https://www.tiktok.com/@" + username)
	req.Header.SetUserAgent(UserAgent)

	err := fasthttp.Do(req, res)
	if err != nil {
		return "", err
	}

	matches := secUIDReg.FindStringSubmatch(res.String())
	if len(matches) != 2 {
		return "", errors.New("no matches")
	}

	return matches[1], nil
}

func getLikedVideos(u User, count int) ([]video, error) {
	var secUserID string
	var err error

	if cfg.TikTokSecUserID == "" {
		secUserID, err = getSecUserID(u.Social.TikTok)
		if err != nil {
			log.Fatalln("SecUID", err)
		}
	}

	req := &fasthttp.Request{}
	res := &fasthttp.Response{}

	req.Header.SetMethod(http.MethodGet)
	req.SetRequestURI("https://api16-normal-c-alisg.tiktokv.com/aweme/v1/aweme/favorite/?" +
		fmt.Sprintf("aid=%d&device_id=%d&sec_user_id=%s&count=%d",
			Aid, 1000000000+seededRand.Intn(1000000000), secUserID, count))
	req.Header.SetUserAgent(UserAgent)

	err = fasthttp.Do(req, res)
	if err != nil {
		return nil, err
	}

	var tiktokResp awemeFavoriteResponse
	err = json.Unmarshal(res.Body(), &tiktokResp)
	if err != nil {
		return nil, err
	}

	if tiktokResp.StatusMessage != "" {
		return nil, errors.New(tiktokResp.StatusMessage)
	}

	var videos []video
	for _, v := range tiktokResp.AwemeList {
		if len(v.Video.PlayAddr.URLList) == 0 {
			continue
		}
		videos = append(videos, video{
			ID:          v.ID,
			ShareURL:    v.ShareURL,
			DownloadURL: v.Video.PlayAddr.URLList[0],
		})
	}

	for i, j := 0, len(videos)-1; i < j; i, j = i+1, j-1 {
		videos[i], videos[j] = videos[j], videos[i]
	}

	return videos, nil
}

type video struct {
	ID          string
	ShareURL    string
	DownloadURL string
}

type awemeFavoriteResponse struct {
	StatusMessage string `json:"status_msg"`
	AwemeList     []struct {
		ID       string `json:"aweme_id"`
		ShareURL string `json:"share_url"`
		Video    struct {
			PlayAddr struct {
				URLList []string `json:"url_list"`
			} `json:"play_addr"`
		} `json:"video"`
	} `json:"aweme_list"`
}
