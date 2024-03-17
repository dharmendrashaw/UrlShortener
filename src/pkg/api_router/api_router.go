package api_router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/UrlShortener/src/pkg/config"
	"github.com/UrlShortener/src/pkg/db"
	"github.com/UrlShortener/src/pkg/utility"
	"github.com/gorilla/mux"
)

type ShortenUrlRequest struct {
	Url string `json:"url"`
}

func CreateShortUrl(w http.ResponseWriter, r *http.Request) {
	var shorUrl ShortenUrlRequest
	json.NewDecoder(r.Body).Decode(&shorUrl)
	defer r.Body.Close()

	if utility.NodeRange.Curr < utility.NodeRange.End && utility.NodeRange.Curr != 0 {
		utility.NodeRange.Curr = utility.NodeRange.Curr + 1
	} else {
		utility.GetRangeNode()
		utility.NodeRange.Curr = utility.NodeRange.Curr + 1
	}

	hash := utility.GetHash(utility.NodeRange.Curr)
	shortUrlCollection := &db.ShortenUrlCollection{
		Url:         shorUrl.Url,
		Hash:        hash,
		CreatedDate: "",
		Clicks:      0,
	}
	shortUrlCollection.Save()
	var shortenUrl = fmt.Sprintf("http://localhost%s/%s", config.AppConfig.ServerPort, hash)
	w.Write([]byte(shortenUrl))
}

func FetchRedirecUrl(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]
	existing := db.FindOneByHash(hash)
	w.Header().Add("location", existing.Url)
	w.WriteHeader(301)
}
