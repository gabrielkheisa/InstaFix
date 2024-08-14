package handlers

import (
	scraper "instafix/handlers/scraper"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

var videoProxy string

func init() {
	videoProxy = os.Getenv("VIDEO_PROXY")
	if videoProxy == "" {
		return
	}
	if !(strings.HasPrefix(videoProxy, "http://") || strings.HasPrefix(videoProxy, "https://")) {
		panic("VIDEO_PROXY must start with http:// or https://")
	}
	if !strings.HasSuffix(videoProxy, "/") {
		videoProxy += "/"
	}
}

func Videos(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	mediaNum, err := strconv.Atoi(chi.URLParam(r, "mediaNum"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item, err := scraper.GetData(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to image URL
	if mediaNum > len(item.Medias) {
		return
	}
	videoURL := item.Medias[max(1, mediaNum)-1].URL

	// Redirect to proxy if not TelegramBot in User-Agent
	if strings.Contains(r.Header.Get("User-Agent"), "TelegramBot") {
		http.Redirect(w, r, videoURL, http.StatusFound)
		return
	}
	http.Redirect(w, r, videoProxy+videoURL, http.StatusFound)
	return
}
