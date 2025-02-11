package main

//well i dont actually need this entire script, but it can be useful, so i wont delete it mkay?

type exception interface{}

type handle struct {
	try     func()
	catch   func(exception)
	finally func()
}

func throw(ex exception) {
	caba_err(ex)
	panic(-1)
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
