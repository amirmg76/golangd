package main

import (
	"fmt"
	"path/filepath"
	"os"
	"os/exec"
	"io/ioutil"
	"strconv"
	"syscall"
	"time"
	"errors"
)
var ErrSudo error

var (
	bin string
	cmd string
)

const (
	varDir = "/var/golangd/"
	pidFile = "golangd.pid"
	initdFile = "/etc/init.d/golangd"
	outFile = "golangd.log"
	errFile = "golangd.err"
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

func writePid(pid int) (err error) {
	f, err := os.OpenFile(filepath.Join(varDir, pidFile), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = fmt.Fprintf(f, "%d", pid); err != nil {
		return err
	}
	return nil
}

func getPid() (pid int, err error) {
	b, err := ioutil.ReadFile(filepath.Join(varDir, pidFile))
	if err != nil {
		return 0, err
	}
	if pid, err = strconv.Atoi(string(b)); err != nil {
		return 0, fmt.Errorf("Invalid PID value: %s", string(b))
	}
	return pid, nil
}

func runApp() error {
	fmt.Println("RUN")
	for {
		time.Sleep(time.Second)
	}
	return nil
}

func installApp() error {
	_, err := os.Stat(initdFile)
	if err == nil {
		return errors.New("Already installed")
	}
	f, err := os.OpenFile(initdFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		if !os.IsPermission(err) {
			return err
		}
		return ErrSudo
	}
	defer f.Close()
	if _, err = fmt.Fprintf(f, "#!/bin/sh\n\"%s\" $1", bin); err != nil {
		return err
	}
	fmt.Println("Daemon", bin, "installed")
	return nil
}

func uninstallApp() error {
	_, err := os.Stat(initdFile)
	if err != nil && os.IsNotExist(err) {
		return fmt.Errorf("the file is not exist")
	}
	if err = os.Remove(initdFile); err != nil {
		if !os.IsPermission(err) {
			return err
		}
		return ErrSudo
	}
	fmt.Println("Daemon", bin, "uninstalled")
	return nil
}

func statusApp() error {
	var pid int
	defer func() {
		if pid == 0 {
			fmt.Println("status: not active")
			return
		}
		fmt.Println("status: active - pid", pid)
	}()
	pid, err := getPid()
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	p, err := os.FindProcess(pid)
	if err != nil {
		return nil
	}
	if err = p.Signal(syscall.Signal(0)); err != nil {
		fmt.Println(pid, "not found - removing PID file...")
		os.Remove(filepath.Join(varDir, pidFile))
		pid = 0
	}
	return nil
}

func startApp() error {
	const perm = os.O_CREATE|os.O_APPEND|os.O_WRONLY
	var err error
	if err := os.MkdirAll(varDir, 0755); err != nil {
		if !os.IsPermission(err) {
			return err
		}
		return ErrSudo
	}
	cmd := exec.Command(bin, "run")
	cmd.Stdout, err = os.OpenFile(filepath.Join(varDir, outFile), perm, 0644)
	if err != nil {
		return err
	}
	cmd.Stderr, err = os.OpenFile(filepath.Join(varDir, errFile), perm, 0644)
	if err != nil {
		return err
	}
	cmd.Dir = "/"
	if err = cmd.Start(); err != nil {
		return err
	}
	if err := writePid(cmd.Process.Pid); err != nil {
		if err := cmd.Process.Kill(); err != nil {
			fmt.Println("Cannot kill process", cmd.Process.Pid, err)
		} 
		return err
	}
	fmt.Println("Started with PID", cmd.Process.Pid)
	return nil

}

func stopApp() error {
	pid, err := getPid()
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("the process wasn't started")
		}
		return err
	}

	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	if err := p.Signal(os.Kill); err != nil {
		return err
	}

	if err := os.Remove(filepath.Join(varDir, pidFile)); err != nil {
		return err
	}

	fmt.Println("Stopped PID", pid)
	return nil
	
}

func helpApp() {
	fmt.Println("you should use one of these command")
	fmt.Println("install uninstall start stop run")
}

