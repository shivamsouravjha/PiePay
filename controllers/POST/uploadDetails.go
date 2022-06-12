package post

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"piepay/services/es"
	"piepay/structs"
	"strings"
	"time"

	"github.com/olivere/elastic/v7"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

var (
	maxResults = flag.Int64("max-results", 25, "Max YouTube results")
	videotype  = flag.String("type", "video", "Type of video")
)

//AIzaSyD9GAb8G4LOSOOPUUTtkOXGuzX_P15bIuM
//AIzaSyCba0XKmRFR0MWmdpYGANHiYhP_k0n5fxs
const developerKey = "AIzaSyCba0XKmRFR0MWmdpYGANHiYhP_k0n5fxs"

func UploadVideoMetaData() {
	startafter := time.Now().Add(-time.Minute * 30).Format(time.RFC3339)

	flag.Parse()

	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	arraya := []string{"id", "snippet"}
	call := service.Search.List(arraya).
		MaxResults(*maxResults).PublishedAfter(startafter).Type(*videotype)
	fmt.Println(startafter, time.Now(), call)
	response, err := call.Do()
	if err != nil {
		errorText := strings.Split(err.Error(), ", ")
		if errorText[1] == "quotaExceeded" {
			fmt.Println(errorText[1])
		}
	} else {
		fmt.Println(err)
		videos := []structs.Video{}
		for _, item := range response.Items {
			switch item.Id.Kind {
			case "youtube#video":
				videos = append(videos, structs.Video{
					ID:          structs.VideoID{item.Id.VideoId},
					PublishedAt: item.Snippet.PublishedAt,
					ChannelName: item.Snippet.ChannelTitle,
					Title:       item.Snippet.Title,
					Description: item.Snippet.Description,
					ChannelID:   item.Snippet.ChannelId,
				})
			}
		}

		printIDs("Videos", videos)
	}
}

func printIDs(hehe string, matches []structs.Video) {
	bulk := es.Client().Bulk()

	for i := range matches {
		fmt.Println(matches[i])
		idStr := matches[i].ID.ID
		req := elastic.NewBulkIndexRequest().Id(idStr).Index("video").Doc(matches[i])
		bulk = bulk.Add(req)
	}

	resp, _ := bulk.Do(context.Background())
	fmt.Println(resp)
	fmt.Printf("\n\n")

}
