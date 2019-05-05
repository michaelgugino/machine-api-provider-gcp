package v1alpha1

import (
	"reflect"
	"testing"
	"time"

	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestEncodeAndDecodeProviderStatus(t *testing.T) {

	codec, err := NewCodec()
	if err != nil {
		t.Fatal(err)
	}
	time := metav1.Time{
		Time: time.Date(2018, 6, 3, 0, 0, 0, 0, time.Local),
	}

	instanceState := "running"
	instanceID := "id"
	providerStatus := &GCPMachineProviderStatus{
		TypeMeta: metav1.TypeMeta{
			Kind:       "GCPMachineProviderStatus",
			APIVersion: "gcpproviderconfig.openshift.io/v1alpha1",
		},
		InstanceState: &instanceState,
		InstanceID:    &instanceID,
		Conditions: []GCPMachineProviderCondition{
			{
				Type:               "MachineCreation",
				Status:             "True",
				Reason:             "MachineCreationSucceeded",
				Message:            "machine successfully created",
				LastTransitionTime: time,
				LastProbeTime:      time,
			},
		},
	}
	providerStatusEncoded, err := codec.EncodeProviderStatus(providerStatus)
	if err != nil {
		t.Error(err)
	}

	// without deep copy
	{
		providerStatusDecoded := &GCPMachineProviderStatus{}
		codec.DecodeProviderStatus(providerStatusEncoded, providerStatusDecoded)
		if !reflect.DeepEqual(providerStatus, providerStatusDecoded) {
			t.Errorf("failed EncodeProviderStatus/DecodeProviderStatus. Expected: %+v, got: %+v", providerStatus, providerStatusDecoded)
		}
	}

	// with deep copy
	{
		providerStatusDecoded := &GCPMachineProviderStatus{}
		codec.DecodeProviderStatus(providerStatusEncoded.DeepCopy(), providerStatusDecoded)
		if !reflect.DeepEqual(providerStatus, providerStatusDecoded) {
			t.Errorf("failed EncodeProviderStatus/DecodeProviderStatus. Expected: %+v, got: %+v", providerStatus, providerStatusDecoded)
		}
	}
}
