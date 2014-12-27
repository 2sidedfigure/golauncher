package thunder

type Launcher interface {
	Close() error
	LedOff() error
	LedOn() error
	Down() error
	Up() error
	Left() error
	Right() error
	Fire() error
	Stop() error
}
