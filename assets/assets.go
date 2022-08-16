// FOR FUTURE USE I CREATED assets PACKAGE TO BE ABLE TO CREATE / DEPLOY ANY OBJECT FROM YAML FILE.
// CAN BE DONE BY:
// import "https://github.com/tikalk/resource-manager/assets"
// ...
// namespaceDeployment := assets.GetDeploymentFromFile("manifests/namespace_deploy.yaml")

package assets

// Imports the relevant k8s API packages that define the schema for Deployment API objects
import (
	"embed"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

// It initializes a Scheme and a set of codecs that can be used by the API's UniversalDecoder
// in order to know how to convert the []byte data representation of the file to a Go struct
var (
	//go:embed manifests/*
	manifests  embed.FS
	appsScheme = runtime.NewScheme()
	appsCodecs = serializer.NewCodecFactory(appsScheme)
)

func init() {
	if err := appsv1.AddToScheme(appsScheme); err != nil {
		panic(err)
	}
}

// It uses the "namespaceDeployment :=" variable we can declare
// to read the Deployment file under assets/namespace_deploy.yaml
func GetDeploymentFromFile(name string) *appsv1.Deployment {
	deploymentBytes, err := manifests.ReadFile(name)
	if err != nil {
		panic(err)
	}

	// It decodes the []byte data returned from deployment.ReadFile()
	// into an object that can be cast to the Go type for Deployments
	deploymentObject, err := runtime.Decode(
		appsCodecs.UniversalDecoder(appsv1.SchemeGroupVersion),
		deploymentBytes,
	)
	if err != nil {
		panic(err)
	}
	return deploymentObject.(*appsv1.Deployment)
}
