// SPDX-License-Identifier: Apache-2.0
// Copyright 2016-2020 Authors of Cilium

package watchers

import (
	"github.com/cilium/cilium/pkg/k8s"
	"github.com/cilium/cilium/pkg/k8s/informer"
	slim_corev1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/api/core/v1"
	"github.com/cilium/cilium/pkg/lock"
	"github.com/cilium/cilium/pkg/option"

	v1 "k8s.io/api/core/v1"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func (k *K8sWatcher) endpointsInit(k8sClient kubernetes.Interface, swgEps *lock.StoppableWaitGroup, optsModifier func(*v1meta.ListOptions)) {
	epOptsModifier := func(options *v1meta.ListOptions) {
		options.FieldSelector = fields.ParseSelectorOrDie(option.Config.K8sWatcherEndpointSelector).String()
		optsModifier(options)
	}

	_, endpointController := informer.NewInformer(
		cache.NewFilteredListWatchFromClient(k8sClient.CoreV1().RESTClient(),
			"endpoints", v1.NamespaceAll,
			epOptsModifier,
		),
		&slim_corev1.Endpoints{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				var valid, equal bool
				defer func() { k.K8sEventReceived(metricEndpoint, metricCreate, valid, equal) }()
				if k8sEP := k8s.ObjToV1Endpoints(obj); k8sEP != nil {
					valid = true
					err := k.addK8sEndpointV1(k8sEP, swgEps)
					k.K8sEventProcessed(metricEndpoint, metricCreate, err == nil)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				var valid, equal bool
				defer func() { k.K8sEventReceived(metricEndpoint, metricUpdate, valid, equal) }()
				if oldk8sEP := k8s.ObjToV1Endpoints(oldObj); oldk8sEP != nil {
					if newk8sEP := k8s.ObjToV1Endpoints(newObj); newk8sEP != nil {
						valid = true
						if oldk8sEP.DeepEqual(newk8sEP) {
							equal = true
							return
						}

						err := k.updateK8sEndpointV1(oldk8sEP, newk8sEP, swgEps)
						k.K8sEventProcessed(metricEndpoint, metricUpdate, err == nil)
					}
				}
			},
			DeleteFunc: func(obj interface{}) {
				var valid, equal bool
				defer func() { k.K8sEventReceived(metricEndpoint, metricDelete, valid, equal) }()
				k8sEP := k8s.ObjToV1Endpoints(obj)
				if k8sEP == nil {
					return
				}
				valid = true
				err := k.deleteK8sEndpointV1(k8sEP, swgEps)
				k.K8sEventProcessed(metricEndpoint, metricDelete, err == nil)
			},
		},
		nil,
	)
	k.blockWaitGroupToSyncResources(wait.NeverStop, swgEps, endpointController.HasSynced, K8sAPIGroupEndpointV1Core)
	go endpointController.Run(wait.NeverStop)
	k.k8sAPIGroups.AddAPI(K8sAPIGroupEndpointV1Core)
}

func (k *K8sWatcher) addK8sEndpointV1(ep *slim_corev1.Endpoints, swg *lock.StoppableWaitGroup) error {
	return k.updateK8sEndpointV1(nil, ep, swg)
}

func (k *K8sWatcher) updateK8sEndpointV1(oldEP, newEP *slim_corev1.Endpoints, swg *lock.StoppableWaitGroup) error {
	k.K8sSvcCache.UpdateEndpoints(newEP, swg)
	if option.Config.BGPAnnounceLBIP {
		k.bgpSpeakerManager.OnUpdateEndpoints(newEP)
	}
	return nil
}

func (k *K8sWatcher) deleteK8sEndpointV1(ep *slim_corev1.Endpoints, swg *lock.StoppableWaitGroup) error {
	k.K8sSvcCache.DeleteEndpoints(ep, swg)
	return nil
}
