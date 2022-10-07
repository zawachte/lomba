package runner

import (
	"bufio"
	"context"
	"fmt"

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

func (r *runner) Run(cancelCtx context.Context) error {
	err := loki.BringUpPod()
	if err != nil {
		return err
	}

	err = grafana.BringUpPod()
	if err != nil {
		return err
	}

	// get list of pods from all namespaces
	podList, err := r.cs.CoreV1().Pods("").List(cancelCtx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, pod := range podList.Items {
		go r.streamPodLogs(cancelCtx, r.cs, pod)

	}

	return nil
}

func (r *runner) loadLogsToLoki(logLine string, logParser parser.Parser, labels map[string]string) error {
	tm, labelset, err := logParser.Parse(logLine, labels)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		r.logger.Log("Skipping log due to invalid parse", "Error", err.Error())
		return err
	}
	r.lokiClient.PostLog(logLine, tm, labelset)

	return nil
}

// streamPodLogs will stream the pod logs and load the logs to loki with relevant
// labels, loglines and timestamp
func (r *runner) streamPodLogs(cancelCtx context.Context, cs kubernetes.Interface, pod corev1.Pod) error {
	for _, container := range pod.Spec.Containers {
		podLogOptions := &corev1.PodLogOptions{
			Follow:     true,
			Timestamps: true,
		}

		req := cs.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, podLogOptions)
		stream, err := req.Stream(cancelCtx)
		if err != nil {
			return fmt.Errorf("error in streaming logs: %s", err.Error())
		}

		reader := bufio.NewScanner(stream)
		reader.Split(bufio.ScanLines)
		defer stream.Close()

		for reader.Scan() {
			labels := make(map[string]string)
			labels["namespace"] = pod.Namespace
			labels["pod_name"] = pod.Name
			labels["container_name"] = container.Name

			logLine := reader.Text()

			err = r.loadLogsToLoki(logLine, parser.NewContainerParser(), labels)
			if err != nil {
				return fmt.Errorf("error loading logs to loki: %s", err.Error())
			}
		}
	}
	return nil
}
