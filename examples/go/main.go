package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func main() {
	stdin, _ := ioutil.ReadAll(os.Stdin)

	fmt.Println("* This is STDIN: ", string(stdin))
	fmt.Println("* This is STDOUT")

	time.Sleep(500 * time.Millisecond)

	fmt.Println("* This is STDOUT")

	time.Sleep(500 * time.Millisecond)

	fmt.Fprint(os.Stderr, "* This is STDERR\n")

	time.Sleep(500 * time.Millisecond)

	fmt.Println("* This is STDOUT")

	time.Sleep(500 * time.Millisecond)

	fmt.Fprint(os.Stderr, "* This is STDERR\n")

	return
}
