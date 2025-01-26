package main

type exception interface{}

type handle struct {
	try     func()
	catch   func(exception)
	finally func()
}

func throw(ex exception) {
	//naah
	//panic(ex)

	//well this function needs to be rewrited buuut i dont care bout this shit so use this function if you need it somehow

	caba_err(ex)
}

func (hndl handle) do_() {
	if hndl.finally != nil {
		defer hndl.finally()
	}

	if hndl.catch != nil {
		defer func() {
			if r := recover(); r != nil {
				hndl.catch(r)
			}
		}()
	}

	hndl.try()
}

func _check(err error) {
	if err != nil {
		caba_err(err)
	}
}
