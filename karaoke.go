package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type (
	Song struct {
		ID        int    `json:"id"`
		Artist    string `json:"artist"`
		Name      string `json:"name"`
		YoutubeID string `json:"youtube_id"`
		Filename  string `json:"filename"`
	}

	SongHistory struct {
		songs []Song
	}

	Karaoke struct {
		history SongHistory
		Queue   chan Song
	}
)

var (
	songList []Song
	karaoke  = Karaoke{
		history: SongHistory{songs: []Song{}},
		Queue:   make(chan Song, Config.QueueSize),
	}
)

func init() {
	songFilename := Config.MediaFolder + "/songs.csv"
	songFile, openErr := os.Open(songFilename)
	if openErr != nil {
		panic("Could not open CSV file")
	}
	defer songFile.Close()

	reader := csv.NewReader(songFile)
	reader.Comma = '|'
	rows, readErr := reader.ReadAll()
	if readErr != nil {
		panic(readErr)
	}

	for i, row := range rows {
		filename := fmt.Sprintf("%s/%s", Config.MediaFolder, row[3])
		songList = append(songList, Song{i, row[0], row[1], row[2], filename})
	}
}

func newPlayCmd(song Song) *exec.Cmd {
	var stream string
	name := song.Filename
	if name != "" {
		if _, err := os.Stat(name); os.IsNotExist(err) {
			name = song.YoutubeURL()
			stream = "stream"
			log.Println("Streaming file")
		}
	}

	switch Config.MediaPlayer {
	case "vlc":
		return exec.Command(Config.ScriptsFolder+"/vlc/start.sh", name)
	default:
		return exec.Command(Config.ScriptsFolder+"/omxplayer/start.sh", name, stream)
	}

}

//============
// Song
func (s Song) YoutubeURL() string {
	return fmt.Sprintf("https://www.youtube.com/watch?v=%s", s.YoutubeID)
}

func (s Song) String() string {
	return fmt.Sprintf("%s: %s", s.Artist, s.Name)
}

//============
// SongHistory
func (sh *SongHistory) Add(s Song) {
	sh.songs = append(sh.songs, s)
}

func (sh SongHistory) Songs() []Song {
	return sh.songs
}

func (sh *SongHistory) String() string {
	return fmt.Sprintf("%v", sh.songs)
}

//============
// Karaoke
func (k Karaoke) History() SongHistory {
	return k.history
}

func (k Karaoke) Play(song Song) {
	log.Printf("Playing %s: %s (%s)\n", song.Artist, song.Name, song.YoutubeURL())
	newPlayCmd(song).Run()
	log.Printf("Finished: %s: %s\n", song.Artist, song.Name)
}

func (k *Karaoke) Run() {
	for {
		select {
		case song := <-k.Queue:
			log.Printf("%s\n", song)
			k.Play(song)
			k.history.Add(song)
		}
	}
}
