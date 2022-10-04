package main

import (
	"fmt"
	"filepath"
)
var ErrSudo error

var (
	bin string
	cmd string
)

func init() {
	p, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	bin = p
	if len(os.Args) != 1 {
		cmd = os.Args[1]
	}
	ErrSudo = fmt.Errorf("try `sudo %s %s`", bin, cmd)
}

func main() {
	var err error
	switch cmd {
	case "run":
		err = runApp()
	case "install":
		err = installApp()
	case "uninstall":
		err = uninstallApp()
	case "status":
		err = statusApp()
	case "start":
		err = startApp()
	case "stop":
		err = stopApp()
	default:
		helpApp()
	}
	if err != nil {
		fmt.Println(cmd, "error:", err)
	}
}

func runApp() error {
	return nil
}

func installApp() error {
	return nil
}

func uninstallApp() error {
	return nil
}

func statusApp() error {
	return nil
}

func startApp() error {
	return nil
}

func stopApp() error {
	return nil
}

func helpApp() {
	
}