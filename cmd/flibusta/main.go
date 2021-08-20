package main

import (
	"flibusta-cli/pkg/client"
	"fmt"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const defaultBookFormat = "mobi"

func commandSearch(context *cli.Context) error {
	query := strings.Join(context.Args().Slice(), " ")
	flibusta, err := client.FromEnv()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("search book: ", query)
	searchResult, err := flibusta.Search(query)
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range searchResult {
		fmt.Println(item.String())
	}
	return nil
}

func commandGet(context *cli.Context) error {
	bookID := context.Args().First()

	flibusta, err := client.FromEnv()
	if err != nil {
		log.Fatal(err)
	}
	bookFormat := context.String("format")
	fmt.Printf("get book <%s> in `%s` format\n", bookID, bookFormat)
	result, err := flibusta.Download(bookID, bookFormat)
	if err != nil {
		log.Fatal(err)
	}

	if result.Name == "" {
		result.Name = fmt.Sprintf("%s.%s", bookID, bookFormat)
	}
	err = ioutil.WriteFile(result.Name, result.File, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("File saved at", result.Name)
	return nil
}

func main() {
	app := &cli.App{
		Commands: cli.Commands{
			&cli.Command{
				Name:    "search",
				Aliases: []string{"s"},
				Usage:   "Search book",
				Action:  commandSearch,
			},
			&cli.Command{
				Name:    "get",
				Aliases: []string{"g"},
				Usage:   "Get book",
				Action:  commandGet,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "format",
						Aliases: []string{"f"},
						Value:   defaultBookFormat,
						Usage:   "Format to download: mobi|epub|fb2",
						EnvVars: []string{"FLIBUSTA_PREFERRED_FORMAT"},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
