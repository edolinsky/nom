package main

import (
	"fmt"
	"net/http"
	"strconv"
	"bytes"
	"regexp"
	"io/ioutil"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

/*
Database set-up:

	CREATE TABLE `restaurants` (
	  `id` int(11) NOT NULL AUTO_INCREMENT,
	  `name` varchar(45) DEFAULT NULL,
	  `rating` double DEFAULT NULL,
	  `address` varchar(45) DEFAULT NULL,
	  `city` varchar(45) DEFAULT NULL,
	  `region` varchar(45) DEFAULT NULL,
	  `url` varchar(45) DEFAULT NULL,
	  `phone` varchar(45) DEFAULT NULL,
	  PRIMARY KEY (`id`),
	  UNIQUE KEY `id_UNIQUE` (`id`)
	) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;
 */

const (
	BATCH_SIZE = 5000		// number of chefDB location pages to query per thread
	NUM_THREADS = 8			// number of threads to run
	DB = "restaurants"		// name of database in which to record data
	USER = "Erik"			// database username
	PASSWORD = "123456"		// database user password
	HOST = "tcp(127.0.0.1:3306)"	// address and port of database

)

func mine(start, end int, titleRgx, addressRgx, regionRgx, telRgx, rUrlRgx, closedRgx *regexp.Regexp) {

	for idx := start; idx < end; idx++ {

		// Prep URL buffer with index
		var url bytes.Buffer = *bytes.NewBufferString("http://www.chefdb.com/pl/")
		url.WriteString(strconv.Itoa(idx))

		// Execute http request and handle any errors
		resp, err := http.Get(url.String())
		if err != nil {
			fmt.Println("Agh! An http error! " + err.Error())
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

		// Extract matches if available, assign NULL value if not
		if len(adrSub) > 1 {
			address = adrSub[1]
		} else {
			address = []byte("NULL")
		}
		if len(rgnSub) > 1 {
			region = rgnSub[1]
		} else {
			region = []byte("NULL")
		}
		if len(telSub) > 1 {
			tel = telSub[1]
		} else {
			tel = []byte("NULL")
		}
		if len(rUrlSub) > 1 {
			rUrl = rUrlSub[1]
		} else {
			rUrl = []byte("NULL")
		}
		if len(title) > 1 {

			// If page title is available, extract restaurant name and city
			name = title[1]
			city = title[2]

			// Set status according to whether we find the word "closed" in html body
			closed = closedRgx.Match(title[0])
			if closed {
				status = "closed"
			} else {
				status = "open"
			}

			fmt.Println(string(name), string(address), string(city),
				string(region), string(tel), string(rUrl), status)

			// If title is available, we have something to record in the database

			// Connect to database, exit if not possible
			db, err := sql.Open("mysql", USER + ":" + PASSWORD + "@" + HOST + "/" + DB)
			if err != nil {
				fmt.Println("DB connection error. Exiting thread. " + err.Error())
				return
			}

			// Prepare insert query, report on failure
			insert, err := db.Prepare(`INSERT IGNORE INTO ` + DB + `.restaurants
			(name, address, city, region, url, phone)
			VALUES("` + string(name) + `", "` + string(address) + `", "` + string(city) +
			`", "` + string(region) + `", "` + string(rUrl) + `", "` + string(tel) + `");`)
			if err != nil {
				fmt.Println("MySQL insert error. Aborting entry. " + err.Error())
			}

			// Execute insert query, report on failure
			_, err = insert.Exec()
			if err != nil {
				fmt.Println("MySQL execution error. Aborting entry. " + err.Error())
			}

			insert.Close()
			db.Close()
		}
	}

	return
}

func iso8859_1_to_Utf8(iso8859_1_buffer []byte) []byte {

	// slice of runes, to copy from byte array
	buffer := make([]rune, len(iso8859_1_buffer))

	// Iterate over byte array, converting to UTF-8
	for i, b := range iso8859_1_buffer {
		buffer[i] = rune(b)
	}

	return []byte(string(buffer))
}

func main() {

	// Compile required regular expressions
	titleRgx := regexp.MustCompile(`<title>(.+),\s+(.+?)\s+[\(:].+<\/title>`)
	closedRgx := regexp.MustCompile(`\((CLOSED)\)`)
	addressRgx := regexp.MustCompile(`street-address.+title=""?(.+[^"])"`)
	regionRgx := regexp.MustCompile(`region.+title="(.+)["\s]>`)
	telRgx := regexp.MustCompile(`tel.+title="(.+)["\s]>`)
	rUrlRgx := regexp.MustCompile(`url.+title="(.+)["\s]>`)

	// Pass a batch off to each thread
	for i := 0; i < NUM_THREADS; i++ {
		go mine(BATCH_SIZE * i, BATCH_SIZE * (i+1),
			titleRgx, addressRgx, regionRgx, telRgx, rUrlRgx, closedRgx)
	}

	// Collect output
	var input string
	fmt.Scanln(&input)
}
