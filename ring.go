package main

import (
	"bufio"
	"container/ring"
	"flag"
	"fmt"
	"os"

	"github.com/tidwall/wal"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	log, err := wal.Open("mylog", nil)
	defer log.Close()

	inputReader := bufio.NewReader(os.Stdin)
	defaultInput := "-"
	defaultBufSize := 55
	nFlag := flag.String("f", defaultInput, "file to read from defaults to stdin")
	bufsize := flag.Int("n", defaultBufSize, "how many lines to keep in the buffer")
	var inFile *os.File
	inFile.Close()

	flag.Parse()
	fmt.Println("nflag", *nFlag)
	fmt.Println("bufsize", *bufsize)
	if *nFlag != defaultInput {
		fmt.Println("got a cli arg for input file", *nFlag)
		inFile, err = os.Open(*nFlag)
		check(err)

		inputReader = bufio.NewReader(inFile)
	}

	if *bufsize < 1 {
		fmt.Fprintf(os.Stderr, "bufsize must be positive:bufsize:%d:is not a valid bufsize\n", *bufsize)
		os.Exit(1)
	}

	r := ring.New(*bufsize)

	fileScanner := bufio.NewScanner(inputReader)

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
