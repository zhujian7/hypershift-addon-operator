package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	goflag "flag"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stolostron/hypershift-addon-operator/pkg/agent"
	"github.com/stolostron/hypershift-addon-operator/pkg/manager"
	"go.uber.org/zap"
	utilflag "k8s.io/component-base/cli/flag"
	"open-cluster-management.io/addon-framework/pkg/version"
)

const (
	// this should match managedclusteraddon cr's spec value in order to trigger the reconcile
	// for this addon.
	componentName = "hypershift-addon"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	pflag.CommandLine.SetNormalizeFunc(utilflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)

	var logger logr.Logger

	zapLog, err := zap.NewDevelopment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start up logger %v\n", err)
		os.Exit(1)
	}

	logger = zapr.NewLogger(zapLog)

	command := newCommand(logger)

	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func newCommand(logger logr.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hypershift-addon",
		Short: "hypershift addon for acm.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
			}
			os.Exit(1)
		},
	}

	if v := version.Get().String(); len(v) == 0 {
		cmd.Version = "<unknown>"
	} else {
		cmd.Version = v
	}

	cmd.AddCommand(manager.NewManagerCommand(componentName, logger.WithName("manager")))
	cmd.AddCommand(agent.NewAgentCommand(componentName, logger.WithName("agent")))

	return cmd
}
