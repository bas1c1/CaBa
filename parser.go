package main

import (
	"fmt"
	"os"
	"strconv"
)

func parseDbSlice(line string) dbslice {
	ds := dbslice{-1, "", ""}
	l := len(line)

	var buf string

	for i := 0; i < l; i++ {
		if line[i] == ';' {
			if ds.index == -1 {
				ind, err := strconv.Atoi(buf)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				ds.index = ind
				buf = ""
			} else {
				ds.key = buf
				buf = ""
			}
			continue
		}
		buf += string(line[i])
		if i == l-1 {
			ds.value = buf
			break
		}
	}

	return ds
}

func parseRequest(line string) request {
	req := request{"", []string{}}
	var buf string
	l := len(line)
	for i := 0; i < l; i++ {
		if line[i] == '{' {
			req.fn = buf
			i++
			buf = ""
			for ; line[i] != '}'; i++ {
				if line[i] == '"' {
					i++
					for ; line[i] != '"'; i++ {
						buf += string(line[i])
					}
					continue
				} else if line[i] == ',' {
					req.args = append(req.args, buf)
					buf = ""
					continue
				}
				buf += string(line[i])
			}
			req.args = append(req.args, buf)
			buf = ""
		}
		buf += string(line[i])
	}
	return req
}
