/*
Copyright 2024 The Kubernetes Authors.

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

package utils

import (
	"fmt"

	crdv1beta1 "github.com/kubernetes-csi/external-snapshotter/client/v8/apis/volumegroupsnapshot/v1beta1"
	crdv1 "github.com/kubernetes-csi/external-snapshotter/client/v8/apis/volumesnapshot/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// VolumeSnapshotParentGroupIndex is the name of the cache index hosting the relationship
// between volume snapshots and their volume group snapshot owner
const VolumeSnapshotParentGroupIndex = "ByVolumeGroupSnapshotMembership"

// GetVolumeGroupSnapshotParentObjectName returns the name of the parent group snapshot, if present.
// The second return value is true when the parent object have been found, false otherwise.
func GetVolumeGroupSnapshotParentObjectName(snapshot *crdv1.VolumeSnapshot) (string, bool) {
	if snapshot == nil {
		return "", false
	}

	apiVersion := fmt.Sprintf(
		"%s/%s",
		crdv1beta1.SchemeGroupVersion.Group,
		crdv1beta1.SchemeGroupVersion.Version,
	)

	for _, owner := range snapshot.ObjectMeta.OwnerReferences {
		if owner.Kind == "VolumeGroupSnapshot" && owner.APIVersion == apiVersion {
			return owner.Name, true
		}
	}

	return "", false
}

// IsVolumeGroupSnapshotMember returns true if the passed VolumeSnapshot object
// is a member of a VolumeGroupSnapshot.
func IsVolumeGroupSnapshotMember(snapshot *crdv1.VolumeSnapshot) bool {
	_, parentFound := GetVolumeGroupSnapshotParentObjectName((snapshot))
	return parentFound
}

// VolumeSnapshotParentGroupKeyFunc maps a member snapshot to the name
// of the parent VolumeGroupSnapshot
func VolumeSnapshotParentGroupKeyFunc(snapshot *crdv1.VolumeSnapshot) string {
	parentName, found := GetVolumeGroupSnapshotParentObjectName(snapshot)
	if !found {
		return ""
	}

	return VolumeSnapshotParentGroupKeyFuncByComponents(types.NamespacedName{
		Namespace: snapshot.Namespace,
		Name:      parentName,
	})
}

// VolumeSnapshotParentGroupKeyFuncByComponents computes the index key for a certain
// name and namespace pair
func VolumeSnapshotParentGroupKeyFuncByComponents(objectKey types.NamespacedName) string {
	return fmt.Sprintf("%s^%s", objectKey.Namespace, objectKey.Name)
}

// NeedToAddVolumeGroupSnapshotOwnership checks if the passed snapshot is a member
// of a volume group snapshot but the ownership is missing
func NeedToAddVolumeGroupSnapshotOwnership(snapshot *crdv1.VolumeSnapshot) bool {
	if snapshot == nil {
		return false
	}

	if snapshot.Status == nil {
		return false
	}

	if snapshot.Status.VolumeGroupSnapshotName == nil {
		return false
	}

	if len(*snapshot.Status.VolumeGroupSnapshotName) == 0 {
		return false
	}

	apiVersion := fmt.Sprintf(
		"%s/%s",
		crdv1beta1.SchemeGroupVersion.Group,
		crdv1beta1.SchemeGroupVersion.Version,
	)

	for _, owner := range snapshot.ObjectMeta.OwnerReferences {
		if owner.Kind == "VolumeGroupSnapshot" && owner.APIVersion == apiVersion {
			return false
		}
	}

	return true
}

// CreateVolumeGroupSnapshotOwnerReference creates a OwnerReference record declaring an
// object as a child of passed VolumeGroupSnapshot
func CreateVolumeGroupSnapshotOwnerReference(parentGroup *crdv1beta1.VolumeGroupSnapshot) metav1.OwnerReference {
	return metav1.OwnerReference{
		APIVersion: fmt.Sprintf(
			"%s/%s",
			crdv1beta1.SchemeGroupVersion.Group,
			crdv1beta1.SchemeGroupVersion.Version,
		),
		Kind: "VolumeGroupSnapshot",
		Name: parentGroup.Name,
		UID:  parentGroup.UID,
	}
}
