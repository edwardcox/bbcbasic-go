package interpreter

import (
	"fmt"
	"strconv"
	"strings"
)

func (i *Interpreter) tryGosub(text string, nextPC int) (bool, error) {
	upper := strings.ToUpper(strings.TrimSpace(text))
	if !strings.HasPrefix(upper, "GOSUB") {
		return false, nil
	}

	rest := strings.TrimSpace(text[len("GOSUB"):])
	if rest == "" {
		return true, fmt.Errorf("missing line number for GOSUB")
	}

	target, err := strconv.Atoi(rest)
	if err != nil || target <= 0 {
		return true, fmt.Errorf("invalid GOSUB target: %s", rest)
	}

	i.runtime.PushReturn(nextPC)
	i.runtime.SetJump(target)
	return true, nil
}

func (i *Interpreter) tryReturn(text string) (bool, error) {
	upper := strings.ToUpper(strings.TrimSpace(text))
	if upper != "RETURN" {
		return false, nil
	}

	pc, ok := i.runtime.PopReturn()
	if !ok {
		return true, fmt.Errorf("RETURN without GOSUB")
	}

	i.runtime.SetReturnPC(pc)
	return true, nil
}
