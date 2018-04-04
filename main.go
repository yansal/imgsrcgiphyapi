package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY must be set")
	}

	log.Fatal(http.ListenAndServe(":"+port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get(fmt.Sprintf("https://api.giphy.com/v1/gifs/random?api_key=%s&tag=%s", apiKey, r.URL.Query().Get("tag")))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var v struct {
			Data struct {
				Images struct {
					Original struct {
						URL string
					}
				}
			}
		}
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err = http.Get(v.Data.Images.Original.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if _, err := io.Copy(w, resp.Body); err != nil {
			log.Print(err)
		}
	})))
}
