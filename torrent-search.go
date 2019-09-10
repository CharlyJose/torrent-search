package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// YTS is YTS API Response
type YTS struct {
	Status        string `json:"status"`
	StatusMessage string `json:"status_message"`
	Data          struct {
		MovieCount int `json:"movie_count"`
		Limit      int `json:"limit"`
		PageNumber int `json:"page_number"`
		Movies     []struct {
			ID                      int      `json:"id"`
			URL                     string   `json:"url"`
			ImdbCode                string   `json:"imdb_code"`
			Title                   string   `json:"title"`
			TitleEnglish            string   `json:"title_english"`
			TitleLong               string   `json:"title_long"`
			Slug                    string   `json:"slug"`
			Year                    int      `json:"year"`
			Rating                  float32  `json:"rating"`
			Runtime                 int      `json:"runtime"`
			Genres                  []string `json:"genres"`
			Summary                 string   `json:"summary"`
			DescriptionFull         string   `json:"description_full"`
			Synopsis                string   `json:"synopsis"`
			YtTrailerCode           string   `json:"yt_trailer_code"`
			Language                string   `json:"language"`
			MpaRating               string   `json:"mpa_rating"`
			BackgroundImage         string   `json:"background_image"`
			BackgroundImageOriginal string   `json:"background_image_original"`
			SmallCoverImage         string   `json:"small_cover_image"`
			MediumCoverImage        string   `json:"medium_cover_image"`
			LargeCoverImage         string   `json:"large_cover_image"`
			State                   string   `json:"state"`
			Torrents                []struct {
				URL              string `json:"url"`
				Hash             string `json:"hash"`
				Quality          string `json:"quality"`
				Type             string `json:"type"`
				Seeds            int    `json:"seeds"`
				Peers            int    `json:"peers"`
				Size             string `json:"size"`
				SizeBytes        int    `json:"size_bytes"`
				DateUploaded     string `json:"date_uploaded"`
				DateUploadedUnix int    `json:"date_uploaded_unix"`
			} `json:"torrents"`
			DateUploaded     string `json:"date_uploaded"`
			DateUploadedUnix int    `json:"date_uploaded_unix"`
		} `json:"movies"`
	} `json:"data"`
	Meta struct {
		ServerTime     int    `json:"server_time"`
		ServerTimezone string `json:"server_timezone"`
		APIVersion     int    `json:"api_version"`
		ExecutionTime  string `json:"execution_time"`
	} `json:"@meta"`
}

func main() {
	Printer(readJSON(request(userInput())))
}

func userInput() string {
	if len(os.Args) < 2 {
		fmt.Println("At least one argument expected")
		os.Exit(0)
	}
	movieName := os.Args[1]
	url := "https://yts.am/api/v2/list_movies.json?query_term=" + movieName + ""

	return url
}

func request(url string) *http.Response {
	client := http.Client{
		Timeout: time.Second * 20,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "YTS MOVIE TORRENTS")

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	return res
}

// func response(res *http.Response) []byte {
// 	body, readErr := ioutil.ReadAll(res.Body)
// 	if readErr != nil {
// 		log.Fatal(readErr)
// 	}
// 	return body
// }

func readJSON(res *http.Response) YTS {
	var movies YTS
	jsonParser := json.NewDecoder(res.Body)
	err := jsonParser.Decode(&movies)
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}
	return movies
}

func Printer(movies YTS) {

	fmt.Println("Status: ", movies.Status)
	fmt.Println("Status Message: ", movies.StatusMessage)
	fmt.Println("Movie Count: ", movies.Data.MovieCount)

	/*
		var j int
		for i, movielist := range movies.Data.Movies {
			fmt.Println("URL:", movielist.URL)
			fmt.Println("Title: ", movielist.TitleLong)
			fmt.Println("Year: ", movielist.Year)
			fmt.Println("Rating: ", movielist.Rating)
			fmt.Println("Runtime: ", movielist.Runtime)

			for _, genre := range movielist.Genres {
				fmt.Println("Genre: ", genre)
			}

			for _, torrent := range movielist.Torrents {
				fmt.Println("URL: ", torrent.URL)
				fmt.Println("Quality: ", torrent.Quality)
				fmt.Println("Type: ", torrent.Type)
				fmt.Println("Seeds: ", torrent.Seeds)
				fmt.Println("Peers: ", torrent.Peers)
				fmt.Println("Size: ", torrent.Size)
			}
			j = i
			fmt.Println()
		}
		fmt.Println("Movies: ", j)
	*/

	for _, movielist := range movies.Data.Movies {
		fmt.Printf("TITLE: %6s | RATING: %1f\n", movielist.TitleLong, movielist.Rating)

		fmt.Printf("GENRE: ")
		for _, genre := range movielist.Genres {
			fmt.Printf("%6s ", genre)
		}
		fmt.Println()

		fmt.Printf("TORRENTS\n")
		for _, torrent := range movielist.Torrents {
			fmt.Printf("QUALITY: %6s | SEEDS: %3d | PEERS: %3d\n", torrent.Quality, torrent.Seeds, torrent.Peers)
			fmt.Printf("URL: %6s", torrent.URL)
			fmt.Println()
		}

		fmt.Println()
		fmt.Println()
	}
}
