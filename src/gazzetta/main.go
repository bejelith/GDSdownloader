package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

import "flag"

var pattern = "http://archiviostorico.gazzettadelsud.it/gazzettasud/books/messina/%s/%smessina/images/pages/Page-%d.jpg"
var fromDate = flag.String("from", "2021-01-01", "from")
var toDate = flag.String("to", "2021-01-03", "to")
var printHelp = flag.Bool("help", false, "Print help")

func main() {
	flag.Parse()

	if *printHelp {
		flag.Usage()
		os.Exit(0)
	}

	var err error
	var from, to, now time.Time
	var client = http.Client{}
	from, err = time.Parse("2006-01-02", *fromDate)
	to, err = time.Parse("2006-01-02", *toDate)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for now = from; to.Sub(now) >= 0; now = now.Add(time.Hour * 24) {

		nowString := now.Format("20060102")
		year := now.Format("2006")

	pageLoop:
		for page := 0; ; page += 1 {
			notFoundOrEmpty := false
			url := fmt.Sprintf(pattern, year, nowString, page)
			resp, _ := client.Get(url)
			switch resp.StatusCode {
			case 404:
				fmt.Printf("Page %d of %s was not found\n", page, nowString)
				notFoundOrEmpty = true
			case 200, 201, 202:
				if resp.ContentLength <= 0 {
					fmt.Printf("Page %d of %s has no content\n", page, nowString)
					notFoundOrEmpty = true
				}
				reader := bufio.NewReader(resp.Body)
				if f, ferr := os.Create(nowString + "_" + strconv.Itoa(page) + ".jpg"); err == nil {
					writer := bufio.NewWriter(f)
					reader.WriteTo(writer) // I dont really care at this stages
					writer.Flush()
					f.Close()
				} else {
					fmt.Printf(ferr.Error())
				}
				resp.Body.Close()
			}
			if page > 1 && notFoundOrEmpty {
				break pageLoop
			}

		}

	}
}
