package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

var saved = true

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s inputFile searchString\n", os.Args[0])
	}

	chunksize := 1000000

	f, err := os.Open(os.Args[1])

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if err != nil {
		log.Fatal(err)
	}
	r := bufio.NewReader(f)

	o, err := os.Create("output.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := o.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	w := bufio.NewWriter(o)

	sp := strings.Split(os.Args[2], "|")
	log.Printf("Searching for terms: %#v", sp)

	buf := make([]byte, chunksize)
	for i := 0; ; i++ {

		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if n == 0 {
			break
		}

		for _, s := range sp {
			if strings.Contains(string(buf[:n]), s) {
				log.Printf("Found a match for term %s.", s)
				if saved == false {
					saved = true
					log.Printf("Writing chunk to output file.")
					c, err := w.Write(buf[:n])
					if err != nil {
						log.Fatal(err)
					}

					_, err = w.Write([]byte("\n\n[[RawFileSearchMatchSeparator]]\n\n"))
					if err != nil {
						log.Fatal(err)
					}

					log.Printf("Writing %d bytes.\n", c)

					if err = w.Flush(); err != nil {
						log.Fatal(err)
					}
				}
			}
		}
		saved = false
		if i%1000 == 0 && i != 0 {
			log.Printf("Scanned %d MB on %s.", i, os.Args[1])
		}
	}
}
