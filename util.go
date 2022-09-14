package raccoon

import (
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
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

func doCmdUser(t Step, envs *[]string) error {
	client := filepath.Join(GetProgramPath(), t.Exec)
	log.Println("[ cmd]:" + client)
	cmd := exec.Command(client)
	// 发送信号关闭子进程
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGINT, //如果主进程退出，则将 SIGINT 发送给子进程
	}
	// 设置具体运行用户
	user, err := user.Lookup("nobody")
	if err == nil {
		log.Printf("uid=%s,gid=%s", user.Uid, user.Gid)

		uid, _ := strconv.Atoi(user.Uid)
		gid, _ := strconv.Atoi(user.Gid)

		cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}
	}
	cmd.Env = *envs
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
