package main

import (
	"flibusta-go/internal/env"
	"flibusta-go/pkg/client"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	env.Load()

	proxyUrl := os.Getenv("PROXY_URL")

	fmt.Println("Proxy:", proxyUrl)
	flibusta, err := client.FromEnv()
	if err != nil {
		log.Fatal("Bad sign")
	}

	searchResult, err := flibusta.Search("путь джедая")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(searchResult)

	bookID := "42053"
	bookFormat := "mobi"
	result, err := flibusta.Download(bookID, bookFormat)
	if err != nil {
		log.Fatal(err)
	}
	outFileName := result.Name

	if outFileName == "" {
		outFileName = fmt.Sprintf("%s.%s", bookID, bookFormat)
	}
	err = ioutil.WriteFile(outFileName, result.File, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("File saved at", outFileName)
}
