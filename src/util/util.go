package util

import (
	"log"
	"math/rand"
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
