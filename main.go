package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("git", "clone", "--progress", "git@github.com:a-hilaly/blockchain", "ag2")
	buff := &bytes.Buffer{}
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = buff
	cmd.Run()
	cmd.Wait()
	fmt.Println("---", buff.String())
}
