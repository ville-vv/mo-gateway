package server

import (
	"github.com/ville-vv/vilgo/vlog"
	"sync"
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
	for _, v := range sel.serves {
		wg.Add(1)
		vlog.LogI("server [%s] is starting !", v.Name())
		go func(s IServer, group *sync.WaitGroup) {
			group.Done()
			if err := s.Start(args...); err != nil {
				panic(err)
			}
		}(v, &wg)
		vlog.LogI("server [%s] start ok !", v.Name())
	}
	wg.Wait()
	return
}
