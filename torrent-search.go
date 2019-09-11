package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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
			URL       string   `json:"url"`
			ImdbCode  string   `json:"imdb_code"`
			TitleLong string   `json:"title_long"`
			Rating    float32  `json:"rating"`
			Runtime   int      `json:"runtime"`
			Genres    []string `json:"genres"`
			Language  string   `json:"language"`
			Torrents  []struct {
				URL          string `json:"url"`
				Quality      string `json:"quality"`
				Type         string `json:"type"`
				Seeds        int    `json:"seeds"`
				Peers        int    `json:"peers"`
				Size         string `json:"size"`
				DateUploaded string `json:"date_uploaded"`
			} `json:"torrents"`
			DateUploaded string `json:"date_uploaded"`
		} `json:"movies"`
	} `json:"data"`
}

func main() {
	Printer(ReadJSON(RequestServer(GetUserInput())))
}

// GetUserInput get the user input
func GetUserInput() string {
	if len(os.Args) < 2 {
		fmt.Println("At least one argument expected")
		os.Exit(0)
	}
	movieName := os.Args[1]
	url := "https://yts.am/api/v2/list_movies.json?query_term=" + movieName + ""

	return url
}

// RequestServer request server
func RequestServer(url string) *http.Response {
	client := http.Client{
		Timeout: time.Second * 20,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err)

		os.Exit(0)
	}

	req.Header.Set("User-Agent", "YTS MOVIE TORRENTS")

	res, getErr := client.Do(req)
	if getErr != nil {
		// fmt.Println(getErr)
		if strings.HasSuffix(fmt.Sprint(getErr), "no such host") {
			fmt.Println("You might not be connected to the internet")
		}
		if strings.HasSuffix(fmt.Sprint(getErr), "access permissions.") {
			fmt.Println("This application might be blocked by your firewall")
		}
		os.Exit(0)
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

// ReadJSON parse the JSON
func ReadJSON(res *http.Response) YTS {
	var movies YTS
	jsonParser := json.NewDecoder(res.Body)
	err := jsonParser.Decode(&movies)
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}
	if movies.Status != "ok" {
		fmt.Println("Some error at other side")
		os.Exit(0)
	}
	if movies.Data.MovieCount == 0 {
		fmt.Println("No Movies found")
		os.Exit(0)
	}
	return movies
}

// Printer prints the JSON response
func Printer(movies YTS) {

	fmt.Println("Movie Count: ", movies.Data.MovieCount)

	for _, movielist := range movies.Data.Movies {
		fmt.Printf("%9s %-6s | LANGUAGE: %6s | RUNTIME: %3d\n", "TITLE:", movielist.TitleLong, movielist.Language, movielist.Runtime)
		fmt.Printf("%9s %-.1f\n", "RATING:", movielist.Rating)
		fmt.Printf("%9s ", "GENRE:")
		for _, genre := range movielist.Genres {
			fmt.Printf("%6s ", genre)
		}
		fmt.Println()

		fmt.Printf("TORRENTS: \n")
		for i, torrent := range movielist.Torrents {
			fmt.Printf("\t%d%2sQUALITY: %5s | TYPE: %6s | SEEDS: %3d | PEERS: %3d\n", i+1, ". ", torrent.Quality, torrent.Type, torrent.Seeds, torrent.Peers)
			fmt.Printf("\t %2sURL: %6s", " ", torrent.URL)
			fmt.Println()
		}
		fmt.Println("\n")
	}

}
