/*
Copyright 2018 The vault-operator Authors

Commercial software license.
*/
package versioned

import (
	vaultv1alpha1 "github.com/coreos/vault-operator/pkg/generated/clientset/versioned/typed/vault/v1alpha1"
	glog "github.com/golang/glog"
	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	VaultV1alpha1() vaultv1alpha1.VaultV1alpha1Interface
	// Deprecated: please explicitly pick a version if possible.
	Vault() vaultv1alpha1.VaultV1alpha1Interface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	*discovery.DiscoveryClient
	vaultV1alpha1 *vaultv1alpha1.VaultV1alpha1Client
}

// VaultV1alpha1 retrieves the VaultV1alpha1Client
func (c *Clientset) VaultV1alpha1() vaultv1alpha1.VaultV1alpha1Interface {
	return c.vaultV1alpha1
}

// Deprecated: Vault retrieves the default version of VaultClient.
// Please explicitly pick a version.
func (c *Clientset) Vault() vaultv1alpha1.VaultV1alpha1Interface {
	return c.vaultV1alpha1
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}

// NewForConfig creates a new Clientset for the given config.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	var cs Clientset
	var err error
	cs.vaultV1alpha1, err = vaultv1alpha1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(&configShallowCopy)
	if err != nil {
		glog.Errorf("failed to create the DiscoveryClient: %v", err)
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.vaultV1alpha1 = vaultv1alpha1.NewForConfigOrDie(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClientForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.vaultV1alpha1 = vaultv1alpha1.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}
