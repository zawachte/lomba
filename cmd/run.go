package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/go-kit/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/zawachte/lomba/internal/runner"
	"github.com/zawachte/lomba/pkg/grafana"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type runOptions struct {
	kubeconfig string
}

var runOpts = &runOptions{}

var runCmd = &cobra.Command{
	Use:     "run",
	Short:   "run ",
	Long:    "run ",
	Example: "	lomba run",
	//Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRun()
	},
}

func init() {
	runCmd.Flags().StringVarP(&runOpts.kubeconfig, "kubeconfig", "k", "", "kubeconfig")
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

	err = rr.Run(context.Background())
	if err != nil {
		return err
	}

	grafanaEndpoint := grafana.GetOutboundIPOrLocalhost()

	fmt.Printf("Kubernetes logs are injested to Loki. Ready to query at http://%s:3000\n", grafanaEndpoint)

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
