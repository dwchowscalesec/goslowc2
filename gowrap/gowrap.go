package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	//cmd := exec.Command("go", "run", "./helloworld/helloworld.go")
	cmd := exec.Command("go", "run", "./messagebox/messagebox.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
}

//https://medium.com/rungo/executing-shell-commands-script-files-and-executables-in-go-894814f1c0f7
