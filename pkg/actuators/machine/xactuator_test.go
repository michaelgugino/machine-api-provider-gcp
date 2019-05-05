package machine

import (
	"testing"
	//"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	providerconfigv1 "github.com/openshift/machine-api-provider-gcp/pkg/apis/gcpproviderconfig/v1alpha1"
	"k8s.io/client-go/tools/record"
)

func TestCreateActuator(t *testing.T) {
	fakeClient := fake.NewFakeClient()
	//mockCtrl := gomock.NewController(t)

	codec, err := providerconfigv1.NewCodec()
	if err != nil {
		t.Fatalf("unable to build codec: %v", err)
	}

	params := ActuatorParams{
		Client: fakeClient,
		Codec: codec,
		// use empty recorder dropping any event recorded
		EventRecorder: &record.FakeRecorder{},
	}

	_, err2 := NewActuator(params)
	if err2 != nil {
		t.Fatalf("Could not create AWS machine actuator: %v", err)
	}
}

/*
func TestListMachines(t *testing.T) {
	mp := MachineActuatorParams{}
	gclient, err := NewMachineActuator(mp)
    if err != nil {
		t.Fatalf("gclient borke")
	}
	_ = gclient.List()

    err2 := gclient.Create2()
    if err2 != nil {
        t.Fatalf("create fail")
    }
}
*/
