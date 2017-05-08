package main

import (
	"log"
	"os"
)

func main() {
	_, err := os.Open("liwei.txt")
	if err != nil {
		log.Println(err.Error())
		//		return
	}

	f, err := os.Stat("test.txt")
	log.Println(f)
	log.Println(err.Error())
	if os.IsNotExist(err) {
		log.Println("File not exist")
	}

	n, err := os.Create("test.txt")
	if err != nil {
		log.Println("Create file failed")
	}
	n.WriteString("Hello world")

}
