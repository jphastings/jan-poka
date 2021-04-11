package future

type Future chan Result

type Result struct {
	ok  bool
	Err error
}

func New() Future {
	return make(Future)
}

func Exec(exec func() error) Future {
	f := New()
	go func() {
		if err := exec(); err != nil {
			f.Fail(err)
		} else {
			f.Succeed()
		}
	}()
	return f
}

func (r Result) IsOK() bool {
	return r.ok
}

func (p Future) Succeed() {
	p <- Result{ok: true, Err: nil}
}

func (p Future) Fail(err error) {
	p <- Result{ok: false, Err: err}
}
