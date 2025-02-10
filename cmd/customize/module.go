package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
)

var moduleNamePattern = regexp.MustCompile(`^module\s+(.+)$`)

func changeModuleName(newModuleName string) error {
	var (
		buf  bytes.Buffer
		file = "go.mod"
	)
	err := func() error {
		input, err := os.Open("go.mod")
		if err != nil {
			return err
		}
		defer func() {
			_ = input.Close()
		}()
		scanner := bufio.NewScanner(input)
		writeLine := func(line []byte) {
			buf.Write(line)
			buf.WriteString("\n")
		}
		for scanner.Scan() {
			line := scanner.Bytes()
			switch {
			case moduleNamePattern.Match(line):
				buf.WriteString(fmt.Sprintf("module %s\n", newModuleName))
			default:
				writeLine(line)
			}
		}
		if err := scanner.Err(); err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		return fmt.Errorf("error processing file '%s': %w", file, err)
	}
	if err := os.WriteFile(file, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to rewrite file '%s': %w", file, err)
	}
	return nil
}
