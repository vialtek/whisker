package main

import (
	"log"

	"github.com/vialtek/whisker/whisker"
)

func main() {
	log.Println("Whisker is starting...")

	whisker := whisker.NewNode()
	whisker.Run()
}
