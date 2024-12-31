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
	}
	return "ok"
}
