package loki

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
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

func BringUpPod() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(homeDir, "loki-config.yaml"), []byte(staticConfig), 0644)
	if err != nil {
		return err
	}

	cmdString := fmt.Sprintf("docker run --name loki -d -v %s:/mnt/config -p 3100:3100 grafana/loki:2.6.0 -config.file=/mnt/config/loki-config.yaml", homeDir)

	cmd := exec.Command(cmdString)
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
