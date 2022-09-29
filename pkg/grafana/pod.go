package granfana

import (
	"os/exec"
)

func BringUpPod() error {
	cmdString := "docker run -d --name=grafana -p 3000:3000 grafana/grafana"
	cmd := exec.Command(cmdString)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
