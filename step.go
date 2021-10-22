package raccoon

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

// 递归所有步骤(准备工作)
func recursionStep(wg *sync.WaitGroup, envs *Envs, currentSteps []Step) {
	// log.Println("[recu]currentSteps: ", currentSteps)
	t, nextSteps := Step{}, make([]Step, 0)
	if len(currentSteps) > 1 {
		t, nextSteps = currentSteps[0], currentSteps[1:]
	} else if len(currentSteps) > 0 {
		t = currentSteps[0]
	}
	// log.Println("             t:", t)
	// log.Println("     nextSteps:", nextSteps)
	if len(t.Webroot) > 0 {
		go webStep(wg, envs, t, nextSteps)
	} else if len(t.Exec) > 0 {
		go cmdStep(wg, envs, t, nextSteps)
	}
}

func printStepInfo(tag string, info *[]string, doExit bool) {
	log.Printf("==> %s :\n", tag)
	for _, v := range *info {
		log.Println(v)
	}
}

// 处理 web 服务 的步骤
func webStep(wg *sync.WaitGroup, envs *Envs, t Step, nextSteps []Step) {
	webroot, listener, err := getWebSource(t.Webroot, t.Webaddr)
	if err != nil {
		wg.Done()
		panic(err)
	}
	// 打印当前Web服务地址
	port := listener.Addr().(*net.TCPAddr).Port
	info := []string{fmt.Sprintf("Serving %s on HTTP port: %d", t.Webroot, port)}
	// 保存服务地址到环境变量列表
	envs.AddEnv(
		fmt.Sprintf("APP_%s_URL", strings.ToUpper(t.Tag)),
		fmt.Sprintf("http://127.0.0.1:%d", port),
	)
	wg.Done()
	// 递归剩下的步骤
	recursionStep(wg, envs, nextSteps)
	printStepInfo(t.Tag, &info, false)
	err = webServer(webroot, listener)
	if err != nil {
		panic(err)
	}
}

// 处理 cmd 服务 的步骤
func cmdStep(wg *sync.WaitGroup, envs *Envs, t Step, nextSteps []Step) {
	e := append(os.Environ(), t.Envs...)
	e = append(e, envs.GetEnvs()...)
	// 递归剩下的步骤
	info := []string{fmt.Sprintf("[step] envs: %s", envs.GetEnvs())}
	info = append(info, fmt.Sprintf("[ cmd] exec: %s", t.Exec))
	// info = append(info, fmt.Sprintf("[->->]steps: %s", nextSteps))
	wg.Done()
	recursionStep(wg, envs, nextSteps)
	printStepInfo(t.Tag, &info, false)
	err := doCmd(t, &e)
	if err != nil {
		panic(err)
	}
}
