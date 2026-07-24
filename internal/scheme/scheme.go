package scheme

import (
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	externaldnsv1alpha1 "sigs.k8s.io/external-dns/apis/v1alpha1"
)

func NewScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(externaldnsv1alpha1.AddToScheme(scheme))
	return scheme
}
