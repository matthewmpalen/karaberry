package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var songRegex = regexp.MustCompile(`(.*?)-(.*?)-(.*?)\.(mkv|mp4|webm)`)

func GenerateSongList() {
	log.Println("Generating song list")
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

	brokenFile, readErr := os.Open(Config.MediaFolder + "/broken")
	if readErr != nil {
		log.Fatalf("Cannot open broken song file: %v\n", readErr)
	}
	defer brokenFile.Close()

	broken := map[string]struct{}{}
	scanner := bufio.NewScanner(brokenFile)
	for scanner.Scan() {
		broken[scanner.Text()] = struct{}{}
	}

	for _, file := range files {
		match := songRegex.FindStringSubmatch(file.Name())
		if len(match) != 5 {
			log.Printf("Bad format: %s\n", file.Name())
		} else {
			filename := match[0]
			artist := strings.TrimSpace(match[1])
			songName := strings.Replace(match[2], "(Karaoke Version)", "", -1)
			songName = strings.Replace(songName, "(Karaoke)", "", -1)
			songName = strings.TrimSpace(songName)
			youtubeID := match[3]
			// Is the Pi unable to play the file?
			// Override filename so we stream instead
			if _, ok := broken[youtubeID]; ok {
				filename = ""
			}
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

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/history", HistoryHandler)
	http.HandleFunc("/idle", IdleHandler)
	http.HandleFunc("/queue", QueueHandler)
	http.HandleFunc("/skip", SkipHandler)
	http.HandleFunc("/songs", SongsHandler)
	http.HandleFunc("/ws", WebsocketHandler)
	go hub.run()

	log.Fatal(http.ListenAndServe(":8080", nil))

	log.Println("Done")
}
