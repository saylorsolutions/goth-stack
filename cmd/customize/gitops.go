package main

import (
	"os/exec"
)

func gitRollback(git string) {
	_ = exec.Command(git, "reset", "--hard").Run()
	_ = exec.Command(git, "clean", "-fd").Run()
	_ = exec.Command(git, "clean", "-fd").Run()
	_ = exec.Command(git, "clean", "-fX").Run()
}
