package host

type FileMode int

const (
	FileModeRead FileMode = iota
	FileModeWrite
	FileModeUpdate
)

type Host interface {
	Init() error
	Reset() error

	WriteChar(b byte) error
	WriteString(s string) error

	ReadChar() (byte, error)
	ReadLine(prompt string) (string, error)

	ClearScreen() error
}
