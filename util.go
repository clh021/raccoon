package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func GetProgramPath() string {
	ex, err := os.Executable()
	if err == nil {
		return filepath.Dir(ex)
	}

	exReal, err := filepath.EvalSymlinks(ex)
	if err != nil {
		panic(err)
	}
	// fmt.Println("exReal: ", exReal)
	return filepath.Dir(exReal)
}

func doCmd(t Step, envs *[]string) error {
	client := filepath.Join(GetProgramPath(), t.Exec)
	log.Println("[ cmd]:" + client)
	cmd := exec.Command(client)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGINT, //如果主进程退出，则将 SIGINT 发送给子进程
	}
	cmd.Env = *envs
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
