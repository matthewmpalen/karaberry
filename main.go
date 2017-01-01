package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	go karaoke.Run()

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/songs", SongsHandler)
	http.HandleFunc("/queue", QueueHandler)
	http.HandleFunc("/history", HistoryHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

	fmt.Println("Done")
}
