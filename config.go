package main

//import "github.com/caarlos0/env/v6"

type config struct {
	//DBAddr string `env:"DB_ADDR" envDefault:"localhost:6379"`

	TelegramToken string `env:"TG_TOKEN,required"`
	ChannelID     int64  `env:"CHANNEL_ID,required"`

	//TikTokUsername  string `env:"TIKTOK_USERNAME,required"`
	TikTokSecUserID string `env:"TIKTOK_SEC_USER_ID"`

	VKToken string `env:"VK_TOKEN,required"`

	Passcode string `env:"PASSCODE,required"`
}

var cfg config
