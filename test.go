package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

// youtube api key
var developerKey string

func main() {
	fmt.Println(`
	██╗   ██╗██╗████████╗
	╚██╗ ██╔╝██║╚══██╔══╝
	 ╚████╔╝ ██║   ██║   
	  ╚██╔╝  ██║   ██║   
	   ██║   ██║   ██║   
	   ╚═╝   ╚═╝   ╚═╝   
	
    `)

	reader := bufio.NewReader(os.Stdin)

	var queryInput, maxResultsInput string

	fmt.Printf("Your Youtube API Key: ")
	fmt.Scan(&developerKey)
	fmt.Printf("\n Search: ")
	queryInput, _ = reader.ReadString('\n')
	queryInput = strings.Replace(queryInput, "\n", "", -1)
	fmt.Printf(" Max Results (number): ")
	fmt.Scan(&maxResultsInput)
	fmt.Printf(" Searching ... \n\n\t Query: %v\t\t Max Results:  %v\n\n", queryInput, maxResultsInput)
	maxResults, err := strconv.ParseInt(maxResultsInput, 10, 64)
	if err != nil {
		maxResults = 5
	}
	videos := getVideos(queryInput, maxResults)
	mapByIndex := make(map[int]string)
	inumber := 1
	for key, val := range videos {
		mapByIndex[inumber-1] = key
		fmt.Printf("(%v)\t[%v]\t%v\n", inumber, val.PublishedAt, val.Title)
		inumber++
	}
	var lineNumberInput string
	fmt.Print("\n Select Line Number: ")
	fmt.Scan(&lineNumberInput)
	lineNumber, err := strconv.Atoi(lineNumberInput)
	if err != nil {
		log.Fatal(err)
	}
	video := videos[mapByIndex[lineNumber]]
	fmt.Printf("\n Title: %v\n Channel: %v\n Date: %v\n Description: %v\n VideoId: %v\n  ChannelId: %v\n Link: https://www.youtube.com/watch?v=%v\n", video.Title, video.ChannelTitle, video.PublishedAt, video.Description, mapByIndex[lineNumber], video.ChannelId, mapByIndex[lineNumber])
}

func getVideos(query string, MaxResults int64) map[string]*youtube.SearchResultSnippet {
	flag.Parse()
	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}
	service, err := youtube.New(client)
	if err != nil {
		log.Fatalln(err)
	}

	call := service.Search.List([]string{"id", "snippet"}).Q(*flag.String("query", query, "Search term")).MaxResults(*flag.Int64("max-results", MaxResults, "Max Youtube results")).Type("video")
	res, err := call.Do()
	if err != nil {
		log.Fatalln(err)
	}

	videos := make(map[string]*youtube.SearchResultSnippet)

	for _, item := range res.Items {
		videos[item.Id.VideoId] = item.Snippet
	}

	return videos
}
