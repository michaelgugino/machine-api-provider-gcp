/*
Copyright 2018 The Kubernetes Authors.
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

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/golang/glog"
	clusterapis "github.com/openshift/cluster-api/pkg/apis"
	"github.com/openshift/cluster-api/pkg/controller/machine"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/klog"
	machineactuator "github.com/openshift/machine-api-provider-gcp/pkg/actuators/machine"
	"github.com/openshift/machine-api-provider-gcp/pkg/apis/awsproviderconfig/v1beta1"
	awsclient "github.com/openshift/machine-api-provider-gcp/pkg/client"
	"github.com/openshift/machine-api-provider-gcp/pkg/version"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

func main() {
	var printVersion bool
	flag.BoolVar(&printVersion, "version", false, "print version and exit")

	klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(klogFlags)
	flag.Parse()

	if printVersion {
		fmt.Println(version.String)
		os.Exit(0)
	}

	flag.VisitAll(func(f1 *flag.Flag) {
		f2 := klogFlags.Lookup(f1.Name)
		if f2 != nil {
			value := f1.Value.String()
			f2.Value.Set(value)
		}
	})

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		glog.Fatalf("Error getting configuration: %v", err)
	}

	// Create a new Cmd to provide shared dependencies and start components
	syncPeriod := 10 * time.Minute
	mgr, err := manager.New(cfg, manager.Options{
		SyncPeriod: &syncPeriod,
	})
	if err != nil {
		glog.Fatalf("Error creating manager: %v", err)
	}

	// Setup Scheme for all resources
	if err := clusterapis.AddToScheme(mgr.GetScheme()); err != nil {
		glog.Fatalf("Error setting up scheme: %v", err)
	}

	machineActuator, err := initActuator(mgr)
	if err != nil {
		glog.Fatalf("Error initializing actuator: %v", err)
	}

	if err := machine.AddWithActuator(mgr, machineActuator); err != nil {
		glog.Fatalf("Error adding actuator: %v", err)
	}

	// Start the Cmd
	err = mgr.Start(signals.SetupSignalHandler())
	if err != nil {
		glog.Fatalf("Error starting manager: %v", err)
	}
}

func initActuator(mgr manager.Manager) (*machineactuator.Actuator, error) {
	codec, err := v1beta1.NewCodec()
	if err != nil {
		return nil, fmt.Errorf("unable to create codec: %v", err)
	}

	params := machineactuator.ActuatorParams{
		Client:           mgr.GetClient(),
		Config:           mgr.GetConfig(),
		Codec:            codec,
		EventRecorder:    mgr.GetRecorder("aws-controller"),
	}

	actuator, err := machineactuator.NewActuator(params)
	if err != nil {
		return nil, fmt.Errorf("could not create AWS machine actuator: %v", err)
	}

	return actuator, nil
}
