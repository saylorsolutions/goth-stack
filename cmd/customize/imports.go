package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	singleLineImport     = regexp.MustCompile(`^\s*import\s+([^"\s]+\s*)?"([^"]+)"$`)
	multilineImportStart = regexp.MustCompile(`^\s*import\s*\(\s*$`)
	multilineImport      = regexp.MustCompile(`^(\s*)([^"\s]+\s*)?"([^"]+)"`)
	multilineImportEnd   = regexp.MustCompile(`^\s*\)\s*$`)
	replaceFileGlobs     = []string{"*.go", "*.templ"}
	ignoreFileGlobs      []string
)

func replaceDirImports(dir, target, replacement string) error {
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == dir {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		for _, ignoreGlob := range ignoreFileGlobs {
			match, err := filepath.Match(ignoreGlob, filepath.Base(path))
			if err != nil {
				return err
			}
			if match {
				return nil
			}
		}
		for _, matchGlob := range replaceFileGlobs {
			match, err := filepath.Match(matchGlob, filepath.Base(path))
			if err != nil {
				return err
			}
			if match {
				if err := replaceFileImports(path, target, replacement); err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}

func replaceFileImports(file, target, replacement string) error {
	var buf bytes.Buffer
	err := func() error {
		input, err := os.Open(file)
		if err != nil {
			return err
		}
		defer func() {
			_ = input.Close()
		}()
		scanner := bufio.NewScanner(input)
		inMultilineImport := false
		writeLine := func(line []byte) {
			buf.Write(line)
			buf.WriteString("\n")
		}
		for scanner.Scan() {
			line := scanner.Bytes()
			switch {
			case singleLineImport.Match(line):
				groups := singleLineImport.FindSubmatch(line)
				buf.WriteString(fmt.Sprintf("import %s\"%s\"\n", groups[1], []byte(strings.Replace(string(groups[2]), target, replacement, 1))))
			case multilineImportStart.Match(line):
				inMultilineImport = true
				writeLine(line)
			case inMultilineImport && multilineImport.Match(line):
				groups := multilineImport.FindSubmatch(line)
				buf.WriteString(fmt.Sprintf("%s%s\"%s\"\n", groups[1], groups[2], []byte(strings.Replace(string(groups[3]), target, replacement, 1))))
			case multilineImportEnd.Match(line):
				inMultilineImport = false
				writeLine(line)
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
