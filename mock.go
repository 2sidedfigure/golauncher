package thunder

import (
	"log"
	"os"
)

type MockLauncher struct {
	logger *log.Logger
}

func NewMockLauncher() *MockLauncher {
	return &MockLauncher{
		logger: log.New(os.Stderr, "mock launcher", log.LstdFlags),
	}
}

func (ml *MockLauncher) Close() error {
	ml.logger.Println("CLOSE")
	return nil
}

func (ml *MockLauncher) LedOff() error {
	ml.logger.Println("LED OFF")
	return nil
}

func (ml *MockLauncher) LedOn() error {
	ml.logger.Println("LED ON")
	return nil
}

func (ml *MockLauncher) Down() error {
	ml.logger.Println("DOWN")
	return nil
}

func (ml *MockLauncher) Up() error {
	ml.logger.Println("UP")
	return nil
}

func (ml *MockLauncher) Left() error {
	ml.logger.Println("LEFT")
	return nil
}

func (ml *MockLauncher) Right() error {
	ml.logger.Println("RIGHT")
	return nil
}

func (ml *MockLauncher) Fire() error {
	ml.logger.Println("FIRE")
	return nil
}

func (ml *MockLauncher) Stop() error {
	ml.logger.Println("STOP")
	return nil
}
