package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/A-ndrey/tg-agent-bot/internal/bot"
)

func main() {
	var isFile bool
	flag.BoolVar(&isFile, "f", false, "send file")
	flag.Parse()

	if len(os.Args) > 1 && os.Args[1] == "init" {
		err := makeConfig()
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	conf, err := bot.ReadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	b := bot.New(conf.Token, conf.ChatID)

	if isFile {
		if err := b.SendFile(os.Stdin); err != nil {
			fmt.Println(err)
			return
		}
	} else {
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := b.SendMessage(string(content)); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func makeConfig() error {
	token := os.Getenv("TG_TOKEN")
	b := bot.New(token, 0)
	chatID, err := b.GetChatID()
	if err != nil {
		return err
	}

	conf := bot.Config{
		Token:  token,
		ChatID: chatID,
	}

	if err := bot.SaveConfig(conf); err != nil {
		return err
	}

	return nil
}
