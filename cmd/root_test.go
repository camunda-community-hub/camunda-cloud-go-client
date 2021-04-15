package cmd

import (
	"os"
	"os/exec"
	"testing"
)

func Test_checkEnvVars(t *testing.T) {
	if checkEnvVars("", "") != false {
		t.Error("Expected false in case env vars are not exported properly")
	}
}

func Test_Execute_envVarsNotPresent(t *testing.T) {
	if os.Getenv("BE_EXECUTE") == "1" {
		ClientId = ""
		ClientSecret = ""
		Execute()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=Test_Execute_envVarsNotPresent")
	cmd.Env = append(os.Environ(), "BE_EXECUTE=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
