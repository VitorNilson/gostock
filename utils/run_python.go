package utils

import "os/exec"

func RefineResearch() {

	cmd := exec.Command("python3", "refine.py")

	err := cmd.Run()

	if err != nil {
		panic(err)
	}
}
