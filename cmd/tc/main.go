package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"

	"try/thunder"
)

var (
	ENTER = []byte{10, 0, 0}
	SPACE = []byte{32, 0, 0}

	UP    = []byte{27, 91, 65}
	DOWN  = []byte{27, 91, 66}
	RIGHT = []byte{27, 91, 67}
	LEFT  = []byte{27, 91, 68}

	ESC = []byte{27, 0, 0}
	q   = []byte{113, 0, 0}
)

type UnbufferedTTY struct {
	defaults string
}

func (ut *UnbufferedTTY) Start() error {
	def := exec.Command("/bin/stty", "-g")
	def.Stdin = os.Stdin
	out, err := def.Output()
	if err != nil {
		return err
	}

	ut.defaults = string(out)

	set := exec.Command("/bin/stty", "cbreak", "-echo")
	set.Stdin = os.Stdin
	return set.Run()
}

func (ut *UnbufferedTTY) End() error {
	end := exec.Command("/bin/stty", ut.defaults)
	end.Stdin = os.Stdin
	return end.Run()
}

func main() {
	tl, err := thunder.GetConnectedThunderLaunchers()
	if err != nil {
		fmt.Println("There was an error looking for connected launchers: ", err)
		return
	}
	if len(tl) < 1 {
		fmt.Println("No connected launchers found. Connect on and try again.")
		return
	}
	//for now just use the first launcher
	launcher := tl[0]
	for i := 1; i < len(tl); i++ {
		tl[i].Close()
	}
	launcher.LedOn()
	launcher.Stop()
	fmt.Println("Launcher on!")

	fmt.Println("Setting up raw TTY...")
	ut := UnbufferedTTY{}
	if err := ut.Start(); err != nil {
		fmt.Println("ERROR 1: ", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	exit := false
	input := make(chan []byte)

	go func(ch chan []byte) {
		for {
			b := make([]byte, 3)
			reader.Read(b)
			ch <- b
		}
	}(input)

	fmt.Println("Listening for input...")
	for {
		t := time.NewTicker(25 * time.Millisecond)

		select {
		case b := <-input:
			switch {
			case bytes.Equal(b, LEFT):
				launcher.Left()
			case bytes.Equal(b, RIGHT):
				launcher.Right()
			case bytes.Equal(b, UP):
				launcher.Up()
			case bytes.Equal(b, DOWN):
				launcher.Down()
			case bytes.Equal(b, ENTER), bytes.Equal(b, SPACE):
				launcher.Fire()
				time.Sleep(6 * time.Second)
				launcher.Stop()
			case bytes.Equal(b, ESC), bytes.Equal(b, q):
				exit = true
			}
		case <-t.C:
			launcher.Stop()
		}

		t.Stop()
		if exit {
			break
		}
	}

	launcher.Stop()
	launcher.LedOff()
	launcher.Close()
	ut.End()
}
