package machine

import (
	"testing"
	"k8s.io/klog"
	"flag"
	"github.com/golang/mock/gomock"
	//"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	providerconfigv1 "github.com/openshift/machine-api-provider-gcp/pkg/apis/gceproviderconfig/v1alpha1"
	"k8s.io/client-go/tools/record"
)
func init() {
	klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(klogFlags)
}
func TestCreateActuator(t *testing.T) {
	_ = gomock.NewController(t)
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
	service, err := serviceFromJWTSecret("none", "none")
	if err != nil {
		t.Fatalf("Unable to get compute service")
	}
	_ = instanceList(service)
	_, _ = getMetadata2()
	_ = Create2(service)
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
