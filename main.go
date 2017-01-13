package main

import (
	"encoding/csv"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

var songRegex = regexp.MustCompile(`(.*?) - (.*?)-(.*?)\.(mkv|mp4|webm)`)

func GenerateSongList() {
	files, readErr := ioutil.ReadDir(Config.MediaFolder)
	if readErr != nil {
		log.Fatalf("Could not generate song list: %v\n", readErr)
	}

	csvFile, createErr := os.Create(Config.MediaFolder + "/songs.csv")
	if createErr != nil {
		log.Fatalf("Cannot open song file: %v\n", createErr)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	writer.Comma = '|'
	defer writer.Flush()

	for _, file := range files {
		match := songRegex.FindStringSubmatch(file.Name())
		if len(match) != 5 {
			log.Printf("Bad format: %s\n", file.Name())
		} else {
			filename := match[0]
			artist := match[1]
			songName := match[2]
			youtubeID := match[3]
			row := []string{artist, songName, youtubeID, filename}
			writer.Write(row)
		}
	}

	log.Println("Finished building song list")
}

func main() {
	player := flag.String("player", "omxplayer", "media player [omxplayer|vlc]")
	build := flag.Bool("build", false, "rebuild song list based on saved files")
	flag.Parse()
	if *build {
		GenerateSongList()
		return
	}
	Config.MediaPlayer = *player

	go karaoke.Run()

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/songs", SongsHandler)
	http.HandleFunc("/queue", QueueHandler)
	http.HandleFunc("/history", HistoryHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

	log.Println("Done")
}
