package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Response struct {
	Programs []Program
}

type Program struct {
	Name    string
	URL     string
	Bounty  bool
	Swag    bool
	Domains []string
}

func main() {
	var max int
	var bounty, swag bool
	flag.IntVar(&max, "max", 50, "Max targets to dump(Use -1 to dump all)")
	flag.BoolVar(&bounty, "bounty", false, "Only Show Targets that offer bounty")
	flag.BoolVar(&swag, "swag", false, "Only Show Targets that offer Swag")

	flag.Parse()

	data := FetchData()

	targets := map[string]struct{}{}

	if !bounty && !swag {
		for _, v := range data.Programs {
			for _, domain := range v.Domains {
				targets[domain] = struct{}{}
			}
		}
	} else {
		for _, v := range data.Programs {
			if bounty {
				for _, domain := range v.Domains {
					targets[domain] = struct{}{}
				}
			} else if swag {
				for _, domain := range v.Domains {
					targets[domain] = struct{}{}
				}
			}
		}
	}

	rand.Seed(time.Now().UnixNano())
	count := 0

	for k := range targets {
		fmt.Println(k)
		count++
		if max != -1 && count == max {
			return
		}
	}

}

func FetchData() *Response {
	res, err := http.Get("https://github.com/projectdiscovery/public-bugbounty-programs/blob/main/chaos-bugbounty-list.json")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var buff bytes.Buffer

	// Get Code Block data
	doc.Find("tbody .blob-code").Each(func(i int, s *goquery.Selection) {
		buff.WriteString(s.Text() + "\n")
	})

	var p Response
	err = json.Unmarshal(buff.Bytes(), &p)
	if err != nil {
		log.Fatal(err)
	}

	return &p
}
