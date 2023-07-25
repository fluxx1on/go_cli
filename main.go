package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/fluxx1on/go_cli/proto"
	"google.golang.org/grpc"
)

var (
	rootDir, _ = os.Getwd()
	mediaDir   = rootDir + "/media/"
)

const serverAddress = "127.0.0.1:50051"

var asyncFlag = flag.Bool("async", false, "Use async mode for ListThumbnails")
var urlList = flag.String("urls", "", "Comma-separated list of URLs for GetThumbnailRequest")

func main() {
	flag.Parse()

	os.Mkdir(mediaDir, 0755)

	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := proto.NewThumbnailServiceClient(conn)

	urls := strings.Split(*urlList, ",")

	if *asyncFlag {
		requests := make([]*proto.GetThumbnailRequest, 0)
		for _, url := range urls {
			request := &proto.GetThumbnailRequest{
				Url: url,
			}
			requests = append(requests, request)
		}

		listRequest := &proto.ListThumbnailRequest{
			Requests: requests,
		}

		response, err := client.ListThumbnail(context.Background(), listRequest)
		if err != nil {
			log.Fatalf("ListThumbnail failed: %v", err)
		}

		ReadResponse(response.Thumbnails...)
	} else {

		for _, url := range urls {
			request := &proto.GetThumbnailRequest{
				Url: url,
			}

			response, err := client.GetThumbnail(context.Background(), request)
			if err != nil {
				log.Fatalf("GetThumbnail failed: %v", err)
			}
			ReadResponse(response)

		}
	}
}

func ReadResponse(thumbResp ...*proto.ThumbnailResponse) {
	for i, resp := range thumbResp {
		if thumb := resp.GetThumbnail(); thumb != nil {
			file, _ := os.Create(mediaDir + thumb.Title + ".jpg")
			defer file.Close()

			if _, err := file.Write(thumb.File); err != nil {
				log.Printf("%d. file didn't write correctly", i+1)
			}
			log.Printf("%d. %v - %v", i+1, thumb.ChannelTitle, thumb.Title)
		} else {
			err := resp.GetError()
			log.Printf("%d. %v, %v", i+1, err.ErrorMessage, err.Url)
		}
	}
}
