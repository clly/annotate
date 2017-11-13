package main

import (
	"os"
	"os/exec"
	"fmt"
	"time"

	"gitlab.com/clly/annotate"
)

func main() {
	args := os.Args[1:]
	s, err := exec.LookPath(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find command %s: %s", args[0], err)
	}
	cmd := &exec.Cmd {
		Path: s,
		Args: args,
	}
	//o, err := cmd.CombinedOutput()
	//if err != nil {
	//	println(err)
	//}
	//fmt.Printf("%s\n", o)
	//test(cmd)
	err = annotate.Decorate(cmd, os.Stdout, os.Stdout, annotater)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failure %s", err)
	}
}

func annotater(s string) string {
	return fmt.Sprintf("%s %s\n", time.Now().Format("2006/01/02 15:04:05"), s)
}

func test(c *exec.Cmd) {
	err := c.Start()
	if err != nil {
		fmt.Println(err)
	}
	err = c.Wait()
	if err != nil {
		fmt.Println(err)
	}
}