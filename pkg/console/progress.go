package console

import (
	"sync"
	"time"
)

// TODO refactor

type Progress struct {
	message string

	done chan struct{}
	wg   sync.WaitGroup
}

func NewProgress(message string) *Progress {
	return &Progress{
		message: message,
		done:    make(chan struct{}),
		wg:      sync.WaitGroup{},
	}
}

func (p *Progress) Done() {
	close(p.done)
	p.wg.Wait()
}

func (p *Progress) displayAnsi(out Writer) {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	chars := []rune{'|', '/', '-', '\\'}
	idx := 0

	for {
		// return to line start, clear it, print frame
		out.Printf("\r\u001B[0K[%s] %s", string(chars[idx]), p.message)
		idx = (idx + 1) % len(chars)

		select {
		case <-ticker.C:
		case <-p.done:
			// return to line start and clear it
			out.Print("\r\u001B[0K")
			return
		}
	}
}

func (p *Progress) Display(out Writer) {
	p.wg.Add(1)
	defer p.wg.Done()

	if _, ok := out.(ansiWriter); ok {
		p.displayAnsi(out)
	} else {
		out.Printf("%s\n", p.message)
		<-p.done
	}
}
