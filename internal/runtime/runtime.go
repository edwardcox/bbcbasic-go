package runtime

type Runtime struct {
	stopped    bool
	vars       map[string]int
	jumpTarget int
	hasJump    bool
}

func New() *Runtime {
	return &Runtime{
		vars: make(map[string]int),
	}
}

func (r *Runtime) Stop() {
	r.stopped = true
}

func (r *Runtime) IsStopped() bool {
	return r.stopped
}

func (r *Runtime) Reset() {
	r.stopped = false
	r.vars = make(map[string]int)
	r.jumpTarget = 0
	r.hasJump = false
}

func (r *Runtime) SetVar(name string, value int) {
	r.vars[name] = value
}

func (r *Runtime) GetVar(name string) (int, bool) {
	value, ok := r.vars[name]
	return value, ok
}

func (r *Runtime) SetJump(target int) {
	r.jumpTarget = target
	r.hasJump = true
}

func (r *Runtime) ConsumeJump() (int, bool) {
	if !r.hasJump {
		return 0, false
	}
	target := r.jumpTarget
	r.jumpTarget = 0
	r.hasJump = false
	return target, true
}
