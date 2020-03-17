package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fabiodcorreia/gotube"
)

func main() {
	target := "https://www.youtube.com/watch?v=urarTyKn9cg2"

	v, err := gotube.GetVideoDetails(target)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(v.Streams[0].ContentLength)

	file, err := os.Create("./Default-" + v.Title + v.Streams[0].Extension)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	file2, err := os.Create("./Custom-" + v.Title + v.Streams[0].Extension)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	start := time.Now()

	_, err = v.DownloadDefault(file)
	if err != nil {
		log.Fatalln(err)
	}

	end1 := time.Now()

	_, err = v.Download(file2)
	if err != nil {
		log.Fatalln(err)
	}

	end2 := time.Now()

	fmt.Printf("Download with Default Client:  %d\n", (end1.Unix() - start.Unix()))
	fmt.Printf("Download with Custom Client: %d\n", (end2.Unix() - end1.Unix()))

}

// go run cmd/main.go  1,23s user 1,30s system 12% cpu 20,883 total
// go run cmd/main.go  0,85s user 0,85s system 12% cpu 13,614 total
// go run cmd/main.go  0,92s user 0,90s system 13% cpu 13,715 total
// go run cmd/main.go  1,61s user 1,77s system 13% cpu 25,408 total
