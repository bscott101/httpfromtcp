package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const filePath = "messages.txt"

func main() {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Printf("Reading data from %s\n", filePath)
	fmt.Println("==================")

	for line := range getLinesChannel(file) {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		defer f.Close()
		buf := make([]byte, 8)
		line := ""

		for {
			n, err := f.Read(buf)
			if err != nil {
				if line != "" {
					ch <- line
				}
				if err == io.EOF {
					break
				}
				fmt.Printf("Error: %v", err)
				return
			}
			text := string(buf[:n])
			strSplit := strings.Split(text, "\n")
			for i := 0; i < len(strSplit)-1; i++ {
				ch <- fmt.Sprintf("%s%s", line, strSplit[i])
				line = ""
			}
			line += strSplit[len(strSplit)-1]
		}
	}()

	return ch
}
