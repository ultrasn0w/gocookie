package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gookit/color"
)

func main() {
	retChan := make(chan string)
	color.Notice.Println("gocookie 1.0")
	color.Notice.Println(runtime.Version())

	// main logic
	x, y := getXY(retChan)
	println(x + y)
}

func getXY(retChan chan string) (x string, y string) {
	go runCmdOutput(retChan, "xdotool", "getmouselocation")
	posOut := strings.TrimSpace(<-retChan)
	println(posOut)
	splt := strings.Split(posOut, " ")
	if len(splt) == 1 {
		os.Exit(1)
	}
	xfound := false
	yfound := false
	for _, v := range splt {
		if strings.HasPrefix(v, "x:") {
			x = v[2:]
			xfound = true
			continue
		}
		if strings.HasPrefix(v, "y:") {
			y = v[2:]
			yfound = true
			continue
		}
		if xfound && yfound {
			break
		}
	}
	if !(xfound && yfound) {
		os.Exit(1)
	}
	return x, y
}

func runCmdOutput(retChan chan<- string, cmds string, args ...string) {
	cmd := exec.Command(cmds, args...)
	out, err := cmd.CombinedOutput()
	if err != nil && cmd.ProcessState.ExitCode() != 1 {
		retChan <- fmt.Sprintln(color.Error.Sprint(err))
		return
	}
	retChan <- string(out)
}

func runCmd(cmds string, args ...string) {
	cmd := exec.Command(cmds, args...)
	err := cmd.Run()
	if err != nil && cmd.ProcessState.ExitCode() != 1 {
		color.Errorln(err.Error())
	}
}
