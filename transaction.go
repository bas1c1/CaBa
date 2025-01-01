package main

type transaction struct {
	id      int64
	request request
}

func (tr transaction) execute() string {
	switch tr.request.fn {
	case "get":
		v := db.get(maindb, tr.request.args[0])
		if v != nil {
			return v.value
		} else {
			return ""
		}
	case "set":
		v := []kvpair{}
		for i := 0; i < len(tr.request.args); i += 2 {
			v = append(v, kvpair{tr.request.args[i], tr.request.args[i+1]})
		}
		defer db.set(maindb, v)
		return "ok - set"
	case "del":
		db.remove(maindb, tr.request.args[0])
		return "ok - del"
	}
	return "ok"
}
