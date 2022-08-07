package puzzle

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/microhod/adventofcode/internal/markdown"
)

const (
	domain  = "adventofcode.com"
	baseURL = "https://adventofcode.com"
)

var (
	lastRequest = time.Now()
)

type Puzzle struct {
	Name      string
	Readme    string
	TestInput string
	Input     string
}

type Client struct {
	httpClient        *http.Client
	markdownConverter *md.Converter
	token             string
}

func NewClient(token string) *Client {
	converter := md.NewConverter("adventofcode.com", true, &md.Options{
		CodeBlockStyle: "fenced",
	})

	converter.Use(markdown.AdventOfCode())
	
	return &Client{
		httpClient:        http.DefaultClient,
		markdownConverter: converter,
		token:             token,
	}
}

func (client *Client) Get(year, day int) (*Puzzle, error) {
	html, err := client.getHTML(year, day)
	if err != nil {
		return nil, err
	}

	input, err := client.getInput(year, day)
	if err != nil {
		return nil, err
	}

	return &Puzzle{
		Name: client.getName(html),
		Readme: client.getREADME(html),
		TestInput: client.getTestInput(html),
		Input: input,
	}, nil
}

func (client *Client) getHTML(year, day int) (*goquery.Selection, error) {
	path := fmt.Sprintf("%d/day/%d", year, day)
	
	resp, err := client.get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// the 'article' tags are the actual puzzle information
	return doc.Find("article"), nil
}

func (client *Client) getName(html *goquery.Selection) string {
	name := html.Find("h2").First().Text()
	// remove the '---' padded round the name
	name = strings.ReplaceAll(name, "-", "")

	// trim the 'Day <number>: ' from the start
	if parts := strings.Split(name, ":"); len(parts) > 1 {
		name = parts[1]
	}

	return strings.TrimSpace(name)
}

func (client *Client) getREADME(html *goquery.Selection) string {
	return client.markdownConverter.Convert(html)
}

func (client *Client) getTestInput(html *goquery.Selection) string {
	// guess at the first 'pre' tag
	input := html.Find("pre").First().Text()

	return strings.TrimSpace(input)
}

func (client *Client) getInput(year, day int) (string, error) {
	path := fmt.Sprintf("%d/day/%d/input", year, day)

	resp, err := client.get(path)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bytes), err
}

func (client *Client) get(path string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", baseURL, path)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got non-OK status for %s: %d", url, resp.StatusCode)
	}

	return resp, nil
}

func (client *Client) do(req *http.Request) (*http.Response, error) {
	// trottle requests
	// https://www.reddit.com/r/adventofcode/comments/3v64sb/aoc_is_fragile_please_be_gentle/
	for time.Since(lastRequest).Seconds() < 5 {
		time.Sleep(time.Second)
	}

	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: client.token,
	})
	return client.httpClient.Do(req)
}
