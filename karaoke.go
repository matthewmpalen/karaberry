package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
)

type (
	Song struct {
		ID        int    `json:"id"`
		Artist    string `json:"artist"`
		Name      string `json:"name"`
		YoutubeID string `json:"youtube_id"`
		filename  string `json:"filename"`
	}

	SongHistory struct {
		songs []*Song
	}

	SongQueue struct {
		q         []*Song
		size      int
		lock      *sync.Mutex
		addlisten chan int
	}

	Karaoke struct {
		history SongHistory
		Queue   chan *Song
	}
)

var (
	songList []*Song
	karaoke  = Karaoke{
		history: SongHistory{songs: []*Song{}},
		Queue:   make(chan *Song, Config.QueueSize),
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
		songList = append(songList, &Song{i, row[0], row[1], row[2], filename})
	}
}

func newPlayCmd(song *Song) *exec.Cmd {
	if song.filename != "" {
		if _, err := os.Stat(song.filename); err == nil {
			switch Config.MediaPlayer {
			case "vlc":
				return exec.Command("vlc", "--play-and-exit", "--fullscreen", "-I", "dummy", song.filename)
			default:
				return exec.Command("omxplayer", song.filename)
			}
		}
	} else if song.YoutubeURL() != "" {
		log.Println("Streaming file")

		switch Config.MediaPlayer {
		case "vlc":
			return exec.Command("vlc", "--play-and-exit", "--fullscreen", "-I", "dummy", song.YoutubeURL())
		default:
			return exec.Command(Config.ScriptsFolder+"/omxplayer/start.sh", song.YoutubeURL())
		}
	}

	log.Printf("Song %d: Could not build command\n", song.ID)
	return nil
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
func (sh *SongHistory) Add(s *Song) {
	sh.songs = append(sh.songs, s)
}

func (sh SongHistory) Songs() []*Song {
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

func (k Karaoke) Play(song *Song) {
	log.Printf("Playing %s: %s (%s)\n", song.Artist, song.Name, song.YoutubeURL())
	if cmd := newPlayCmd(song); cmd != nil {
		cmd.Run()
	}
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
