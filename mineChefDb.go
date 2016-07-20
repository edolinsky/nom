package main

import (
	"fmt"
	"net/http"
	"strconv"
	"bytes"
	"regexp"
	"io/ioutil"
)

func mine(start, end int, titleRgx, addressRgx, regionRgx, telRgx, rUrlRgx, closedRgx *regexp.Regexp) {

	for idx := start; idx < end; idx++ {

		// Prep URL buffer with index
		var url bytes.Buffer = *bytes.NewBufferString("http://www.chefdb.com/pl/")
		url.WriteString(strconv.Itoa(idx))

		resp, err := http.Get(url.String())

		if err != nil {
			fmt.Println("Agh! An http error!" + err.Error())
			return
		}

		// Read body and close connection
		iso_8859_body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		// Convert iso8859-1 body to UTF-8
		body := iso8859_1_to_Utf8(iso_8859_body)

		var name, city, address, region, tel, rUrl []byte
		var closed bool
		var status string
		var title, adrSub, rgnSub, telSub, rUrlSub [][]byte

		// If body is available, match against individual regular expressions
		if len(body) > 0 {
			title = titleRgx.FindSubmatch(body)
			adrSub = addressRgx.FindSubmatch(body)
			rgnSub = regionRgx.FindSubmatch(body)
			telSub = telRgx.FindSubmatch(body)
			rUrlSub = rUrlRgx.FindSubmatch(body)
		}

		// Extract matches if available
		if len(adrSub) > 1 {
			address = adrSub[1]
		}
		if len(rgnSub) > 1 {
			region = rgnSub[1]
		}
		if len(telSub) > 1 {
			tel = telSub[1]
		}
		if len(rUrlSub) > 1 {
			rUrl = rUrlSub[1]
		}
		if len(title) > 1 {
			// compile regexes and find name, city within header title
			name = title[1]

			city = title[2]

			closed = closedRgx.Match(title[0])
			if closed {
				status = "closed"
			} else {
				status = "open"
			}

			fmt.Println(string(name), string(address), string(city),
				string(region), string(tel), string(rUrl), status)
		}



	}
	return
}

func iso8859_1_to_Utf8(iso8859_1_buffer []byte) string {

	// slice of runes, to copy from byte array
	buffer := make([]rune, len(iso8859_1_buffer))

	// Iterate over byte array, converting to UTF-8
	for i, b := range iso8859_1_buffer {
		buffer[i] = rune(b)
	}

	return string(buffer)
}

func main() {

	batch_size := 5
	titleRgx := regexp.MustCompile(`<title>(.+),\s+(.+?)\s+[\(:].+<\/title>`)
	closedRgx := regexp.MustCompile(`\((CLOSED)\)`)
	addressRgx := regexp.MustCompile(`street-address.+title="(.+)"`)
	regionRgx := regexp.MustCompile(`region.+title="(.+)["\s]>`)
	telRgx := regexp.MustCompile(`tel.+title="(.+)["\s]>`)
	rUrlRgx := regexp.MustCompile(`url.+title="(.+)["\s]>`)

	for i := 0; i < 3; i++ {
		go mine(batch_size * i, batch_size * (i+1), titleRgx, addressRgx, regionRgx, telRgx, rUrlRgx, closedRgx)
	}
	var input string
	fmt.Scanln(&input)
}
