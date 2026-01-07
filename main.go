package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"github.com/Bention99/gator/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Println("Error reading File:")
		fmt.Printf(" - %v\n", err)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please enter your Username: ")
	input, _ := reader.ReadString('\n')

	input = strings.TrimSpace(input)

	c.DBURL = "postgres://example"
	c.CurrentUserName = input

	e := config.SetUser(c)
	if e != nil {
		fmt.Println("Error setting Username:")
		fmt.Printf(" - %v\n", e)
	}

	fmt.Printf("Username set to: %v\n", input)

	c2, err := config.Read()
	fmt.Printf("DBURL: %v User: %v\n", c2.DBURL, c2.CurrentUserName)
}