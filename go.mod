module github.com/danielfoehrkn/kubeswitch

go 1.12

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.19.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v0.11.0
	// v0.1.0 is incompatible with azidentity v0.11.0. Wait for azidentity to be updated.
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice v0.1.0
	github.com/aws/aws-sdk-go-v2/config v1.9.0
	github.com/aws/aws-sdk-go-v2/service/eks v1.11.0
	github.com/disiqueira/gotree v1.0.0
	github.com/gardener/gardener v1.18.0
	github.com/hashicorp/go-multierror v1.1.0
	github.com/hashicorp/vault/api v1.0.4
	github.com/jedib0t/go-pretty/v6 v6.1.0
	github.com/karrick/godirwalk v1.16.1
	github.com/ktr0731/go-fuzzyfinder v0.5.0
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.5
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.2.1
	golang.org/x/net v0.0.0-20211101193420-4a448f8816b3 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.7 // indirect
	google.golang.org/api v0.57.0
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api v0.19.6
	k8s.io/apimachinery v0.19.6
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/kube-openapi v0.0.0-20210305164622-f622666832c1 // indirect
	k8s.io/utils v0.0.0-20200912215256-4140de9c8800
	sigs.k8s.io/controller-runtime v0.7.1
)

replace (
	//github.com/Azure/azure-sdk-for-go/sdk/azcore => github.com/Azure/azure-sdk-for-go/sdk/azcore v0.20.0
	//github.com/Azure/azure-sdk-for-go/sdk/internal => github.com/Azure/azure-sdk-for-go/sdk/internal v0.7.0
	// pin to same version to avoid dependency issues
	k8s.io/api => k8s.io/api v0.19.6
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.6
	k8s.io/apiserver => k8s.io/apiserver v0.19.6
	k8s.io/client-go => k8s.io/client-go v0.19.6
)
