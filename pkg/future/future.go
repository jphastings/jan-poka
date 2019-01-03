package future

type Future chan Result

type Result struct {
	ok  bool
	Err error
}

func New() Future {
	return make(Future)
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

func (p Future) Bubble(r Result) {
	p <- r
}
