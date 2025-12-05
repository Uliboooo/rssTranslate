# rssTranslate

current version: 0.2.0

`rssTranslate` is a simple web server that fetches an RSS feed, translates the titles to Japanese using the Gemini API, and displays them on a clean, readable webpage.

## Prerequisites

-   Go programming language installed.
-   A Gemini API key.

## Setup

1.  **Clone the repository:**
    ```bash
    git clone <repository-url>
    cd rssTranslate
    ```

2.  **Set your Gemini API Key:**
    The application requires a `GEMINI_API_KEY` environment variable to be set.
    ```bash
    export GEMINI_API_KEY="YOUR_API_KEY"
    ```

## How to Run

1.  **Start the server:**
    Run the following command in the project's root directory:
    ```bash
    go run main.go
    ```
    By default, the server will start on port `58877` and use the Hacker News "newest" RSS feed.

2.  **View the application:**
    Open the `index.html` file in your web browser to see the translated RSS feed.

## Configuration

You can customize the server's behavior using command-line flags:

-   `-p`: Set the port for the server.
    -   Example: `go run main.go -p 8080`
    -   Default: `58877`
-   `-url`: Specify the RSS feed URL to fetch.
    -   Example: `go run main.go -url "https://example.com/feed.xml"`
    -   Default: `https://hnrss.org/newest`
-   `-api`: Set the API URL.
    -   Default: `http://localhost`

## How it Works

The Go backend starts a web server with a `/api/rts` endpoint. When this endpoint is called, it fetches the specified RSS feed, extracts the titles, and sends them to the Gemini API for translation. The translated titles, along with their links and publication dates, are then returned as a JSON object.

The frontend consists of a simple `index.html` and `script.js`. The JavaScript makes a `fetch` request to the `/api/rts` endpoint and dynamically populates the page with the translated content.
