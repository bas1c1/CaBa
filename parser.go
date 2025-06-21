package main

import (
	"slices"
	"strconv"
	"strings"
)

func parseRequests(line string) []request {
	var requests []request

	req := request{"", []string{}}
	var buf string
	l := len(line)
	for i := 0; i < l; i++ {
		for ;i < l;i++ {
			if line[i] == '{' {
				req.fn = buf
				i++
				parseArgs(line, &i, &req)
				buf = ""
				break
			}
			buf += string(line[i])
		}
		if i >= l {
			break
		}
		requests = append(requests, req)
		req = request{"", []string{}}
		buf = ""
	}
	return requests
}

func parseArgs(line string, offset *int, req *request) {
	buf := ""

	i := *offset

	for ; line[i] != '}'; i++ {
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

	if buf != "" {
		req.args = append(req.args, buf)
	}

	*offset = i
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
				passkey := parseString(line, &i, false)
				if len(passkey) > 32 {
					passkey = passkey[:32]
				} else if len(passkey) < 32 {
					caba_err("Not valid passkey")
					return
				}
				config_.passkey = []byte(passkey)
			} else if wrd == "CACHE_SIZE" && line[i+1] == '"' {
				i++
				config_.cache_size, _ = strconv.Atoi(parseString(line, &i, false))
			} else if wrd == "HASH_KEYS" && line[i+1] == '"' {
				i++
				parsed_int, _ := strconv.Atoi(parseString(line, &i, false))
				config_.hash_keys = parsed_int != 0
			} else if wrd == "CACHING" && line[i+1] == '"' {
				i++
				parsed_int, _ := strconv.Atoi(parseString(line, &i, false))
				config_.caching = parsed_int != 0
			}
		}
	}
}
