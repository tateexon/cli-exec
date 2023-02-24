package command

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

type CommandOptions struct {
	WaitForStdPipe bool // Should we wait for stdout and stderr to finish before returning
	StdOutHandler  func(string)
	StdErrHandler  func(string)
}

// ExecuteCommand will execute an os exec while capturing the output
func ExecuteCommand(command string,
	args []string,
	opts CommandOptions,
) error {
	cmd := exec.Command(command, args...) // #nosec G204
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	var wgOut sync.WaitGroup
	wgOut.Add(1)
	go func() {
		if opts.WaitForStdPipe {
			defer wgOut.Done()
		}
		HandleOutput(stdout, opts.StdOutHandler)
	}()

	var wgErr sync.WaitGroup
	wgErr.Add(1)
	go func() {
		if opts.WaitForStdPipe {
			defer wgErr.Done()
		}
		HandleOutput(stderr, opts.StdErrHandler)
	}()

	// Should we wait for stdout and stderr to complete
	// sometimes we need to not wait since a pipe doesn't get closed
	if opts.WaitForStdPipe {
		wgOut.Wait()
		wgErr.Wait()
	}

	return cmd.Wait()
}

// HandleOutput will handle the output of a command
func HandleOutput(stdPipe io.ReadCloser, outputHandlerFunc func(string)) {
	scanner := bufio.NewScanner(stdPipe)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		if outputHandlerFunc != nil {
			outputHandlerFunc(m)
		} else {
			DefaultHandler(m)
		}
	}
}

func DefaultHandler(message string) {
	fmt.Println(message)
}
