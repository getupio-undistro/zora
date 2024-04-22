// Code generated by informer-gen. DO NOT EDIT.

package v1alpha2

import (
	internalinterfaces "github.com/undistro/zora/pkg/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// CustomChecks returns a CustomCheckInformer.
	CustomChecks() CustomCheckInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// CustomChecks returns a CustomCheckInformer.
func (v *version) CustomChecks() CustomCheckInformer {
	return &customCheckInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
