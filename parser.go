package main

import (
	"slices"
	"strings"
)

func parseDbSlice(line string) dbslice {
	ds := dbslice{"", ""}
	l := len(line)

	var buf string

	for i := 0; i < l-1; i++ {
		if line[i] == ';' {
			ds.key = buf
			buf = ""
			continue
		} else if line[i] == '"' {
			buf += parseString(line, &i, true)
			continue
		}

		buf += string(line[i])
	}

	ds.value = buf

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
			parseArgs(line, i, &req)
			buf = ""
		}
		buf += string(line[i])
	}
	return req
}

func parseArgs(line string, offset int, req *request) {
	buf := ""

	for i := offset; line[i] != '}'; i++ {
		if line[i] == '"' {
			buf = parseString(line, &i, true)
			continue
		} else if line[i] == ',' {
			req.args = append(req.args, buf)
			buf = ""
			continue
		}
		buf += string(line[i])

	}

	req.args = append(req.args, buf)
}

func parseString(line string, offset *int, slash bool) string {
	buf := ""
	i := *offset + 1
	for ; line[i] != '"' && i < len(line)-1; i++ {
		if line[i] == '\\' {
			if slash {
				buf += "\\"
			}
			i++
			buf += string(line[i])
			continue
		}
		buf += string(line[i])
	}
	*offset = i
	return buf
}

func parseWord(line string, offset *int) string {
	if slices.Contains(strings.Split("qwertryuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM", ""), string(line[0])) {
		tmp := ""
		for ; line[*offset] != '='; *offset++ {
			tmp += string(line[*offset])
		}
		return tmp
	}
	return ""
}

func parseConfig(line string) {
	for i := 0; i < len(line); i++ {
		if wrd := parseWord(line, &i); wrd != "" {
			if wrd == "PASSKEY" && line[i+1] == '"' {
				i++
				config_.passkey = []byte(parseString(line, &i, false))
				if len(config_.passkey) > 32 {
					config_.passkey = config_.passkey[:32]
				} else if diff := 32 - len(config_.passkey); diff > 0 {
					for i := 0; i < diff; i++ {
						config_.passkey = append(config_.passkey, byte(i))
					}
				}
			}
		}
	}
}
