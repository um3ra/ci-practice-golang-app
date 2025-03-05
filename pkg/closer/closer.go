package closer

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

var globalCloser = NewCloser()

func Add(funcs ...func() error) {
	globalCloser.add(funcs...)
}

func Wait() {
	globalCloser.wait()
}

func CloseAll() {
	globalCloser.closeAll()
}

func NewCloser(sig ...os.Signal) *Closer {
	cl := &Closer{done: make(chan struct{})}
	ch := make(chan os.Signal, 1)

	if len(sig) > 0 {
		go func() {
			signal.Notify(ch, sig...)
			<-ch
			signal.Stop(ch)
			cl.closeAll()
		}()
	}
	return cl
}

type Closer struct {
	once  sync.Once
	mu    sync.Mutex
	funcs []func() error
	done  chan struct{}
}

func (cl *Closer) wait() {
	<-cl.done
}

func (cl *Closer) add(funcs ...func() error) {
	cl.mu.Lock()
	cl.funcs = append(cl.funcs, funcs...)
	cl.mu.Unlock()
}

func (cl *Closer) closeAll() {
	cl.once.Do(func() {
		cl.mu.Lock()
		funcs := cl.funcs
		cl.funcs = nil
		cl.mu.Unlock()

		errs := make(chan error, len(funcs))
		for _, f := range funcs {
			go func(f func() error) {
				errs <- f()
			}(f)
		}
		for i := 0; i < cap(errs); i++ {
			if err := <-errs; err != nil {
				log.Println("error returned from Closer")
			}
		}
	})
}
