package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func geturl() {
	url := flag.String("p", "", "enter the subdomains file path")
	outputfile := flag.String("o", "", "path to save the output file")
	flag.Parse()

	if *url == "" {
		fmt.Println("error: invalid path")
		return
	}

	path, err := filepath.Abs(*url)
	if err != nil {
		fmt.Println("error getting absolute path:", err)
		return
	}

	x, err := os.Open(path)
	if err != nil {
		fmt.Println("error opening file:", err)
		return
	}
	defer x.Close()

	var output *os.File
	if *outputfile != "" {
		output, err = os.Create(*outputfile)
		if err != nil {
			fmt.Println("error creating output file:", err)
			return
		}
		defer output.Close()
	}

	c := bufio.NewScanner(x)
	for c.Scan() {
		line := c.Text()
		if !strings.HasPrefix(line, "https://") {
			line = "https://" + line
		}

		resp, err := http.Get(line)
		if err != nil {
			fmt.Println(err)

		}

		status := fmt.Sprintf("[+] %v : %v\n", line, resp.StatusCode)
		fmt.Print(status)

		if output != nil {
			_, err := output.WriteString(status)
			if err != nil {
				fmt.Println(err)
			}
		}

		resp.Body.Close()
	}

	if err := c.Err(); err != nil {
		fmt.Println("error reading file:", err)
	}
}

func main() {
	geturl()
}
