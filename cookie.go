package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/gookit/color"
)

func main() {
	// stay alive
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	color.Notice.Println("gocookie 1.0")
	color.Notice.Println(runtime.Version())
	color.Notice.Printf("Running on %d threads\n", runtime.NumCPU())

	// main logic
	x, y := getXY()
	color.Infof("X: %s Y: %s\n", x, y)
	go checkPos(x, y)
	// Spawn thread for every CPU
	for i := 0; i < runtime.NumCPU()-1; i++ {
		color.Infof("Launch click thread %d\n", i)
		go runClick()
	}

	<-sc
	os.Exit(0)
}

func checkPos(startx, starty string) {
	x, y := startx, starty
	for range time.Tick(time.Second) {
		x, y = getXY()
		if startx != x || starty != y {
			color.Infoln("Mouse moved")
			os.Exit(0)
		}
	}
}

func getXY() (x string, y string) {
	splt := strings.Split(strings.TrimSpace(runCursorPos()), " ")
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

func runCursorPos() string {
	cmd := exec.Command("xdotool", "getmouselocation")
	out, err := cmd.CombinedOutput()
	if err != nil && cmd.ProcessState.ExitCode() != 1 {
		return fmt.Sprintln(color.Error.Sprint(err))
	}
	return string(out)
}

func runClick() {
	for range time.Tick(time.Millisecond) {
		cmd := exec.Command("xdotool", "click", "1")
		cmd.Run()
	}
}
