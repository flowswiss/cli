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

	out.Print("\u001B[s") // save current cursor position
	for {
		out.Print("\u001B[u\u001B[0K") // restore cursor position and clear line
		out.Printf("[%s] %s\n", string(chars[idx]), p.message)
		idx = (idx + 1) % len(chars)

		select {
		case <-ticker.C:
		case <-p.done:
			out.Print("\u001B[u\u001B[0K") // restore cursor position and clear line
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
