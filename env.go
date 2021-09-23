package main

import "fmt"

type Envs struct {
	env map[string]string
	c   chan string
}

func NewEnvs() *Envs {
	return &Envs{
		env: make(map[string]string),
		c:   make(chan string, 2),
	}
}
func (e *Envs) GetEnvs() []string {
	r := make([]string, 0)
	for k, v := range e.env {
		r = append(r, fmt.Sprintf("%s=%s", k, v))
	}
	return r
}

func (e *Envs) AddEnv(k string, v string) *Envs {
	e.env[k] = v
	e.c <- k
	return e
}

func (e *Envs) WaitEnv(name string) string {
	if val, found := e.env[name]; found {
		return val
	}
	return <-e.c
}
