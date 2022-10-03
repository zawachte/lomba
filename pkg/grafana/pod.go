package grafana

import (
	"os/exec"
	"strings"
)

func BringDownPod() error {
	cmdString := "kill grafana"
	cmd := exec.Command("docker", strings.Split(cmdString, " ")...)
	err := cmd.Run()
	if err != nil {
		return err
	}

	cmdString = "rm grafana"
	cmd = exec.Command("docker", strings.Split(cmdString, " ")...)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func BringUpPod() error {

	BringDownPod()
	cmdString := "run -d --name=grafana -p 3000:3000 grafana/grafana"
	cmd := exec.Command("docker", strings.Split(cmdString, " ")...)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
