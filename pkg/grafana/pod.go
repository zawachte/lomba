package grafana

import (
	"os/exec"
	"strings"
)

func BringUpPod() error {
	cmdString := "kill grafana"
	cmd := exec.Command("docker", strings.Split(cmdString, " ")...)
	cmd.Run()

	cmdString = "rm grafana"
	cmd = exec.Command("docker", strings.Split(cmdString, " ")...)
	cmd.Run()

	cmdString = "run -d --name=grafana -p 3000:3000 grafana/grafana"
	cmd = exec.Command("docker", strings.Split(cmdString, " ")...)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
