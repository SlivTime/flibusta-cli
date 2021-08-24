package client

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"io"
	"regexp"
	"strings"
)

const listItemsSelector = "//div[@id='main']/ul/li"

var itemIdRe = regexp.MustCompile(`[0-9]+$`)
var stripSpacesRe = regexp.MustCompile(`\s+`)

type ListItem struct {
	Title   string
	Authors []string
	ID      string
}

func (item *ListItem) String() string {
	return fmt.Sprintf("%s: %s <%s>", item.ID, item.Title, strings.Join(item.Authors, ", "))
}

func ParseSearch(stream io.Reader) (result *[]ListItem, err error) {
	doc, _ := htmlquery.Parse(stream)

	list := htmlquery.Find(doc, listItemsSelector)
	if list == nil {
		return nil, errors.New("list with items not found")
	}
	result = &[]ListItem{}
	for _, listItem := range list {
		titleNode := htmlquery.FindOne(listItem, "//a[1]")
		authorNodes := htmlquery.Find(listItem, "//a[position()>1]")
		itemHref := htmlquery.SelectAttr(titleNode, "href")

		title := &bytes.Buffer{}
		collectText(titleNode, title)
		*result = append(*result, ListItem{
			ID:      itemIdRe.FindString(itemHref),
			Title:   title.String(),
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

func collectText(n *html.Node, buf *bytes.Buffer) {
	if n.Type == html.TextNode {
		txt := stripSpacesRe.ReplaceAllString(n.Data, ` `)
		buf.WriteString(txt)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectText(c, buf)
	}
}
