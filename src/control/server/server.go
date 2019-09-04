package server

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"vilgo/vlog"
)

type Initializer interface {
	Init(...interface{}) error
	UnInit()
}

func Init(args ...interface{}) error {
	for _, av := range args {
		switch t := av.(type) {
		case Initializer:
			err := t.Init()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type IServer interface {
	Name() string
	Start(...interface{}) error
	Stop()
}

type Serve struct {
	serves map[string]IServer
}

func NewServe(args ...IServer) *Serve {
	s := new(Serve)
	s.serves = make(map[string]IServer)
	for _, v := range args {
		s.serves[v.Name()] = v
	}
	return s
}

func (sel *Serve) Start(args ...interface{}) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(w *sync.WaitGroup) {
		wg.Done()
		for _, v := range sel.serves {
			vlog.LogI("server [%s] have starting !", v.Name())
			if err := v.Start(args...); err != nil {
				panic(err)
			}
			vlog.LogI("server [%s] have started !", v.Name())
		}
	}(&wg)
	sgc := make(chan os.Signal, 1)
	signal.Notify(sgc, os.Interrupt, os.Kill, syscall.SIGQUIT)
	sg := <-sgc
	fmt.Println("mo-gateway exit ", sg)
	return
}
