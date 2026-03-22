package runtime

type Runtime struct {
	stopped      bool
	vars         map[string]int
	jumpTarget   int
	hasJump      bool
	returnStack  []int
	hasReturnPC  bool
	returnTarget int
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
	r.returnStack = nil
	r.hasReturnPC = false
	r.returnTarget = 0
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

func (r *Runtime) PushReturn(pc int) {
	r.returnStack = append(r.returnStack, pc)
}

func (r *Runtime) PopReturn() (int, bool) {
	if len(r.returnStack) == 0 {
		return 0, false
	}

	last := len(r.returnStack) - 1
	pc := r.returnStack[last]
	r.returnStack = r.returnStack[:last]
	return pc, true
}

func (r *Runtime) SetReturnPC(pc int) {
	r.returnTarget = pc
	r.hasReturnPC = true
}

func (r *Runtime) ConsumeReturnPC() (int, bool) {
	if !r.hasReturnPC {
		return 0, false
	}

	pc := r.returnTarget
	r.returnTarget = 0
	r.hasReturnPC = false
	return pc, true
}
