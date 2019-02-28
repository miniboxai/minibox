package main

import "minibox.ai/pkg/cmd"

func main() {
	apiserver := cmd.NewApiServer()
	apiserver.Execute()
}
