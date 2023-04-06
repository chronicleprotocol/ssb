/*
 * SSB Tools
 *     Copyright (C) 2023 Chronicle Labs, Inc.
 *
 *     This program is free software: you can redistribute it and/or modify
 *     it under the terms of the GNU Affero General Public License as published
 *     by the Free Software Foundation, either version 3 of the License, or
 *     (at your option) any later version.
 *
 *     This program is distributed in the hope that it will be useful,
 *     but WITHOUT ANY WARRANTY; without even the implied warranty of
 *     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *     GNU Affero General Public License for more details.
 *
 *     You should have received a copy of the GNU Affero General Public License
 *     along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	log.Println(os.Args[1:])

	var scale int64
	flag.Int64Var(&scale, "s", 100, "scale of progress")
	flag.Parse()

	name := flag.Arg(0)
	info, err := os.Stat(name)
	if err != nil {
		log.Panic(err)
	}
	size := info.Size()

	f, err := os.Open(name)
	if err != nil {
		log.Panic(err)
	}

	data := bufio.NewScanner(f)
	data.Split(bufio.ScanBytes)

	tok := []byte(`{"key":"`)
	tl, ti, depth := len(tok), 0, 0

	scaleS := fmt.Sprintf("%d", scale)
	scaleS = "0/" + scaleS[2:]
	if scale == 100 {
		scaleS = "%"
	} else if scale == 1000 {
		scaleS = "‰"
	}
	inJSON := false
	var counter, bytesRead, progress, p int64

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
					if p = bytesRead * scale / size; p != progress {
						progress = p
						log.Println(name, "-", bytesRead, "bytes", progress, scaleS, "of", size, "objects:", counter)
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
