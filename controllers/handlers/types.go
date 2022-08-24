package handlers

import (
	"context"
	"github.com/go-logr/logr"
	resourcemanagmentv1alpha1 "github.com/tikalk/resource-manager/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//  keep everything in a struct

type Obj struct {
	Name types.NamespacedName
	C    client.Client
	Ctx  context.Context
	L    logr.Logger
	Spec resourcemanagmentv1alpha1.ResourceManagerSpec
}
