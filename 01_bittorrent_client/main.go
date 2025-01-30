package main

import (
	"bittorrent-client/torrentfile"
	"log"
	"os"
)

func main() {
	inPath := os.Args[1]
	outPath := os.Args[2]

	tf, err := torrentfile.Open(inPath)

	if err != nil {
		log.Fatal(err)
	}

	if err := tf.DownloadToFile(outPath); err != nil {
		log.Fatal(err)
	}
}
