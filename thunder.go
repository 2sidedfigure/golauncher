// Package golauncher provides a means to control USB connected Dream Cheeky
// Thunder Launchers (http://dreamcheeky.com/thunder-missile-launcher).
package golauncher

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

// ThunderLauncher provides funcs to control a USB connected Thunder Launcher.
type ThunderLauncher struct {
	device *usb.Device
	ledOn  bool
}

// GetConnectedThunderLaunchers returns a slice of *ThunderLaunchers, each
// member of the slice corresponding to a connected Thunder Launcher.
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

// Close the USB connection to the Thunder Launcher.
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

// Turn off the Thunder Launcher's LED.
func (tl *ThunderLauncher) LedOff() error {
	err := tl.setLed(LED_OFF)

	if err != nil {
		tl.ledOn = false
	}

	return err
}

// Turn on the Thunder Launcher's LED.
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

// Down starts moving the Thunder Launcher down.
func (tl *ThunderLauncher) Down() error {
	return tl.do(DOWN)
}

// Up starts moving the Thunder Launcher up.
func (tl *ThunderLauncher) Up() error {
	return tl.do(UP)
}

// Left starts moving the Thunder Launcher left.
func (tl *ThunderLauncher) Left() error {
	return tl.do(LEFT)
}

// Right starts moving the Thunder Launcher right.
func (tl *ThunderLauncher) Right() error {
	return tl.do(RIGHT)
}

// Fire starts the process of firing the Thunder Launcher.
// BUG(ryan): Need to add the appropriate timing so Fire will perform a "complete" fire
func (tl *ThunderLauncher) Fire() error {
	return tl.do(FIRE)
}

// Stop ceases the last command sent to the Thunder Launcher. Only LedOff and
// LedOn don't require Stop to be called after their invocation.
func (tl *ThunderLauncher) Stop() error {
	return tl.do(STOP)
}
