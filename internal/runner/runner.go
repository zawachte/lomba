package runner

import (
	"bufio"
	"bytes"
	"context"
	"io"

	"github.com/go-kit/log"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/zawachte/lomba/pkg/grafana"
	"github.com/zawachte/lomba/pkg/loki"

	"github.com/zawachte/lomba/pkg/parser"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Runner interface {
	Run(context.Context) error
}

type runner struct {
	lokiClient loki.Client
	cs         kubernetes.Interface
	logger     log.Logger
}

type RunnerParams struct {
	URI       string
	Logger    log.Logger
	ClientSet kubernetes.Interface
}

func NewRunner(params RunnerParams) (Runner, error) {
	lokiClient, err := loki.NewClient(loki.ClientParams{
		URI:    "http://localhost:3100/loki/api/v1/push",
		Logger: params.Logger,
	})
	if err != nil {
		return nil, err
	}

	return &runner{
		lokiClient: lokiClient,
		logger:     params.Logger,
		cs:         params.ClientSet}, nil
}

func (r *runner) Run(ctx context.Context) error {

	err := loki.BringUpPod()
	if err != nil {
		return err
	}

	err = grafana.BringUpPod()
	if err != nil {
		return err
	}

	namespaceList, err := r.cs.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, ns := range namespaceList.Items {
		podList, err := r.cs.CoreV1().Pods(ns.Name).List(ctx, metav1.ListOptions{})
		if err != nil {
			return err
		}

		for _, pod := range podList.Items {

			for _, container := range pod.Spec.Containers {
				req := r.cs.CoreV1().Pods(ns.Name).GetLogs(pod.Name, &corev1.PodLogOptions{
					Timestamps: true,
					Container:  container.Name,
				})

				podLogs, err := req.Stream(ctx)
				if err != nil {
					return err
				}
				defer podLogs.Close()

				buf := new(bytes.Buffer)

				_, err = io.Copy(buf, podLogs)
				if err != nil {
					return err
				}

				labels := make(map[string]string)
				labels["namespace"] = ns.Name
				labels["pod_name"] = pod.Name
				labels["container_name"] = container.Name

				err = r.loadLogsToLoki(buf, parser.NewContainerParser(), labels)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (r *runner) loadLogsToLoki(rawLogs *bytes.Buffer, logParser parser.Parser, labels map[string]string) error {

	scanner := bufio.NewScanner(rawLogs)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		log_line := scanner.Text()
		tm, labels, err := logParser.Parse(log_line, labels)
		if err != nil {
			r.logger.Log("Skipping log due to invalid parse", "Error", err.Error())
			continue
		}
		r.lokiClient.PostLog(log_line, tm, labels)
	}

	return nil
}
