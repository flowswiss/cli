package output

import "time"

func (o *Output) ansiProgress(message string, done <-chan string, sync chan<- struct{}) {
	chars := []rune{'|', '/', '-', '\\'}
	idx := 0

	o.Print("\u001B[s") // save current cursor position
	for {
		o.Print("\u001B[u\u001B[0K") // restore cursor position and clear line

		select {
		case message = <-done:
			o.Println(message)
			return
		default:
			o.Printf("[%s] %s ", string(chars[idx]), message)
			idx = (idx + 1) % len(chars)

			time.Sleep(200 * time.Millisecond)
		}
	}
}

func (o *Output) Progress(message string, done <-chan string, sync chan<- struct{}) {
	if o.EnableColors {
		o.ansiProgress(message, done, sync)
	} else {
		o.Printf("%s\n", message)
		o.Printf("%s\n", <-done)
	}

	close(sync)
}
