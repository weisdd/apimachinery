/*
Copyright The Kmodules Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KEP: https://github.com/kubernetes/enhancements/blob/ced773ab59f0ff080888a912ab99474245623dad/keps/sig-api-machinery/1623-standardize-conditions/README.md

type ConditionStatus string

// These are valid condition statuses. "ConditionTrue" means a resource is in the condition.
// "ConditionFalse" means a resource is not in the condition. "ConditionUnknown" means kubernetes
// can't decide if a resource is in the condition or not. In the future, we could add other
// intermediate conditions, e.g. ConditionDegraded.
const (
	ConditionTrue    ConditionStatus = "True"
	ConditionFalse   ConditionStatus = "False"
	ConditionUnknown ConditionStatus = "Unknown"
)

type Condition struct {
	// Type of condition in CamelCase or in foo.example.com/CamelCase.
	// Many .condition.type values are consistent across resources like Available, but because arbitrary conditions can be
	// useful (see .node.status.conditions), the ability to deconflict is important.
	// +required
	Type string `json:"type" protobuf:"bytes,1,opt,name=type"`
	// Status of the condition, one of True, False, Unknown.
	// +required
	Status ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status"`
	// If set, this represents the .metadata.generation that the condition was set based upon.
	// For instance, if .metadata.generation is currently 12, but the .status.condition[x].observedGeneration is 9, the condition is out of date
	// with respect to the current state of the instance.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,3,opt,name=observedGeneration"`
	// Last time the condition transitioned from one status to another.
	// This should be when the underlying condition changed.  If that is not known, then using the time when the API field changed is acceptable.
	// +required
	LastTransitionTime metav1.Time `json:"lastTransitionTime" protobuf:"bytes,4,opt,name=lastTransitionTime"`
	// The reason for the condition's last transition in CamelCase.
	// The specific API may choose whether or not this field is considered a guaranteed API.
	// This field may not be empty.
	// +required
	Reason string `json:"reason" protobuf:"bytes,5,opt,name=reason"`
	// A human readable message indicating details about the transition.
	// This field may be empty.
	// +required
	Message string `json:"message" protobuf:"bytes,6,opt,name=message"`
}

// HasCondition returns "true" if the desired condition provided in "condType" is present in the condition list.
// Otherwise, it returns "false".
func HasCondition(conditions []Condition, condType string) bool {
	for i := range conditions {
		if conditions[i].Type == condType {
			return true
		}
	}
	return false
}

// GetCondition returns a pointer to the desired condition referred by "condType". Otherwise, it returns nil.
func GetCondition(conditions []Condition, condType string) (int, *Condition) {
	for i := range conditions {
		c := conditions[i]
		if c.Type == condType {
			return i, &c
		}
	}
	return -1, nil
}

// SetCondition add/update the desired condition to the condition list. It does nothing if the condition is already in
// its desired state.
func SetCondition(conditions []Condition, newCondition Condition) []Condition {
	idx, curCond := GetCondition(conditions, newCondition.Type)
	// If the current condition is in its desired state, we have nothing to do. Just return the original condition list.
	if curCond != nil &&
		curCond.Status == newCondition.Status &&
		curCond.Reason == newCondition.Reason &&
		curCond.Message == newCondition.Message &&
		curCond.ObservedGeneration == newCondition.ObservedGeneration {
		return conditions
	}
	// The desired conditions is not in the condition list or is not in its desired state.
	// Update it if present in the condition list, or append the new condition if it does not present.
	newCondition.LastTransitionTime = metav1.Now()
	if idx != -1 {
		conditions[idx] = newCondition
	} else {
		conditions = append(conditions, newCondition)
	}
	return conditions
}

// RemoveCondition remove a condition from the condition list referred by "condType" parameter.
func RemoveCondition(conditions []Condition, condType string) []Condition {
	idx, _ := GetCondition(conditions, condType)
	if idx == -1 {
		// The desired condition is not present in the condition list. So, nothing to do.
		return conditions
	}
	return append(conditions[:idx], conditions[idx+1:]...)
}

// IsConditionTrue returns "true" if the desired condition is in true state.
// It returns "false" if the desired condition is not in "true" state or is not in the condition list.
func IsConditionTrue(conditions []Condition, condType string) bool {
	for i := range conditions {
		if conditions[i].Type == condType && conditions[i].Status == ConditionTrue {
			return true
		}
	}
	return false
}
