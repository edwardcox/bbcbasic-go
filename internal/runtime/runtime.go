package runtime

type Runtime struct {
	stopped bool
}

func New() *Runtime {
	return &Runtime{}
}

func (r *Runtime) Stop() {
	r.stopped = true
}

func (r *Runtime) IsStopped() bool {
	return r.stopped
}

func (r *Runtime) Reset() {
	r.stopped = false
}
