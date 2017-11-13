package annotate

import (
	"io"
	"os/exec"
	"sync"
	"github.com/pkg/errors"
	"bufio"
	"fmt"
	"time"
)

const bufferSize = 1024*64*4

type Annotater func(s string) string

type AnnotateBytes func(b []byte) string

func Decorate(cmd *exec.Cmd, stderr, stdout io.Writer, a Annotater) error {
	// Create pipes for synchronous reading and writing
	// Look at using separate functions and byte buffers + cmd.Stdout() which may be easier/simpler
	r, w := io.Pipe()
	cmd.Stdout = w
	//er, err := cmd.StderrPipe()
	//if err != nil {
	//	return err
	//}

	// Create channels for moving between reading and writing
	stdoutChan := make(chan string)
	//stderrChan := make(chan string)

	// We have 4 readers/writers so we just need to wait on that
	var sg sync.WaitGroup
	sg.Add(2)

	// Create go routines for Reading/Writing stdout and stderr
	go ReadStd(&sg, r, stdoutChan)
	//go ReadStd(&sg, er, stderrChan)
	go WriteStd(&sg, stdout, stdoutChan, a)
	//go WriteStd(&sg, stderr, stderrChan, a)
	err := cmd.Start()
	if err != nil {
		return errors.Wrap(err, "Failed to start command")
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	w.Close()
	sg.Wait()
	return nil
}

func ReadStd(sg *sync.WaitGroup, reader io.Reader, c chan<- string) {
	scanner := bufio.NewScanner(reader)
	b := make([]byte, 0, bufferSize)
	scanner.Buffer(b, bufferSize)
	for scanner.Scan() {
		c <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	sg.Done()
	close(c)
}

func WriteStd(sg *sync.WaitGroup, writer io.Writer, c <-chan string, a Annotater) {
	var ok = true
	var s string
	for ok {
		select {
		case s, ok = <-c:
			if ok {
				fmt.Fprint(writer, a(s))
			}
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}
	sg.Done()
}

// asyncAnnotate is not implemented
func asyncAnnotate(cmd exec.Cmd, stderr, stdout io.Writer, a Annotater) <-chan struct{} {
	return nil
}