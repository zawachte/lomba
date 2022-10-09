package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/zawachte/lomba/internal/runner"
	"github.com/zawachte/lomba/pkg/grafana"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type runOptions struct {
	kubeconfig     string
	streamDuration time.Duration
}

var runOpts = &runOptions{}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run ",
	Long:  "run ",
	Example: "	lomba run\n" +
		"	lomba run --kubeconfig <path to kubeconfig> --stream-duration 10m\n" +
		"	lomba run -k <path to kubeconfig> -s 15m30s",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRun()
	},
}

func init() {
	runCmd.Flags().StringVarP(&runOpts.kubeconfig, "kubeconfig", "k", os.Getenv("HOME")+"/.kube/config", "kubeconfig")
	runCmd.Flags().DurationVarP(&runOpts.streamDuration, "stream-duration", "s", 1*time.Hour, "time duration to "+
		"stream logs into loki")
	RootCmd.AddCommand(runCmd)
}

func runRun() error {

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	cs, err := createClientSet(runOpts.kubeconfig)
	if err != nil {
		return err
	}
	rr, err := runner.NewRunner(runner.RunnerParams{
		Logger:    logger,
		ClientSet: cs,
	})
	if err != nil {
		return err
	}

	cancelCtx, cancelFunc := context.WithTimeout(context.Background(), runOpts.streamDuration)
	err = rr.Run(cancelCtx)
	if err != nil {
		cancelFunc()
		return err
	}

	grafanaEndpoint := grafana.GetOutboundIPOrLocalhost()

	fmt.Printf("Kubernetes logs will be streamed to Loki for next %s minutes. Ready to query at http://%s:3000\n",
		runOpts.streamDuration, grafanaEndpoint)

	// sleep for duration as much as set in stream-duration flag to keep the goroutines active
	time.Sleep(runOpts.streamDuration)
	cancelFunc()
	return nil
}

// move to pkg
func createClientSet(kubeconfig string) (kubernetes.Interface, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	if kubeconfig != "" {
		rules.ExplicitPath = kubeconfig
	}

	config, err := rules.Load()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load Kubeconfig")
	}

	configOverrides := &clientcmd.ConfigOverrides{}
	restConfig, err := clientcmd.NewDefaultClientConfig(*config, configOverrides).ClientConfig()
	if err != nil {
		if strings.HasPrefix(err.Error(), "invalid configuration:") {
			return nil, errors.New(strings.Replace(err.Error(), "invalid configuration:", "invalid kubeconfig file:", 1))
		}
		return nil, err
	}

	restConfig.UserAgent = "lomba"
	restConfig.QPS = 20
	restConfig.Burst = 100

	cs, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	return cs, nil
}
