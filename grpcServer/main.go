package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"

	pb "github.com/KseniaErsh/VideoFromPlaylist/proto"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type GRPCServer struct {
	pb.UnsafeGetVideoListServer
}

func (s *GRPCServer) GetPlaylistItems(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	playlistId := req.GetPlaylistID()
	log.Printf("Received: %v", playlistId)
	videoInfo, err := GetPlaylistItems(playlistId)
	if err != nil {
		return &pb.Response{
			VideoList: nil,
		}, nil
	}
	return &pb.Response{
		VideoList: videoInfo,
	}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGetVideoListServer(s, &GRPCServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

const YOUTUNE_PLAYLIST_ITEMS_URL = "https://www.googleapis.com/youtube/v3/playlistItems"
const YOUTUBE_API_KEY = "AIzaSyBvw_EMjvevGfsD9BlblDmlvFZ6fue7vIs"
const YOUTUBE_VIDEO_URL = "https://www.youtube.com/watch?v="

type RestResponse struct {
	Items []Item `json:"items"`
}

type Item struct {
	Snippet SnippetInfo `json:"snippet"`
}

type SnippetInfo struct {
	Title      string       `json:"title"`
	ResourceId resourceInfo `json:"resourceId"`
}

type resourceInfo struct {
	VideoID string `json:"videoId"`
}

//Получение списка видео
func GetPlaylistItems(playlistID string) ([]string, error) {
	items, err := retrieveVideos(playlistID)
	if err != nil {
		return nil, err
	}
	if len(items) < 1 {
		return nil, errors.New("Playlist not found")
	}
	result := make([]string, len(items))
	for i := range items {
		title_video := items[i].Snippet.Title
		if title_video == "" {
			break
		}
		url_video := YOUTUBE_VIDEO_URL + items[i].Snippet.ResourceId.VideoID
		result[i] = title_video + ": " + url_video
	}
	return result, nil
}

//Получение информации о каждом видео из плейлиста
func retrieveVideos(playlistID string) ([]Item, error) {
	req, err := makeRequest(playlistID, 50)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return restResponse.Items, nil
}

//Создание запроса
func makeRequest(playlistID string, maxResult int) (*http.Request, error) {
	req, err := http.NewRequest("GET", YOUTUNE_PLAYLIST_ITEMS_URL, nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Add("part", "snippet")
	query.Add("maxResults", strconv.Itoa(maxResult))
	query.Add("playlistId", playlistID)
	query.Add("key", YOUTUBE_API_KEY)
	req.URL.RawQuery = query.Encode()
	return req, nil
}
