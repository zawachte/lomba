package grafana

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	staticConfigTemplate = `# config file version
apiVersion: 1

# list of datasources to insert/update depending
# whats available in the database
datasources:
  # <string, required> name of the datasource. Required
- name: Loki
  # <string, required> datasource type. Required
  type: loki
  # <string, required> access mode. direct or proxy. Required
  access: direct
  # <int> org id. will default to orgId 1 if not specified
  orgId: 1
  # <string> url
  url: http://%s:3100`
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
	tempDir := os.TempDir()
	configFilePath := path.Join(tempDir, "datasource.yml")

	lokiEndpoint := GetOutboundIPOrLocalhost()
	staticConfig := fmt.Sprintf(staticConfigTemplate, lokiEndpoint)

	err := ioutil.WriteFile(configFilePath, []byte(staticConfig), 0777)
	if err != nil {
		return err
	}

	cmdString := fmt.Sprintf("run -d -v %s:/etc/grafana/provisioning/datasources/datasource.yml --name=grafana -p 3000:3000 grafana/grafana", configFilePath)
	cmd := exec.Command("docker", strings.Split(cmdString, " ")...)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// Get preferred outbound ip of this machine
func GetOutboundIPOrLocalhost() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "localhost"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}
