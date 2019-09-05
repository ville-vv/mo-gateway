package health

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// HealthCheckor
type ServerCheckor struct {
	inputs   []Input
	outputs  []Output
	stop     chan int
	interval time.Duration
}

// NewCheckor
// interval is used to gather data
func NewServerCheckor(inputs []Input, outputs []Output, interval time.Duration) *ServerCheckor {
	if interval == 0 {
		interval = time.Second * 20
	}
	return &ServerCheckor{inputs: inputs, outputs: outputs, stop: make(chan int), interval: interval}
}

// Run is an starter which check servers
func (s *ServerCheckor) Run(ctx context.Context) error {
	go func() {
		err := s.run(ctx)
		if err != nil {
			fmt.Println("server checkor run error ", err)
		}
	}()
	return nil
}

// Run
func (s *ServerCheckor) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		select {
		case <-s.stop:
			cancel()
		}
	}()

	// 初始化插件
	if err := s.initPlugins(); err != nil {
		return err
	}

	if err := s.connectOutput(ctx); err != nil {
		return err
	}

	inputC := make(chan Metric, 100)

	if err := s.startServerInputs(ctx, inputC); err != nil {
		return err
	}

	wg.Add(1)
	go func(dst chan Metric) {
		wg.Done()
		if err := s.runInputs(ctx, dst); err != nil {
			fmt.Println("inputs run error ", err)
		}
		s.stopServerInput()
		close(inputC)
	}(inputC)

	outputC := inputC
	wg.Add(1)
	go func(src chan Metric) {
		wg.Done()
		s.runOutputs(src)
	}(outputC)

	s.closePlugins()
	return nil
}

func (s *ServerCheckor) Stop() {
	close(s.stop)
}

// initInputs
func (s *ServerCheckor) initPlugins() error {

	for _, v := range s.inputs {

		if input, ok := v.(Initiator); ok {
			if err := input.Init(); err != nil {
				return err
			}
		}
	}

	for _, v := range s.outputs {
		if output, ok := v.(Initiator); ok {
			if err := output.Init(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *ServerCheckor) connectOutput(ctx context.Context) error {
	for _, v := range s.outputs {
		if err := v.Connect(); err != nil {
			// 第一次链接不上等待  15秒再次尝试连接
			if err := SleepContext(ctx, time.Second*15); err != nil {
				return err
			}

			if err = v.Connect(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *ServerCheckor) runInputs(ctx context.Context, dst chan<- Metric) error {
	var wg sync.WaitGroup
	for _, input := range s.inputs {
		acc := NewAccumulator(dst)
		wg.Add(1)
		go func(iupt Input) {
			defer wg.Done()
			s.toGather(ctx, acc, iupt, s.interval)
		}(input)
	}
	wg.Wait()
	return nil
}

func (s *ServerCheckor) toGather(ctx context.Context, acc Accumulator, input Input, interval time.Duration) {
	panicRecover()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		if err := input.Gather(acc); err != nil {
			acc.AddError("heath server check", err)
		}
		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			return
		}
	}
}

func (s *ServerCheckor) startServerInputs(ctx context.Context, dst chan<- Metric) error {
	started := []ServerInput{}
	for _, input := range s.inputs {
		if st, ok := input.(ServerInput); ok {
			acc := NewAccumulator(dst)
			err := st.Start(acc)
			if err != nil {
				// 出现错误要把已经启动的给 关闭掉
				for _, si := range started {
					si.Stop()
				}
				return err
			}
			started = append(started, st)
		}
	}
	return nil
}

func (s *ServerCheckor) runOutputs(src <-chan Metric) {
	for metric := range src {
		for _, output := range s.outputs {
			output.Report(metric)
		}
	}
	return
}

func (s *ServerCheckor) closePlugins() {
	for _, p := range s.outputs {
		p.Close()
	}
}

func (s *ServerCheckor) stopServerInput() {
	for _, input := range s.inputs {
		if it, ok := input.(ServerInput); ok {
			it.Stop()
		}
	}
}

// SleepContext
func SleepContext(ctx context.Context, duration time.Duration) error {
	if duration == 0 {
		return nil
	}

	t := time.NewTimer(duration)
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		t.Stop()
		return ctx.Err()
	}
}

// panicRecover displays an error if an input panics.
func panicRecover() {
	if err := recover(); err != nil {
		trace := make([]byte, 2048)
		runtime.Stack(trace, true)
		fmt.Printf("panicRecover %s\n", string(trace))
	}
}
