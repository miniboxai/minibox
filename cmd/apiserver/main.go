package main

import "minibox.ai/minibox/pkg/cmd"

func main() {
	apiserver := cmd.NewApiServer()
	apiserver.Execute()
}
