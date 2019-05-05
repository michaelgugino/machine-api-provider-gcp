package machine

import (
	"fmt"
    "context"
    "k8s.io/klog"
    corev1 "k8s.io/api/core/v1"
	machinev1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	providerconfigv1 "github.com/openshift/machine-api-provider-gcp/pkg/apis/gcpproviderconfig/v1alpha1"
    "k8s.io/apimachinery/pkg/types"
)

// providerConfigFromMachine gets the machine provider config MachineSetSpec from the
// specified cluster-api MachineSpec.
func providerConfigFromMachine(machine *machinev1.Machine, codec *providerconfigv1.GCPProviderConfigCodec) (*providerconfigv1.GCPMachineProviderConfig, error) {
	if machine.Spec.ProviderSpec.Value == nil {
		return nil, fmt.Errorf("unable to find machine provider config: Spec.ProviderSpec.Value is not set")
	}

	var config providerconfigv1.GCPMachineProviderConfig
	if err := codec.DecodeProviderSpec(&machine.Spec.ProviderSpec, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// isMaster returns true if the machine is part of a cluster's control plane
func (a *Actuator) isMaster(machine *machinev1.Machine) (bool, error) {
	if machine.Status.NodeRef == nil {
		klog.Errorf("NodeRef not found in machine %s", machine.Name)
		return false, nil
	}
	node := &corev1.Node{}
	nodeKey := types.NamespacedName{
		Namespace: machine.Status.NodeRef.Namespace,
		Name:      machine.Status.NodeRef.Name,
	}

	err := a.client.Get(context.Background(), nodeKey, node)
	if err != nil {
		return false, fmt.Errorf("failed to get node from machine %s", machine.Name)
	}

	if _, exists := node.Labels["node-role.kubernetes.io/master"]; exists {
		return true, nil
	}
	return false, nil
}
