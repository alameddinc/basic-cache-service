package main

import (
	"github.com/alameddinc/ysc/routes"
	"github.com/alameddinc/ysc/storage"
	"log"
)

func main() {
	if err := storage.ReadFileStorage(); err != nil {
		log.Fatalln(err)
		return
	}
	go storage.Sync()
	routes.Handler()
}
