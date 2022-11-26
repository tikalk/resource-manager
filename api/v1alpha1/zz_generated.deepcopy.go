//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2022.

Dev in DevOps course @ Tikal
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterResourceManager) DeepCopyInto(out *ClusterResourceManager) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterResourceManager.
func (in *ClusterResourceManager) DeepCopy() *ClusterResourceManager {
	if in == nil {
		return nil
	}
	out := new(ClusterResourceManager)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterResourceManager) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterResourceManagerList) DeepCopyInto(out *ClusterResourceManagerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ClusterResourceManager, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterResourceManagerList.
func (in *ClusterResourceManagerList) DeepCopy() *ClusterResourceManagerList {
	if in == nil {
		return nil
	}
	out := new(ClusterResourceManagerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterResourceManagerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterResourceManagerSpec) DeepCopyInto(out *ClusterResourceManagerSpec) {
	*out = *in
	in.ResourceManagerSpec.DeepCopyInto(&out.ResourceManagerSpec)
	if in.NamespaceSelector != nil {
		in, out := &in.NamespaceSelector, &out.NamespaceSelector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterResourceManagerSpec.
func (in *ClusterResourceManagerSpec) DeepCopy() *ClusterResourceManagerSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterResourceManagerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterResourceManagerStatus) DeepCopyInto(out *ClusterResourceManagerStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterResourceManagerStatus.
func (in *ClusterResourceManagerStatus) DeepCopy() *ClusterResourceManagerStatus {
	if in == nil {
		return nil
	}
	out := new(ClusterResourceManagerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Expiration) DeepCopyInto(out *Expiration) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Expiration.
func (in *Expiration) DeepCopy() *Expiration {
	if in == nil {
		return nil
	}
	out := new(Expiration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceManager) DeepCopyInto(out *ResourceManager) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceManager.
func (in *ResourceManager) DeepCopy() *ResourceManager {
	if in == nil {
		return nil
	}
	out := new(ResourceManager)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceManager) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceManagerList) DeepCopyInto(out *ResourceManagerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ResourceManager, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceManagerList.
func (in *ResourceManagerList) DeepCopy() *ResourceManagerList {
	if in == nil {
		return nil
	}
	out := new(ResourceManagerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceManagerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceManagerSpec) DeepCopyInto(out *ResourceManagerSpec) {
	*out = *in
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	out.Condition = in.Condition
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceManagerSpec.
func (in *ResourceManagerSpec) DeepCopy() *ResourceManagerSpec {
	if in == nil {
		return nil
	}
	out := new(ResourceManagerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceManagerStatus) DeepCopyInto(out *ResourceManagerStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceManagerStatus.
func (in *ResourceManagerStatus) DeepCopy() *ResourceManagerStatus {
	if in == nil {
		return nil
	}
	out := new(ResourceManagerStatus)
	in.DeepCopyInto(out)
	return out
}
