package console

import "time"

// TODO refactor

type Progress struct {
	Message string
	Done    chan string
	Sync    chan struct{}
}

func NewProgress(message string) *Progress {
	return &Progress{
		Message: message,
		Done:    make(chan string),
		Sync:    make(chan struct{}),
	}
}

func (p *Progress) Complete(message string) {
	p.Done <- message
	<-p.Sync
}

func (p *Progress) displayAnsi(output *Console) {
	chars := []rune{'|', '/', '-', '\\'}
	idx := 0

	output.Print("\u001B[s") // save current cursor position
	for {
		output.Print("\u001B[u\u001B[0K") // restore cursor position and clear line

		select {
		case message := <-p.Done:
			output.Println(message)
			return
		default:
			output.Printf("[%s] %s ", string(chars[idx]), p.Message)
			idx = (idx + 1) % len(chars)

			time.Sleep(200 * time.Millisecond)
		}
	}
}

func (p *Progress) Display(output *Console) {
	if output.EnableColors {
		p.displayAnsi(output)
	} else {
		output.Printf("%s\n", p.Message)
		output.Printf("%s\n", <-p.Done)
	}

	close(p.Sync)
}
