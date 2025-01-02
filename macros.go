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
