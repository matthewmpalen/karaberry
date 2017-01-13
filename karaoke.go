package main

import (
	"encoding/csv"
	"fmt"
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
	filename := Config.MediaFolder + "/songs.csv"
	file, openErr := os.Open(filename)
	if openErr != nil {
		panic("Could not open CSV file")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '|'
	rows, readErr := reader.ReadAll()
	if readErr != nil {
		panic(readErr)
	}

	for i, row := range rows {
		songList = append(songList, Song{i, row[0], row[1], row[2], ""})
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
	fmt.Printf("Playing %s: %s (%s)\n", song.Artist, song.Name, song.YoutubeURL())
	var cmd *exec.Cmd
	//cmd := exec.Command("vlc", "--play-and-exit", "--fullscreen", "-I", "dummy", song.YoutubeURL())
	if song.Filename != "" {
		cmd = exec.Command("omxplayer", song.Filename)
	} else {
		cmd = exec.Command("./stream.sh", song.YoutubeURL())
	}

	fmt.Printf("Command: %s\n", cmd)
	cmd.Run()
	fmt.Printf("Finished: %s: %s\n", song.Artist, song.Name)
}

func (k *Karaoke) Run() {
	for {
		select {
		case song := <-k.Queue:
			fmt.Printf("%s\n", song)
			k.Play(song)
			k.history.Add(song)
		}
	}
}
