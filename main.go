package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

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
	prompt := fmt.Sprintf("Translate this list to Japanese. Please translate each line independently and output them separated by '|||' and Only show translated string.\n%s", titleLst)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatal(err)
	}

	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		res := fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0])
		res = strings.TrimSpace(res)
		res = strings.TrimSuffix(res, "|||")
		resLst := strings.Split(res, "|||")
		for i, x := range resLst {
			c.Lst[i].Title = strings.TrimSpace(x)
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

func reJSON(url string) (*Contents, error) {
	// url := "https://hnrss.org/newest"

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	fp := gofeed.NewParser()
	feed, err := fp.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	var contentsSrc []Content
	for _, item := range feed.Items {
		var pubDate string

		if item.PublishedParsed != nil {
			pubDate = item.PublishedParsed.Format("2006-01-02 15:04:05")
		} else {
			pubDate = ""
		}

		contentsSrc = append(contentsSrc, Content{item.Title, item.Link, pubDate})
	}
	contens := Contents{contentsSrc}

	contens.Translate()
	// fmt.Printf("%s", contens)

	return &contens, nil
}

func RssHandleWithURL(url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := reJSON(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		jsonBytes, err := json.MarshalIndent(*data, "", "  ")
		if err != nil {
			// http error 500
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)

	}
}

var Version string

func main() {
	portShortPtr := flag.Int("p", 58877, "port")
	urlPtr := flag.String("url", "https://hnrss.org/newest", "rss url")
	apiURLPtr := flag.String("api", "localhost", "api url")
	cliModePrt := flag.Bool("cli", false, "enable cli-mode(without api server)")
	ver := flag.Bool("version", false, "show version")
	flag.Parse()

	if *ver {
		if Version == "" {
			Version = "dev"
		}
		fmt.Printf("version:: %s", Version)
		os.Exit(0)
	}

	if *cliModePrt {
		data, err := reJSON(*urlPtr)
		if err != nil {
			fmt.Printf("Error %v", err)
		}
		jsonB, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Printf("Error: %v", err)
			os.Exit(1)
		}
		fmt.Printf("%s", jsonB)
		os.Exit(0)
	}

	http.HandleFunc("/api/rts", RssHandleWithURL(*urlPtr))
	port := fmt.Sprintf("%s:%d", *apiURLPtr, *portShortPtr)
	fmt.Printf("Server starting on %s\n", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
