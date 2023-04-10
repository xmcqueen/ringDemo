package main

import (
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

	log, _ := wal.Open("mylog", nil)
	defer log.Close()

	defaultBufSize := 55
	var fs FileScanner
	flag.Var(&fs, "f", "file to read from defaults to stdin")
	bufsize := flag.Int("n", defaultBufSize, "how many lines to keep in the buffer")

	flag.Parse()
	fmt.Println("bufsize", *bufsize)

	if *bufsize < 1 {
		fmt.Fprintf(os.Stderr, "bufsize must be positive:bufsize:%d:is not a valid bufsize\n", *bufsize)
		os.Exit(1)
	}

	r := ring.New(*bufsize)

	for fs.Get().Scan() {
		r.Value = fs.Get().Text()
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
