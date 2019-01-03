package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter port to connect to: (larger than 4000 please):")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)

}