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
	"text/template"
)

const (
	listItemsSelector = "//div[@id='main']/ul/li"
	itemBodySelector  = "//div[@id='main']"
)

var (
	ItemInListIdRe        = regexp.MustCompile(`[0-9]+$`)
	ItemInDescriptionIdRe = regexp.MustCompile(`b/([0-9]+)/read$`)
	stripSpacesRe         = regexp.MustCompile(`\s+`)
)

type ListItem struct {
	Title   string
	Authors []string
	ID      string
}

func (item *ListItem) String() string {
	return fmt.Sprintf("%s: %s <%s>", item.ID, item.Title, strings.Join(item.Authors, ", "))
}

func (info *InfoResult) String() string {
	const tpl = `
	{{.Title}}
	ID: {{.ID}}
	Size: {{.Size}}
	Formats: {{range .Formats}} {{.}} {{end}}

	{{.Annotation}}
	`
	t, err := template.New("bookInfo").Parse(tpl)
	check(err)
	buf := &bytes.Buffer{}
	err = t.Execute(buf, info)
	check(err)

	return buf.String()
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
			ID:      ItemInListIdRe.FindString(itemHref),
			Title:   title.String(),
			Authors: getAuthors(authorNodes),
		})
	}

	return result, nil
}

func ParseInfo(stream io.Reader) (result *InfoResult, err error) {
	doc, _ := htmlquery.Parse(stream)

	list := htmlquery.Find(doc, itemBodySelector)
	if list == nil {
		return nil, errors.New("item body not found")
	}

	getText := func(node *html.Node) (data string) {
		buf := &bytes.Buffer{}
		if node != nil {
			collectText(node, buf)
		}
		return buf.String()
	}

	id := getID(doc)
	if id == "" {
		// It is not item page
		return nil, errors.New("item not found")
	}

	result = &InfoResult{
		ID:         id,
		Title:      getText(htmlquery.FindOne(doc, "//div[@id='main']/h1/text()")),
		Genre:      getText(htmlquery.FindOne(doc, "//p[@class=\"genre\"]")),
		Annotation: getText(htmlquery.FindOne(doc, "//div[@id='main']/p/text()")),
		Size:       getText(htmlquery.FindOne(doc, "//span[@style=\"size\"]/text()")),
		Formats:    getFormats(doc),
	}
	return
}

func getAuthors(nodes []*html.Node) (authors []string) {
	for _, node := range nodes {
		authors = append(authors, htmlquery.InnerText(node))
	}
	return authors
}

func getFormats(doc *html.Node) (formats []string) {
	for _, format := range validFormats {
		available := htmlquery.FindOne(doc, "//a[contains(@href, \""+format+"\")]")
		if available != nil {
			formats = append(formats, format)
		}
	}
	return formats
}

func getID(doc *html.Node) (ID string) {
	readUrl := htmlquery.SelectAttr(
		htmlquery.FindOne(doc, "//a[contains(@href, \"read\")]"),
		"href",
	)
	if readUrl != "" {
		return ItemInDescriptionIdRe.FindStringSubmatch(readUrl)[1]
	}
	return
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
