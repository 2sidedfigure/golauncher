package thunder

import (
	"fmt"

	"github.com/kylelemons/gousb/usb"
)

const (
	DOWN = 1 << iota
	UP
	LEFT
	RIGHT
	FIRE
	STOP
)

const (
	LED_OFF = iota
	LED_ON
)

type ThunderLauncher struct {
	device *usb.Device
	ledOn  bool
}

func GetConnectedThunderLaunchers() ([]*ThunderLauncher, error) {
	ctx := usb.NewContext()
	//defer ctx.Close()

	devices, err := ctx.ListDevices(func(d *usb.Descriptor) bool {
		return d.Vendor == 0x2123 && d.Product == 0x1010
	})
	if err != nil {
		return nil, err
	}
	if len(devices) == 0 {
		return nil, fmt.Errorf("No connected Thunder Launcher devices found")
	}

	tl := make([]*ThunderLauncher, 0)
	for _, d := range devices {
		tl = append(tl, newThunderLauncher(d))
	}

	return tl, nil
}

func newThunderLauncher(device *usb.Device) *ThunderLauncher {
	tl := &ThunderLauncher{device: device}

	tl.LedOff()

	return tl
}

func (tl *ThunderLauncher) Close() error {
	return tl.device.Close()
}

func (tl *ThunderLauncher) control(msg []byte) error {
	_, err := tl.device.Control(0x21, 0x09, 0, 0, msg)
	return err
}

func (tl *ThunderLauncher) setLed(state byte) error {
	return tl.control([]byte{3, state})
}

func (tl *ThunderLauncher) LedOff() error {
	err := tl.setLed(LED_OFF)

	if err != nil {
		tl.ledOn = false
	}

	return err
}

func (tl *ThunderLauncher) LedOn() error {
	err := tl.setLed(LED_ON)

	if err != nil {
		tl.ledOn = true
	}

	return err
}

func (tl *ThunderLauncher) do(action byte) error {
	return tl.control([]byte{2, action})
}

func (tl *ThunderLauncher) Down() error {
	return tl.do(DOWN)
}

func (tl *ThunderLauncher) Up() error {
	return tl.do(UP)
}

func (tl *ThunderLauncher) Left() error {
	return tl.do(LEFT)
}

func (tl *ThunderLauncher) Right() error {
	return tl.do(RIGHT)
}

func (tl *ThunderLauncher) Fire() error {
	return tl.do(FIRE)
}

func (tl *ThunderLauncher) Stop() error {
	return tl.do(STOP)
}
