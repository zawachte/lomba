package loki

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	staticConfig = `auth_enabled: false

server:
  http_listen_port: 3100
  grpc_listen_port: 9096

common:
  path_prefix: /tmp/loki
  storage:
    filesystem:
      chunks_directory: /tmp/loki/chunks
      rules_directory: /tmp/loki/rules
  replication_factor: 1
  ring:
    instance_addr: 127.0.0.1
    kvstore:
      store: inmemory

schema_config:
  configs:
    - from: 2020-10-24
      store: boltdb-shipper
      object_store: filesystem
      schema: v11
      index:
        prefix: index_
        period: 24h

limits_config:
  enforce_metric_name: false
  reject_old_samples: false
  reject_old_samples_max_age: 43800h

ruler:
  alertmanager_url: http://localhost:9093

# By default, Loki will send anonymous, but uniquely-identifiable usage and configuration
# analytics to Grafana Labs. These statistics are sent to https://stats.grafana.org/
#
# Statistics help us better understand how Loki is used, and they show us performance
# levels for most users. This helps us prioritize features and documentation.
# For more information on what's sent, look at
# https://github.com/grafana/loki/blob/main/pkg/usagestats/stats.go
# Refer to the buildReport method to see what goes into a report.
#
# If you would like to disable reporting, uncomment the following lines:
#analytics:
#  reporting_enabled: false
`
)

func BringDownPod() error {

	cmdString := "kill loki"
	cmd := exec.Command("docker", strings.Split(cmdString, " ")...)
	err := cmd.Run()
	if err != nil {
		return err
	}

	cmdString = "rm loki"
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

	configFilePath := path.Join(tempDir, "loki-config.yaml")

	err := ioutil.WriteFile(configFilePath, []byte(staticConfig), 0777)
	if err != nil {
		return err
	}

	cmdString := fmt.Sprintf("run --name loki -d -v %s:/test/loki-config.yaml -p 3100:3100 grafana/loki:2.6.0 -config.file=/test/loki-config.yaml", configFilePath)
	cmd := exec.Command("docker", strings.Split(cmdString, " ")...)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

/**
func BringUpPod(ctx context.Context) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	imageName := "grafana/loki:2.6.0"

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return "", err
	}

	defer out.Close()
	io.Copy(os.Stdout, out)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, nil, "")
	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	return resp.ID, nil
}
**/
