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
	sgc := make(chan os.Signal, 1)
	signal.Notify(sgc, os.Interrupt, os.Kill, syscall.SIGQUIT)
	for _, v := range sel.serves {
		wg := sync.WaitGroup{}
		wg.Add(1)
		vlog.LogI("server [%s] is starting !", v.Name())
		go func(s IServer, group *sync.WaitGroup) {
			group.Done()
			if err := s.Start(args...); err != nil {
				panic(err)
			}
		}(v, &wg)
		wg.Wait()
		vlog.LogI("server [%s] start ok !", v.Name())
	}
	sg := <-sgc
	fmt.Println("mo-gateway exit ", sg)
	return
}
