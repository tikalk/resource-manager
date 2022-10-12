package handlers

import (
	"context"
	"fmt"
	"k8s.io/client-go/tools/cache"

	"github.com/go-logr/logr"
	resourcemanagmentv1alpha1 "github.com/tikalk/resource-manager/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//  keep everything in a struct

type ResourceManagerHandler struct {
	resourceManager resourcemanagmentv1alpha1.ResourceManager
	objInformer     cache.SharedIndexInformer
	objHandlers     map[string]*ObjectHandler

	//Name            types.NamespacedName
	client client.Client
	ctx    context.Context
	log    logr.Logger

	stop chan struct{}
}

func InitObj(rm resourcemanagmentv1alpha1.ResourceManager, c client.Client, ctx context.Context, l logr.Logger) ResourceManagerHandler {
	stop := make(chan struct{})

	//name := types.NamespacedName{
	//	Name:      rm.Name,
	//	Namespace: rm.Namespace,
	//}

	o := ResourceManagerHandler{
		resourceManager: rm,
		//Name:            name,
		client: c,
		ctx:    ctx,
		log:    l,
		stop:   stop,
	}
	return o
}

func (o ResourceManagerHandler) Stop() {
	o.stop <- struct{}{}
}

func (o ResourceManagerHandler) Run() {
	fmt.Printf("Processing object: %s \n", o.resourceManager.Name)
	for {
		select {
		case <-o.stop:
			o.log.Info(fmt.Sprintf("%s Got stop signal!\n", o.resourceManager.Name))
			return
		default:
			switch o.resourceManager.Spec.Resources { // here we decide which handler to use
			case "namespace":
				o.HandleNamespaceObj()
			}
		}
	}
}
