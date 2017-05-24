package util

import (
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

func SaveToFile(name string, data []byte) {
	file, err := os.Create(name)
	if err != nil {
		log.Println("Cannot create file: ", name, " ", err.Error())
		return
	}

	file.Write(data)
	defer file.Close()
}

func GenerateDeviceID() string {
	var id = "0x"

	rand.Seed(time.Now().Unix())
	result := rand.Perm(13)
	for _, i := range result {
		id = id + strconv.Itoa(i)
	}

	return id
}

func Download(link, name string) error {
	if name == "" || link == "" {
		return errors.New("Invalid paramenters: " + link + ":" + name)
	}

	resp, err := http.Get(link)
	if err != nil {
		return errors.New("Do http request error")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("Do http request error")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Read http request body error: " + err.Error())
	}

	if err := ioutil.WriteFile(name, body, 0644); err != nil {
		return errors.New("Write content to file error: " + err.Error())
	}

	return nil
}
