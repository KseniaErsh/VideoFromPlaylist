package main

import (
	pb "GitHub/VideoFromPlaylist/proto"
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr       = flag.String("addr", "localhost:50051", "the address to connect to")
	playlistId = "PLGtCetCIU8w2rFUP0CxdhRz9k-yRAmJcm"
)

func getVideoList() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGetVideoListClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetPlaylistItems(ctx, &pb.Request{PlaylistID: playlistId})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	result := r.GetVideoList()
	for i := range result {
		fmt.Println(result[i])
	}
}

func main() {
	getVideoList()
}
