package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	input := "ls -al"
	cmd := exec.Command("sh", "-c", input)
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(stdout))
}
