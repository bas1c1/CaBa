package main

type transaction struct {
	id      int64
	request request
}

func (tr transaction) execute() string {
	switch tr.request.fn {
	case "create_db":
		if len(tr.request.args) <= 0 {
			caba_err("ERR GET - NOT ENOUGH ARGS")
			return "not ok"
		}
		v := create_db(tr.request.args[0])
		if v != nil {
			return v.name
		}
		return "not ok"
	case "choose_db":
		if len(tr.request.args) <= 0 {
			caba_err("ERR GET - NOT ENOUGH ARGS")
			return "not ok"
		}
		maindb = db{tr.request.args[0]}
		cache_.clear()
		return tr.request.args[0]
	case "get":
		if len(tr.request.args) <= 0 {
			caba_err("ERR GET - NOT ENOUGH ARGS")
			return "not ok"
		}
		v := db.get(maindb, tr.request.args[0])
		if v != _zeroslice {
			return v.value
		}
		return ""
	case "multiget":
		if len(tr.request.args) <= 0 {
			caba_err("ERR GET - NOT ENOUGH ARGS")
			return "not ok"
		}

		v := db.multiget(maindb, tr.request.args)
		if v == nil {
			return "[]"
		}

		rv := "["
		for _, k := range v {
			rv += k.value + ","
		}

		rv = rv[:len(rv)-1] + "]"

		return rv
	case "list":
		v := db.list(maindb)
		if v == nil {
			return "[]"
		}

		rv := "["
		for _, k := range v {
			rv += "key: " + "\"" + k.key + "\"" + " - " + "value: " + "\"" + k.value + "\"" + ","
		}

		rv = rv[:len(rv)-1] + "]"

		return rv
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
	case "clearcache":
		cache_.clear()
		return "ok - cleared"
	case "updatecache":
		cache_.save_cache()
		return "ok - updated"
	case "loadcache":
		cache_.load_cache()
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
	case "asave":
		go save_backup_async()
		return "ok - async updating backup started"
	case "loadcfg":
		if len(tr.request.args) <= 0 {
			caba_err("ERR LOADFROM - NOT ENOUGH ARGS")
			return "not ok"
		}
		load_cfg(tr.request.args[0])
		return "ok - loaded config " + tr.request.args[0]
	case "zip":
		dname := ""
		if len(tr.request.args) <= 0 {
			dname = maindb.name
		} else {
			dname = tr.request.args[0]
		}
		createZip(dname+".zip", dname)
		return dname
	case "unzip":
		dname := ""
		if len(tr.request.args) <= 0 {
			dname = maindb.name
		} else {
			dname = tr.request.args[0]
		}
		unzip(dname+".zip", dname)
		return dname
	}
	return "ok"
}
