package main

import (
	"github.com/benodiwal/docker_ssh/pkg/env"
	"github.com/benodiwal/docker_ssh/pkg/ssh"
)

func main() {
	env.Load()
	ssh.Init()
}
