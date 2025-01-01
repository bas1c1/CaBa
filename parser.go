package main

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
			buf += parseString(line, &i)
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
			buf = parseString(line, &i)
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

func parseString(line string, offset *int) string {
	buf := ""
	i := *offset + 1
	for ; line[i] != '"' && i < len(line)-1; i++ {
		if line[i] == '\\' {
			buf += "\\"
			i++
			buf += string(line[i])
			continue
		}
		buf += string(line[i])
	}
	*offset = i
	return buf
}
