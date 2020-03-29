package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

// Routine Pattern - Generator
// <-chan - read only channel
func titles(urls ...string) <-chan string {
	channel := make(chan string)

	for _, url := range urls {
		go func(url string) {
			resp, _ := http.Get(url)
			html, _ := ioutil.ReadAll(resp.Body)

			r, _ := regexp.Compile("<title>(.*?)<\\/title>")

			channel <- r.FindStringSubmatch(string(html))[1]

		}(url)
	}

	return channel
}

// Routine Pattern - Multiplexer
func multiplexer(firstEntry, secondEntry <-chan string) <-chan string {
	channel := make(chan string)

	go foward(firstEntry, channel)
	go foward(secondEntry, channel)

	return channel
}

func foward(origin <-chan string, destiny chan string) {
	for {
		destiny <- <-origin
	}
}

// Routine Pattern - Select - Control structure
func fastest(urlA, urlB, urlC string) string {
	chA, chB, chC := titles(urlA), titles(urlB), titles(urlC)

	select {
	case titleA := <-chA:
		return titleA

	case titleB := <-chB:
		return titleB

	case titleC := <-chC:
		return titleC

	case <-time.After(3500 * time.Millisecond):
		return "Timeout"

		// default:
		// 	return "No response yet"
	}
}

func main() {
	google := "https://google.com"
	twitch := "https://twitch.com"
	microsoft := "https://microsoft.com"
	youtube := "https://youtube.com"
	amazon := "https://amazon.com"

	// Multiplex and Generator patterns
	uniqueChannel := multiplexer(
		titles(google, twitch),
		titles(microsoft, youtube),
	)

	fmt.Println("First:", <-uniqueChannel, "|", <-uniqueChannel)
	fmt.Println("Second:", <-uniqueChannel, "|", <-uniqueChannel)

	// Select : control strucuture pattern
	faster := fastest(
		microsoft,
		amazon,
		twitch,
	)

	fmt.Println(faster)
}
