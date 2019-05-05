/*
Copyright 2019 The Kubernetes Authors.

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

package machine

import (
	"context"
	"fmt"
    "errors"
//	"time"
    "strings"
	corev1 "k8s.io/api/core/v1"
//	"k8s.io/apimachinery/pkg/api/equality"
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	errorutil "k8s.io/apimachinery/pkg/util/errors"
//	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"

	clusterv1 "github.com/openshift/cluster-api/pkg/apis/cluster/v1alpha1"
	machinev1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	//clustererror "github.com/openshift/cluster-api/pkg/controller/error"
	apierrors "github.com/openshift/cluster-api/pkg/errors"
	providerconfigv1 "github.com/openshift/machine-api-provider-gcp/pkg/apis/gcpproviderconfig/v1alpha1"

	"sigs.k8s.io/controller-runtime/pkg/client"

    compute "google.golang.org/api/compute/v1"
    "github.com/openshift/cluster-api-provider-gcp/pkg/cloud/google/machinesetup"
    gcecloud "github.com/openshift/cluster-api-provider-gcp/pkg/cloud/google"
    "k8s.io/apimachinery/pkg/runtime"
)

const (
	userDataSecretKey         = "userData"
	requeueAfterSeconds       = 20

	// MachineCreationSucceeded indicates success for machine creation
	MachineCreationSucceeded = "MachineCreationSucceeded"

	// MachineCreationFailed indicates that machine creation failed
	MachineCreationFailed = "MachineCreationFailed"
)

const (
	ProjectAnnotationKey = "gcp-project"
	ZoneAnnotationKey    = "gcp-zone"
	NameAnnotationKey    = "gcp-name"

	BootstrapLabelKey = "bootstrap"

	// This file is a yaml that will be used to create the machine-setup configmap on the machine controller.
	// It contains the supported machine configurations along with the startup scripts and OS image paths that correspond to each supported configuration.
	MachineSetupConfigsFilename = "machine_setup_configs.yaml"
	ProviderName                = "google"
)

// Actuator is the GCP-specific actuator for the Cluster API machine controller
type Actuator struct {
	client             client.Client
	codec              *providerconfigv1.GCPProviderConfigCodec
	eventRecorder      record.EventRecorder
}

// ActuatorParams holds parameter information for Actuator
type ActuatorParams struct {
	Client           client.Client
	Codec            *providerconfigv1.GCPProviderConfigCodec
	EventRecorder    record.EventRecorder
}

// NewActuator returns a new GCP Actuator
func NewActuator(params ActuatorParams) (*Actuator, error) {
	actuator := &Actuator{
		client:           params.Client,
		codec:            params.Codec,
		eventRecorder:    params.EventRecorder,
	}
	return actuator, nil
}

const (
	createEventAction = "Create"
	updateEventAction = "Update"
	deleteEventAction = "Delete"
	noEventAction     = ""
)

// Set corresponding event based on error. It also returns the original error
// for convenience, so callers can do "return handleMachineError(...)".
func (a *Actuator) handleMachineError(machine *machinev1.Machine, err *apierrors.MachineError, eventAction string) error {
	if eventAction != noEventAction {
		a.eventRecorder.Eventf(machine, corev1.EventTypeWarning, "Failed"+eventAction, "%v", err.Reason)
	}

	klog.Errorf("Machine error: %v", err.Message)
	return err
}

// Create runs a new GCP instance
func (a *Actuator) Create(context context.Context, cluster *clusterv1.Cluster, machine *machinev1.Machine) error {
	klog.Info("creating machine")
    /*
	instance, err := a.CreateMachine(cluster, machine)
	if err != nil {
		klog.Errorf("error creating machine: %v", err)
		updateConditionError := a.updateMachineProviderConditions(machine, providerconfigv1.MachineCreation, MachineCreationFailed, err.Error())
		if updateConditionError != nil {
			klog.Errorf("error updating machine conditions: %v", updateConditionError)
		}
		return err
	}
	// return a.updateStatus(machine, instance)
    */
    return errors.New("Not Implemented")

}

/*
func (a *Actuator) updateMachineStatus(machine *machinev1.Machine, awsStatus *providerconfigv1.GCPMachineProviderStatus, networkAddresses []corev1.NodeAddress) error {
	awsStatusRaw, err := a.codec.EncodeProviderStatus(awsStatus)
	if err != nil {
		klog.Errorf("error encoding GCP provider status: %v", err)
		return err
	}

	machineCopy := machine.DeepCopy()
	machineCopy.Status.ProviderStatus = awsStatusRaw
	if networkAddresses != nil {
		machineCopy.Status.Addresses = networkAddresses
	}

	oldGCPStatus := &providerconfigv1.GCPMachineProviderStatus{}
	if err := a.codec.DecodeProviderStatus(machine.Status.ProviderStatus, oldGCPStatus); err != nil {
		klog.Errorf("error updating machine status: %v", err)
		return err
	}

	// TODO(vikasc): Revisit to compare complete machine status objects
	if !equality.Semantic.DeepEqual(awsStatus, oldGCPStatus) || !equality.Semantic.DeepEqual(machine.Status.Addresses, machineCopy.Status.Addresses) {
		klog.Infof("machine status has changed, updating")
		time := metav1.Now()
		machineCopy.Status.LastUpdated = &time

		if err := a.client.Status().Update(context.Background(), machineCopy); err != nil {
			klog.Errorf("error updating machine status: %v", err)
			return err
		}
	} else {
		klog.Info("status unchanged")
	}

	return nil
}
*/

/*
// updateMachineProviderConditions updates conditions set within machine provider status.
func (a *Actuator) updateMachineProviderConditions(machine *machinev1.Machine, conditionType providerconfigv1.GCPMachineProviderConditionType, reason string, msg string) error {

	klog.Info("updating machine conditions")

	awsStatus := &providerconfigv1.GCPMachineProviderStatus{}
	if err := a.codec.DecodeProviderStatus(machine.Status.ProviderStatus, awsStatus); err != nil {
		klog.Errorf("error decoding machine provider status: %v", err)
		return err
	}

	awsStatus.Conditions = setGCPMachineProviderCondition(awsStatus.Conditions, conditionType, corev1.ConditionTrue, reason, msg, updateConditionIfReasonOrMessageChange)

	if err := a.updateMachineStatus(machine, awsStatus, nil); err != nil {
		return err
	}

	return nil
}
*/

// CreateMachine starts a new GCP instance as described by the cluster and machine resources
func (a *Actuator) CreateMachine(cluster *clusterv1.Cluster, machine *machinev1.Machine) (*compute.Instance, error) {
    /*
	machineProviderConfig, err := providerConfigFromMachine(machine, a.codec)
	if err != nil {
		return nil, a.handleMachineError(machine, apierrors.InvalidMachineConfiguration("error decoding MachineProviderConfig: %v", err), createEventAction)
	}


	credentialsSecretName := ""
	if machineProviderConfig.CredentialsSecret != nil {
		credentialsSecretName = machineProviderConfig.CredentialsSecret.Name
	}


    // create gcp client
	// awsClient, err := a.awsClientBuilder(a.client, credentialsSecretName, machine.Namespace, machineProviderConfig.Placement.Region)
	if err != nil {
		klog.Errorf("unable to obtain GCP client: %v", err)
		return nil, a.handleMachineError(machine, apierrors.CreateMachine("error creating gcp services: %v", err), createEventAction)
	}


	// We explicitly do NOT want to remove stopped masters.
	isMaster, err := a.isMaster(machine)
	// Unable to determine if a machine is a master machine.
	// Yet, it's only used to delete stopped machines that are not masters.
	// So we can safely continue to create a new machine since in the worst case
	// we just don't delete any stopped machine.
	if err != nil {
		klog.Errorf("Error determining if machine is master: %v", err)
	} else {
		if !isMaster {
			// Prevent having a lot of stopped nodes sitting around.
			err = removeStoppedMachine(machine, awsClient)
			if err != nil {
				klog.Errorf("unable to remove stopped machines: %v", err)
				return nil, fmt.Errorf("unable to remove stopped nodes: %v", err)
			}
		}
	}


	userData := []byte{}
	if machineProviderConfig.UserDataSecret != nil {
		var userDataSecret corev1.Secret
		err := a.client.Get(context.Background(), client.ObjectKey{Namespace: machine.Namespace, Name: machineProviderConfig.UserDataSecret.Name}, &userDataSecret)
		if err != nil {
			return nil, a.handleMachineError(machine, apierrors.CreateMachine("error getting user data secret %s: %v", machineProviderConfig.UserDataSecret.Name, err), createEventAction)
		}
		if data, exists := userDataSecret.Data[userDataSecretKey]; exists {
			userData = data
		} else {
			klog.Warningf("Secret %v/%v does not have %q field set. Thus, no user data applied when creating an instance.", machine.Namespace, machineProviderConfig.UserDataSecret.Name, userDataSecretKey)
		}
	}

	//instance, err := launchInstance(machine, machineProviderConfig, userData, awsClient)
    instance, err := Create2x()
	if err != nil {
		return nil, a.handleMachineError(machine, apierrors.CreateMachine("error launching instance: %v", err), createEventAction)
	}

	a.eventRecorder.Eventf(machine, corev1.EventTypeNormal, "Created", "Created Machine %v", machine.Name)
	return instance, nil
    */
    return nil, errors.New("Not Implemented")
}

// Delete deletes a machine and updates its finalizer

func (a *Actuator) Delete(context context.Context, cluster *clusterv1.Cluster, machine *machinev1.Machine) error {
    /*
	klog.Info("deleting machine")
	if err := a.DeleteMachine(cluster, machine); err != nil {
		klog.Errorf("error deleting machine: %v", err)
		return err
	}
    */
	return nil
}

func Create2x() (*compute.Instance, error) {
    return nil, errors.New("error 1")
}

type klogLogger struct{}

func (gl *klogLogger) Log(v ...interface{}) {
	klog.Info(v...)
}

func (gl *klogLogger) Logf(format string, v ...interface{}) {
	klog.Infof(format, v...)
}

// DeleteMachine deletes an GCP instance
/*
func (a *Actuator) DeleteMachine(cluster *clusterv1.Cluster, machine *machinev1.Machine) error {
	machineProviderConfig, err := providerConfigFromMachine(machine, a.codec)
	if err != nil {
		return a.handleMachineError(machine, apierrors.InvalidMachineConfiguration("error decoding MachineProviderConfig: %v", err), deleteEventAction)
	}

	region := machineProviderConfig.Placement.Region
	credentialsSecretName := ""
	if machineProviderConfig.CredentialsSecret != nil {
		credentialsSecretName = machineProviderConfig.CredentialsSecret.Name
	}
	client, err := a.awsClientBuilder(a.client, credentialsSecretName, machine.Namespace, region)
	if err != nil {
		errMsg := fmt.Errorf("error getting EC2 client: %v", err)
		klog.Error(errMsg)
		return errMsg
	}

	instances, err := getRunningInstances(machine, client)
	if err != nil {
		klog.Errorf("error getting running instances: %v", err)
		return err
	}
	if len(instances) == 0 {
		klog.Warningf("no instances found to delete for machine")
		return nil
	}

	err = terminateInstances(client, instances)
	if err != nil {
		return a.handleMachineError(machine, apierrors.DeleteMachine(err.Error()), noEventAction)
	}
	a.eventRecorder.Eventf(machine, corev1.EventTypeNormal, "Deleted", "Deleted machine %v", machine.Name)

	return nil
}
*/

// Update attempts to sync machine state with an existing instance. Today this just updates status
// for details that may have changed. (IPs and hostnames) We do not currently support making any
// changes to actual machines in GCP. Instead these will be replaced via MachineDeployments.

func (a *Actuator) Update(context context.Context, cluster *clusterv1.Cluster, machine *machinev1.Machine) error {
	klog.Info("updating machine")

    return nil
/*
	machineProviderConfig, err := providerConfigFromMachine(machine, a.codec)
	if err != nil {
		return a.handleMachineError(machine, apierrors.InvalidMachineConfiguration("error decoding MachineProviderConfig: %v", err), updateEventAction)
	}

	region := machineProviderConfig.Placement.Region
	klog.Info("obtaining EC2 client for region")
	credentialsSecretName := ""
	if machineProviderConfig.CredentialsSecret != nil {
		credentialsSecretName = machineProviderConfig.CredentialsSecret.Name
	}
	client, err := a.awsClientBuilder(a.client, credentialsSecretName, machine.Namespace, region)
	if err != nil {
		errMsg := fmt.Errorf("error getting EC2 client: %v", err)
		klog.Error(errMsg)
		return errMsg
	}

	instances, err := getRunningInstances(machine, client)
	if err != nil {
		klog.Errorf("error getting running instances: %v", err)
		return err
	}
	klog.Infof("found %d instances for machine", len(instances))

	// Parent controller should prevent this from ever happening by calling Exists and then Create,
	// but instance could be deleted between the two calls.
	if len(instances) == 0 {
		klog.Warningf("attempted to update machine but no instances found")

		a.handleMachineError(machine, apierrors.UpdateMachine("no instance found, reason unknown"), updateEventAction)

		// Update status to clear out machine details.
		if err := a.updateStatus(machine, nil); err != nil {
			return err
		}

		errMsg := "attempted to update machine but no instances found"
		klog.Error(errMsg)
		return fmt.Errorf(errMsg)
	}

	klog.Info("instance found")

	// In very unusual circumstances, there could be more than one machine running matching this
	// machine name and cluster ID. In this scenario we will keep the newest, and delete all others.
	sortInstances(instances)
	if len(instances) > 1 {
		err = terminateInstances(client, instances[1:])
		if err != nil {
			klog.Errorf("Unable to remove redundant instances: %v", err)
			return err
		}
	}

	newestInstance := instances[0]

	err = a.updateLoadBalancers(client, machineProviderConfig, newestInstance)
	if err != nil {
		a.handleMachineError(machine, apierrors.CreateMachine("Error updating load balancers: %v", err), updateEventAction)
		return err
	}

	a.eventRecorder.Eventf(machine, corev1.EventTypeNormal, "Updated", "Updated machine %v", machine.Name)

	// We do not support making changes to pre-existing instances, just update status.
	return a.updateStatus(machine, newestInstance)
}
*/
}

// Exists determines if the given machine currently exists. For GCP we query for instances in
// running state, with a matching name tag, to determine a match.
func (a *Actuator) Exists(_ context.Context, cluster *clusterv1.Cluster, machine *machinev1.Machine) (bool, error) {
	klog.Info("Checking if machine exists")

	instances, err := a.getMachineInstances(cluster, machine)
	if err != nil {
		klog.Errorf("Error getting running instances: %v", err)
		return false, err
	}
	if len(instances) == 0 {
		klog.Info("Instance does not exist")
		return false, nil
	}

	// If more than one result was returned, it will be handled in Update.
	//klog.Infof("Instance exists as %q", *instances[0].InstanceId)
	return true, nil
}

// Describe provides information about machine's instance(s)
func (a *Actuator) Describe(cluster *clusterv1.Cluster, machine *machinev1.Machine) (*compute.Instance, error) {
	klog.Infof("Checking if machine exists")

	instances, err := a.getMachineInstances(cluster, machine)
	if err != nil {
		klog.Errorf("Error getting running instances: %v", err)
		return nil, err
	}
	if len(instances) == 0 {
		klog.Info("Instance does not exist")
		return nil, nil
	}

	return instances[0], nil
}

func (a *Actuator) getMachineInstances(cluster *clusterv1.Cluster, machine *machinev1.Machine) ([]*compute.Instance, error) {
    /*
	machineProviderConfig, err := providerConfigFromMachine(machine, a.codec)
	if err != nil {
		klog.Errorf("Error decoding MachineProviderConfig: %v", err)
		return nil, err
	}
    */
    return nil, errors.New("Not Implemented")
    /*
	region := machineProviderConfig.Placement.Region
	credentialsSecretName := ""
	if machineProviderConfig.CredentialsSecret != nil {
		credentialsSecretName = machineProviderConfig.CredentialsSecret.Name
	}
	client, err := a.awsClientBuilder(a.client, credentialsSecretName, machine.Namespace, region)
	if err != nil {
		klog.Errorf("Error getting EC2 client: %v", err)
		return nil, fmt.Errorf("error getting EC2 client: %v", err)
	}

	return getRunningInstances(machine, client)
    */
}

/*
// updateLoadBalancers adds a given machine instance to the load balancers specified in its provider config
func (a *Actuator) updateLoadBalancers(client awsclient.Client, providerConfig *providerconfigv1.GCPMachineProviderConfig, instance *compute.Instance) error {
	if len(providerConfig.LoadBalancers) == 0 {
		klog.V(4).Infof("Instance %q has no load balancers configured. Skipping", *instance.InstanceId)
		return nil
	}
	errs := []error{}
	classicLoadBalancerNames := []string{}
	networkLoadBalancerNames := []string{}
	for _, loadBalancerRef := range providerConfig.LoadBalancers {
		switch loadBalancerRef.Type {
		case providerconfigv1.NetworkLoadBalancerType:
			networkLoadBalancerNames = append(networkLoadBalancerNames, loadBalancerRef.Name)
		case providerconfigv1.ClassicLoadBalancerType:
			classicLoadBalancerNames = append(classicLoadBalancerNames, loadBalancerRef.Name)
		}
	}

	var err error
	if len(classicLoadBalancerNames) > 0 {
		err := registerWithClassicLoadBalancers(client, classicLoadBalancerNames, instance)
		if err != nil {
			klog.Errorf("Failed to register classic load balancers: %v", err)
			errs = append(errs, err)
		}
	}
	if len(networkLoadBalancerNames) > 0 {
		err = registerWithNetworkLoadBalancers(client, networkLoadBalancerNames, instance)
		if err != nil {
			klog.Errorf("Failed to register network load balancers: %v", err)
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errorutil.NewAggregate(errs)
	}
	return nil
}
*/

/*
// updateStatus calculates the new machine status, checks if anything has changed, and updates if so.
func (a *Actuator) updateStatus(machine *machinev1.Machine, instance *compute.Instance) error {

	klog.Info("Updating status")

	return nil
}
*/

func getClusterID(machine *machinev1.Machine) (string, bool) {
	clusterID, ok := machine.Labels[providerconfigv1.ClusterIDLabel]
	// NOTE: This block can be removed after the label renaming transition to machine.openshift.io
	if !ok {
		clusterID, ok = machine.Labels["sigs.k8s.io/cluster-api-cluster"]
	}
	return clusterID, ok
}

var MachineActuator *GCEClient

type GCEClientMachineSetupConfigGetter interface {
	GetMachineSetupConfig() (machinesetup.MachineSetupConfig, error)
}

type GCEClient struct {
	computeService           gcecloud.GCEClientComputeService
	client                   client.Client
	machineSetupConfigGetter GCEClientMachineSetupConfigGetter
	eventRecorder            record.EventRecorder
	scheme                   *runtime.Scheme
}

func newDisks2() []*compute.AttachedDisk{
	var disks []*compute.AttachedDisk
	//for idx, disk := range config.Disks {
	//diskSizeGb := disk.InitializeParams.DiskSizeGb
	d := compute.AttachedDisk{
		AutoDelete: true,
		InitializeParams: &compute.AttachedDiskInitializeParams{
			DiskSizeGb: 30,
			DiskType:   fmt.Sprintf("zones/%s/diskTypes/%s", "us-east1-c", "pd-ssd"),
		},
	}
	d.InitializeParams.SourceImage = "projects/debian-cloud/global/images/debian-9-stretch-v20190423"
	d.Boot = true
	disks = append(disks, &d)

	return disks
}

func (gce *GCEClient) getMetadata2() (*compute.Metadata, error) {
	metadataMap := make(map[string]string)
	{
		var b strings.Builder

		project := "openshift-gce-devel"

		clusterName := "mgdev"
		nodeTag := clusterName + "-worker"

		network := "default"
		subnetwork := "kubernetes"

		fmt.Fprintf(&b, "[global]\n")
		fmt.Fprintf(&b, "project-id = %s\n", project)
		fmt.Fprintf(&b, "network-name = %s\n", network)
		fmt.Fprintf(&b, "subnetwork-name = %s\n", subnetwork)
		fmt.Fprintf(&b, "node-tags = %s\n", nodeTag)
		metadataMap["cloud-config"] = b.String()
	}

	var metadataItems []*compute.MetadataItems
	for k, v := range metadataMap {
		v := v // rebind scope to avoid loop aliasing below
		metadataItems = append(metadataItems, &compute.MetadataItems{
			Key:   k,
			Value: &v,
		})
	}
	metadata := compute.Metadata{
		Items: metadataItems,
	}
	return &metadata, nil
}
