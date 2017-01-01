package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	context := map[string]interface{}{"songs": songList}
	RenderTemplate(w, "home", context)
}

func SongsHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{
		"count": len(songList),
		"songs": songList,
	}
	JSON(w, resp, http.StatusOK)
}

func QueueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResp := map[string]string{"error": "Requires POST"}
		JSON(w, errorResp, http.StatusMethodNotAllowed)
		return
	}

	id := r.PostFormValue("songId")
	songID, parseErr := strconv.Atoi(id)
	if parseErr != nil || songID < 0 || songID > len(songList)-1 {
		msg := fmt.Sprintf("Invalid song: %s", id)
		errorResp := map[string]string{"error": msg}
		JSON(w, errorResp, http.StatusBadRequest)
		return
	}

	song := songList[songID]
	karaoke.Queue <- song

	resp := map[string]interface{}{
		"count": len(karaoke.Queue),
		"added": song,
	}
	JSON(w, resp, http.StatusOK)
}

func HistoryHandler(w http.ResponseWriter, r *http.Request) {
	songs := karaoke.History().Songs()
	resp := map[string]interface{}{
		"count": len(songs),
		"songs": songs,
	}
	JSON(w, resp, http.StatusOK)
}
