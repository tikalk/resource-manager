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
func (in *Condition) DeepCopyInto(out *Condition) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Condition.
func (in *Condition) DeepCopy() *Condition {
	if in == nil {
		return nil
	}
	out := new(Condition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExpiryCondition) DeepCopyInto(out *ExpiryCondition) {
	*out = *in
	out.Condition = in.Condition
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExpiryCondition.
func (in *ExpiryCondition) DeepCopy() *ExpiryCondition {
	if in == nil {
		return nil
	}
	out := new(ExpiryCondition)
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
	if in.Condition != nil {
		in, out := &in.Condition, &out.Condition
		*out = make([]ExpiryCondition, len(*in))
		copy(*out, *in)
	}
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
