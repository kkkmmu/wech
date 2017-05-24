package handler

import (
	"gopkg.in/redis.v5"
	"io/ioutil"
	"log"
	"math/rand"
	"message"
	"os"
	"path/filepath"
	"strings"
	"util"
)

var BeautifulHandler = &Handler{}

var redisClient *redis.Client

var LinkFile = "asset/link.text"

var contentFile = "c"

var beautiPrefix = "%be"

func BeautifulHandle(msg *message.LocalMessage) *message.LocalMessage {
	content := string(msg.Content())
	content = strings.TrimSpace(content)
	var local message.LocalMessage
	if strings.HasPrefix(content, "%be") {
		file, err := ioutil.ReadFile(LinkFile)
		if err != nil {
			return local.SetFrom(msg.To()).SetTo(msg.From()).SetType(message.TEXT).SetContent([]byte(err.Error()))
		}
		links := strings.Split(string(file), "\n")
		link := links[rand.Int31n(int32(len(links)))]

		err = util.Download(link, contentFile+filepath.Base(link))
		if err != nil {
			return local.SetFrom(msg.To()).SetTo(msg.From()).SetType(message.TEXT).SetContent([]byte(err.Error()))
		}

		return local.SetFrom(msg.To()).SetTo(msg.From()).SetType(message.IMAGE).SetContent([]byte(contentFile + filepath.Base(link)))
	}

	return nil
}

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	links, err := redisClient.HGetAll("SPIDER:RESULT:CACHE").Result()
	if err != nil {
		log.Println("Error happend when Get picture reposity")
		return
	}

	file, err := os.Create(LinkFile)
	if err != nil {
		log.Println("Error happend when create link file")
		return
	}
	defer file.Close()
	for k, _ := range links {
		file.WriteString(k + "\n")
	}

	BeautifulHandler.SetName("Beautiful").SetType(message.TEXT).SetState(true).SetHandleFunc(BeautifulHandle)
}
