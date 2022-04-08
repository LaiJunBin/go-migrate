package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/laijunbin/go-migrate/cmd"
)

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "init" {
		cmd.Execute()
		return
	}

	if _, err := os.Stat("cmd/migrate/migrate.go"); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Error: please run 'go-migrate init'.")
		return
	}

	args := []string{"run", "cmd/migrate/migrate.go"}
	args = append(args, os.Args[1:]...)

	b, err := exec.Command("go", args...).Output()
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	fmt.Print(string(b))

}
