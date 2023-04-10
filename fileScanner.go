package main

import (
	"bufio"
	"os"
)

type FileScanner struct {
	Validate func(os.FileInfo, error) error
	Name     string

	Value *bufio.Scanner
}

func (fv *FileScanner) Set(v string) error {
	_, err := os.Stat(v)
	if err != nil {
		return err
	}

	inFile, err := os.Open(v)
	if err != nil {
		return err
	}

	fv.Name = v
	inputReader := bufio.NewReader(inFile)
	fileScanner := bufio.NewScanner(inputReader)
	fv.Value = fileScanner

	return err
}

func (fv *FileScanner) Get() *bufio.Scanner {
	return fv.Value
}

func (fv *FileScanner) String() string {
	return fv.Name
}
