package main

import (
	pb "KseniaErsh/VideoFromPlaylist/proto"
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr       = flag.String("addr", "localhost:50051", "the address to connect to")
	playlistId = ""
)

//Получение списка видео от backend-модуля
func getVideoList() []string {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return nil
	}
	defer conn.Close()
	c := pb.NewGetVideoListClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetPlaylistItems(ctx, &pb.Request{PlaylistID: playlistId})
	if err != nil {
		log.Fatalf("%v", err)
		return nil
	}
	result := r.GetVideoList()
	if len(result) < 1 {
		return nil
	}
	return result
}

func main() {
	http.HandleFunc("/", home_page)
	http.HandleFunc("/get/", get_page)
	http.ListenAndServe(":8080", nil)
}

//Отображение основной страницы
func home_page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./webServer/homePage.html")
	if err != nil {
		fmt.Println(err)
	}
	tmpl.Execute(w, nil)
}

//Отображение страницы с результатом
func get_page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./webServer/getPage.html")
	if err != nil {
		fmt.Println(err)
	}
	id := r.FormValue("playlistId")
	playlistId = id
	videoList := getVideoList()
	tmpl.Execute(w, videoList)

}
