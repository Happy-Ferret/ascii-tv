package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type frame struct {
	lines    []string
	duration time.Duration
}

func main() {
	moviePtr := flag.String("movie", "resources/sw1.txt", "path to ASCII movie file")
	addrPtr := flag.String("addr", ":8080", "TCP address to listen on")
	flag.Parse()
	data, err := ioutil.ReadFile(*moviePtr)
	if err != nil {
		fmt.Printf("Failed to load file %s\n", *moviePtr)
	}
	lines := strings.Split(string(data), "\n")
	frameHeight := 13
	var frames []frame
	for i := range lines {
		if i%(frameHeight+1) == 0 {
			frameDurationStr := lines[i]
			frameDurationInt, err := strconv.ParseInt(frameDurationStr, 0, 64)
			if err != nil {
				fmt.Printf("Failed to parse frame duration from line: %s", frameDurationStr)
			}
			frames = append(frames, frame{lines[i+1 : i+1+frameHeight], time.Duration(frameDurationInt)})
		}
	}
	fmt.Printf("Extracted %d frames from %s\n", len(frames), *moviePtr)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, frame := range frames {
			// Clear terminal and move cursor to position (1,1)
			fmt.Fprint(w, "\033[2J\033[1;1H")
			for _, line := range frame.lines {
				fmt.Fprintln(w, line)
			}
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			time.Sleep(frame.duration * time.Second / 15)
		}
	})

	fmt.Printf("Listening on %s\n", *addrPtr)
	log.Fatal(http.ListenAndServe(*addrPtr, nil))
}
