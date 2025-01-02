package main

type transaction struct {
	id      int64
	request request
}

func (tr transaction) execute() string {
	switch tr.request.fn {
	case "get":
		if len(tr.request.args) <= 0 {
			caba_err("ERR GET - NOT ENOUGH ARGS")
			return "not ok"
		}
		v := db.get(maindb, tr.request.args[0])
		if v != nil {
			return v.value
		} else {
			return ""
		}
	case "set":
		if len(tr.request.args) <= 0 {
			caba_err("ERR SET - NOT ENOUGH ARGS")
			return "not ok"
		}
		v := []dbslice{}
		for i := 0; i < len(tr.request.args); i += 2 {
			v = append(v, dbslice{tr.request.args[i], tr.request.args[i+1]})
		}
		defer db.set(maindb, v)
		return "ok - set"
	case "del":
		if len(tr.request.args) <= 0 {
			caba_err("ERR DEL - NOT ENOUGH ARGS")
			return "not ok"
		}
		db.remove(maindb, tr.request.args[0])
		return "ok - del"
	case "update":
		update_backup()
		return "ok - updated"
	case "load":
		load_backup()
		return "ok - loaded"
	case "loadfrom":
		if len(tr.request.args) <= 0 {
			caba_err("ERR LOADFROM - NOT ENOUGH ARGS")
			return "not ok"
		}
		load_backup_file(tr.request.args[0])
		return "ok - loaded from " + tr.request.args[0]
	case "save":
		return save_backup()
	}
	return "ok"
}
