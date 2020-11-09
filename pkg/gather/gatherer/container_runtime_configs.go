package gatherer

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	_ "k8s.io/apimachinery/pkg/runtime/serializer/yaml"

	"github.com/openshift/insights-operator/pkg/record"
)

// GatherContainerRuntimeConfig collects ContainerRuntimeConfig  information
//
// The Kubernetes api https://github.com/openshift/machine-config-operator/blob/master/pkg/apis/machineconfiguration.openshift.io/v1/types.go#L402
// Response see https://docs.okd.io/latest/rest_api/machine_apis/containerruntimeconfig-machineconfiguration-openshift-io-v1.html
//
// Location in archive: config/containerruntimeconfigs/
func GatherContainerRuntimeConfig(i *Gatherer) func() ([]record.Record, []error) {
	return func() ([]record.Record, []error) {
		crc := schema.GroupVersionResource{Group: "machineconfiguration.openshift.io", Version: "v1", Resource: "containerruntimeconfigs"}
		containerRCs, err := i.dynamicClient.Resource(crc).List(i.ctx, metav1.ListOptions{})
		if errors.IsNotFound(err) {
			return nil, nil
		}
		if err != nil {
			return nil, []error{err}
		}

		records := []record.Record{}
		for _, i := range containerRCs.Items {
			records = append(records, record.Record{
				Name: fmt.Sprintf("config/containerruntimeconfigs/%s", i.GetName()),
				Item: record.JSONMarshaller{Object: i.Object},
			})
		}
		return records, nil
	}
}
