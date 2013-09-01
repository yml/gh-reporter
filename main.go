package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
)

const (
	EXIT_SUCCESS = 0
	EXIT_FAILURE = 1
)

var accessToken string
var exitStatus = EXIT_SUCCESS
var exitMu sync.Mutex

func setExitStatus(n int) {
	exitMu.Lock()
	if exitStatus < n {
		exitStatus = n
	}
	exitMu.Unlock()
}

func init() {
	accessToken = os.Getenv("GITHUB_TOKEN")
}

var cmds = NewCommands(
	"CLI to interact to the github API",
	cmdIssues,
)

func main() {
	fmt.Println("Initialising the cli")
	fmt.Println("AccessToken: ", accessToken)
	flag.Parse()
	args := flag.Args()
	fmt.Println("args: ", args)
	cmds.Parse(args)

	os.Exit(exitStatus)
}
