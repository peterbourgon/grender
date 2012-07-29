package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func RecursiveCopy(from, to string) error {
	cmd := exec.Command("cp", "-r", from+"/", to+"/")
	buf, err := cmd.CombinedOutput()
	if err != nil {
		s := strings.Replace(string(buf), "\n", " ", 0)
		return fmt.Errorf("%s (%s)", s, err)
	}
	return nil
}
