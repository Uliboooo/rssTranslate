package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/mmcdole/gofeed"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

type Content struct {
	Title string `json:"title"`
	Link  string `json:"link"`
	Date  string `json:"date"`
}
type Contents struct {
	Lst []Content `json:"list"`
}

func (c *Contents) Translate() {
	ctx := context.Background()

	api := os.Getenv("GEMINI_API_KEY")
	if api == "" {
		log.Fatal("not found api key: GEMINI_API_KEY")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(api))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	var titleLst string
	for _, x := range c.Lst {
		titleLst += (x.Title + "\n")
	}
	model := client.GenerativeModel("gemini-2.5-flash-lite")
	prompt := fmt.Sprintf("Translate this list to Japanese. Please translate each line independently and output them separated by commas.\n%s", titleLst)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("-------------------\n%s\n\n\n\n-------------------\n", resp.Candidates[0].Content.Parts[0])

	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		res := fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0])
		resLst := strings.Split(res, ",")
		for i, x := range resLst {
			c.Lst[i].Title = x
		}
	} else {
		fmt.Printf("gemini not resp")
	}
}

func (c Content) String() string {
	return fmt.Sprintf("title: %s\n  link: %s\n  date: %s\n", c.Title, c.Link, c.Date)
}

func (c Contents) String() string {
	var res string
	for _, x := range c.Lst {
		res += fmt.Sprintf("%s\n", x.String())
	}
	return res
}

func main() {
	url := "https://hnrss.org/newest"
	fmt.Println("Fetching RSS feed from:", url)

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("http reqest error: %v\n", err)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("http status error: %s\n", resp.Status)
		return
	}
	fmt.Println("Successfully fetched RSS feed. Parsing feed...")

	fp := gofeed.NewParser()
	feed, err := fp.Parse(resp.Body)
	if err != nil {
		fmt.Printf("parse error of feed: %v\n", err)
		return
	}
	fmt.Println("Successfully parsed RSS feed.")

	fmt.Printf("## feed title: %s\n", feed.Title)

	var contentsSrc []Content
	for _, item := range feed.Items {
		var pubDate string

		if item.PublishedParsed != nil {
			pubDate = item.PublishedParsed.Format("2006-01-02 15:04:05")
		} else {
			pubDate = ""
		}

		contentsSrc = append(contentsSrc, Content{item.Title, item.Link, pubDate})
		// contentsSrc = append(contentsSrc, fmt.Sprintf("\ntitle: %s\n  link: %s\n  date: %s\n", item.Title, item.Link, pubDate))
	}
	contens := Contents{contentsSrc}

	fmt.Println("Translating content...")
	contens.Translate()
	fmt.Println("Translation complete.")
	fmt.Printf("%s", contens)

	u, err := json.Marshal(contens)
	if err != nil {
		fmt.Printf("convertion json error: %v", err)
		return
	}

	current, _ := os.Getwd()
	fileName := filepath.Join(current, "res.json")
	e := os.WriteFile(fileName, u, 0644)
	if e != nil {
		fmt.Printf("failed to write to file: %v", e)
	}
}
