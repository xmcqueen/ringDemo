package main

import (
	"bufio"
	"container/ring"
	"flag"
	"fmt"
	"os"

	"github.com/tidwall/wal"
)

type FileReader struct {
	Validate func(os.FileInfo, error) error
	Name     string

	Value *bufio.Reader
}

func (fv *FileReader) Set(v string) error {
	_, err := os.Stat(v)
	if err != nil {
		return err
	}

	inFile, err := os.Open(v)
	if err != nil {
		return err
	}

	inputReader := bufio.NewReader(inFile)
	fv.Name = v
	fv.Value = inputReader

	return err
}

func (fv *FileReader) Get() *bufio.Reader {
	return fv.Value
}

func (fv *FileReader) String() string {
	return fv.Name
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	log, _ := wal.Open("mylog", nil)
	defer log.Close()

	defaultBufSize := 55
	var fr FileReader
	flag.Var(&fr, "f", "file to read from defaults to stdin")
	bufsize := flag.Int("n", defaultBufSize, "how many lines to keep in the buffer")

	flag.Parse()
	fmt.Println("bufsize", *bufsize)

	if *bufsize < 1 {
		fmt.Fprintf(os.Stderr, "bufsize must be positive:bufsize:%d:is not a valid bufsize\n", *bufsize)
		os.Exit(1)
	}

	r := ring.New(*bufsize)

	fileScanner := bufio.NewScanner(fr.Get())

	for fileScanner.Scan() {
		r.Value = fileScanner.Text()
		r = r.Next()
	}

	index := uint64(1)
	r.Do(func(p any) {
		if p != nil {
			fmt.Println(p.(string))
			log.Write(index, []byte(p.(string)))
			index++
		}
	})

}
