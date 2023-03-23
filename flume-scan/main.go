package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
)

func main() {
	name := os.Args[1]
	info, err := os.Stat(name)
	if err != nil {
		log.Panic(err)
	}
	size := info.Size()

	f, err := os.ReadFile(name)
	if err != nil {
		log.Panic(err)
	}

	data := bufio.NewScanner(bytes.NewReader(f))
	data.Split(bufio.ScanBytes)

	tok := []byte(`{"key":"`)
	tl, ti, depth, counter := len(tok), 0, 0, 0

	inJSON := false
	var bytesRead, progress, p int64
	for data.Scan() {
		bytesRead++
		if inJSON {
			fmt.Print(data.Text())
			if data.Bytes()[0] == '{' {
				depth++
			} else if data.Bytes()[0] == '}' {
				depth--
				if depth == 0 {
					inJSON = false
					ti = 0
					fmt.Println()

					counter++
					if p = bytesRead * 100 / size; p != progress {
						progress = p
						log.Println(name, "-", bytesRead, "bytes", progress, "% of", size, "objects:", counter)
					}
				}
			}
			continue
		}
		if data.Bytes()[0] == tok[ti] {
			ti++
			if ti == tl {
				inJSON = true
				depth = 1
				fmt.Print(string(tok))
			}
		}
	}

	log.Printf("Found %d JSON objects", counter)
}
