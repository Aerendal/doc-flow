package cli

import (
	"fmt"
	"os"
	"time"
)

type spinner struct {
	frames  []string
	stop    chan struct{}
	stopped chan struct{}
}

func newSpinner(label string) *spinner {
	s := &spinner{
		frames:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		stop:    make(chan struct{}),
		stopped: make(chan struct{}),
	}
	go func() {
		defer close(s.stopped)
		if !colorEnabled() {
			// No TTY — just print the label once and wait
			fmt.Fprintf(os.Stderr, "%s...\n", label)
			<-s.stop
			return
		}
		i := 0
		for {
			select {
			case <-s.stop:
				fmt.Fprintf(os.Stderr, "\r\033[K") // clear line
				return
			default:
				fmt.Fprintf(os.Stderr, "\r%s %s", cyan(s.frames[i%len(s.frames)]), label)
				i++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
	return s
}

func (s *spinner) Stop() {
	close(s.stop)
	<-s.stopped
}
