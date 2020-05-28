/*
Copyright 2018 The OpenEBS Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package command

import (
	goflag "flag"

	"github.com/openebs/node-disk-manager/cmd/ndm_daemonset/controller"
	"github.com/openebs/node-disk-manager/pkg/features"
	"github.com/openebs/node-disk-manager/pkg/util"
	"github.com/openebs/node-disk-manager/pkg/version"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog"
)

// options is the options with which the daemon is to be run
var options controller.NDMOptions

// NewNodeDiskManager creates a new ndm.
func NewNodeDiskManager() (*cobra.Command, error) {

	// Create a new command
	cmd := &cobra.Command{
		Use:   "ndm",
		Short: "ndm controls the Node-Disk-Manager ",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			util.CheckErr(RunNodeDiskManager(cmd), util.Fatal)

			// set the feature gates on NDM daemon
			err := features.FeatureGates.SetFeatureFlag(options.FeatureGate)
			if err != nil {
				klog.Fatalf("error setting feature gate: %v", err)
			}
		},
	}

	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	cmd.PersistentFlags().StringVar(&options.ConfigFilePath, "config",
		controller.DefaultConfigFilePath,
		"Path to config file")
	cmd.PersistentFlags().StringSliceVar(&options.FeatureGate, "feature-gates",
		nil,
		"FeatureGates to be enabled or disabled")
	_ = goflag.CommandLine.Parse([]string{})

	cmd.AddCommand(
		NewCmdBlockDevice(), //Add new command on block device
		NewCmdStart(),       //Add new command to start the ndm controller
	)

	return cmd, nil
}

// RunNodeDiskManager logs the starting of NDM daemon
func RunNodeDiskManager(cmd *cobra.Command) error {
	klog.Infof("Starting Node Device Manager Daemon...")
	klog.Infof("Version Tag : %s", version.GetVersion())
	klog.Infof("GitCommit : %s", version.GetGitCommit())
	return nil
}
