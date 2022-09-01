package closer

import (
	"fmt"
	"sync"
)

type ICloser interface {
	Add(...func() error)
	CloseAll()
}

type closer struct {
	funcs []func() error
	mutex *sync.Mutex
}

func NewCloser() ICloser {
	return &closer{
		mutex: &sync.Mutex{},
	}
}

func (cl *closer) Add(functions ...func() error) {
	cl.mutex.Lock()

	cl.funcs = append(cl.funcs, functions...)

	cl.mutex.Unlock()
}

func (cl *closer) CloseAll() {
	cl.mutex.Lock()
	funcs := cl.funcs
	cl.funcs = nil
	cl.mutex.Unlock()

	errCh := make(chan error, len(funcs))
	defer close(errCh)

	for _, f := range cl.funcs {
		go func(f func() error) {
			errCh <- f()
		}(f)
	}

	for err := range errCh { // TO DO, подумать, что тут можно сделать. Возможно, приложение не выйдет из цикла
		if err != nil {
		}
		fmt.Println(err)
	}

}
