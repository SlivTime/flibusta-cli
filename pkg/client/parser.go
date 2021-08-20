package client

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"io"
	"regexp"
	"strings"
)

const listItemsSelector = "//div[@id='main']/ul/li"

var itemIdRe = regexp.MustCompile(`[0-9]+`)

type ListItem struct {
	Title   string
	Authors []string
	ID      string
}

func (item *ListItem) String() string {
	return fmt.Sprintf("%s: %s <%s>", item.ID, item.Title, strings.Join(item.Authors, ", "))
}

func ParseSearch(stream io.Reader) (result []ListItem, err error) {
	doc, err := htmlquery.Parse(stream)
	if err != nil {
		return
	}

	result = []ListItem{}

	list := htmlquery.Find(doc, listItemsSelector)
	for _, listItem := range list {
		titleNode := htmlquery.FindOne(listItem, "//a[1]")
		authorNodes := htmlquery.Find(listItem, "//a[position()>1]")
		itemHref := htmlquery.SelectAttr(titleNode, "href")
		result = append(result, ListItem{
			ID:      itemIdRe.FindString(itemHref),
			Title:   htmlquery.InnerText(titleNode),
			Authors: getAuthors(authorNodes),
		})
	}

	return result, nil
}

func getAuthors(nodes []*html.Node) (authors []string) {
	for _, node := range nodes {
		authors = append(authors, htmlquery.InnerText(node))
	}
	return authors
}
