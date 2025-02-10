package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func replaceAppName(root, oldAppName, newAppName string) error {
	if fi, err := os.Stat(root); err != nil {
		return fmt.Errorf("unable to stat path '%s': %w", root, err)
	} else {
		if !fi.IsDir() {
			return replaceFileAppName(root, oldAppName, newAppName)
		}
	}

	fileRepl := map[string]string{}
	dirRepl := map[string]string{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(filepath.Base(path), oldAppName) {
			if d.IsDir() {
				dirRepl[path] = filepath.Join(filepath.Dir(path), strings.ReplaceAll(filepath.Base(path), oldAppName, newAppName))
			} else {
				fileRepl[path] = filepath.Join(filepath.Dir(path), strings.ReplaceAll(filepath.Base(path), oldAppName, newAppName))
			}
		}
		if d.IsDir() {
			return nil
		}
		return replaceFileAppName(path, oldAppName, newAppName)
	})
	if err != nil {
		return err
	}
	for curName, newName := range fileRepl {
		if err := os.Rename(curName, newName); err != nil {
			return err
		}
	}
	for curName, newName := range dirRepl {
		if err := os.Rename(curName, newName); err != nil {
			return err
		}
	}

	return nil
}

func replaceFileAppName(file, oldAppName, newAppName string) error {
	var (
		buf bytes.Buffer
	)
	err := func() error {
		input, err := os.Open(file)
		if err != nil {
			return err
		}
		defer func() {
			_ = input.Close()
		}()
		scanner := bufio.NewScanner(input)
		writeLine := func(line []byte) {
			if bytes.Contains(line, []byte(oldAppName)) {
				line = bytes.ReplaceAll(line, []byte(oldAppName), []byte(newAppName))
			}
			buf.Write(line)
			buf.WriteString("\n")
		}
		for scanner.Scan() {
			writeLine(scanner.Bytes())
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
