package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

// 支持启动时显示构建日期和构建版本
// 需要通过命令 ` go build -ldflags "-X main.build=`git rev-parse HEAD`" ` 打包
var build = "not set"

func main() {
	fmt.Printf("Build: %s\n", build)
	steps := *getConf()

	wg := sync.WaitGroup{}
	// 最后一项客户端启动前等待所有准备服务完成
	wg.Add(len(steps) - 1)

	// 循环所有服务准备的步骤(去掉了最后一项客户端步骤)
	envs := NewEnvs()
	recursionStep(&wg, envs, steps[:len(steps)-1])
	wg.Wait()

	// 启动最后一个客户端访问服务的步骤设定
	lastStep := steps[len(steps)-1]
	if len(lastStep.Webroot) > 0 {
		panic(fmt.Errorf("config error: The web service should not be the last step"))
	} else {
		e := append(os.Environ(), lastStep.Envs...)
		log.Println("[last] envs.env:", envs.GetEnvs())
		e = append(e, envs.GetEnvs()...)
		err := doCmd(lastStep, &e)

		if err != nil {
			panic(err)
		}
	}
}
