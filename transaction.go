package main

type transaction struct {
	id      int64
	request request
}

// type operation func(transaction) interface{}
func (tr transaction) execute() string {
	switch tr.request.fn {
	case "get":
		v := db.get(maindb, tr.request.args[0]).value
		return v
	case "set":
		v := []kvpair{}
		for i := 0; i < len(tr.request.args); i += 2 {
			v = append(v, kvpair{tr.request.args[i], tr.request.args[i+1]})
		}
		db.set(maindb, v)
		return "ok - set"
	}
	return "ok"
}
