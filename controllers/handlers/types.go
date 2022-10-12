package handlers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	resourcemanagmentv1alpha1 "github.com/tikalk/resource-manager/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//  keep everything in a struct

type Obj struct {
	Name types.NamespacedName
	c    client.Client
	ctx  context.Context
	l    logr.Logger
	rm   resourcemanagmentv1alpha1.ResourceManager

	stop chan bool
}

func InitObj(rm resourcemanagmentv1alpha1.ResourceManager, c client.Client, ctx context.Context, l logr.Logger) *Obj {
	stop := make(chan bool)

	name := types.NamespacedName{
		Name:      rm.Name,
		Namespace: rm.Namespace,
	}

	o := &Obj{
		Name: name,
		c:    c,
		ctx:  ctx,
		l:    l,
		stop: stop,
	}
	return o
}

func (o *Obj) Stop() {
	o.stop <- true
}

func (o *Obj) Run() {
	for {
		select {
		case <-o.stop:
			o.l.Info(fmt.Sprintf("%s Got stop signal!\n", h.Name))
			return
		default:
			switch o.rm.Spec.Resources { // here we decide which handler to use
			case "namespace":
				o.HandleNamespaceObj()
			}
		}
	}
}
