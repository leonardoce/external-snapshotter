/*
Copyright 2019 The Kubernetes Authors.

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

package common_controller

import (
	"testing"

	crdv1 "github.com/kubernetes-csi/external-snapshotter/client/v8/apis/volumesnapshot/v1"
	"github.com/kubernetes-csi/external-snapshotter/v8/pkg/utils"
	v1 "k8s.io/api/core/v1"
)

func TestDeleteGroupSnapshotSync(t *testing.T) {
	tests := []controllerTest{
		{
			name: "6-1 - group snapshot have been deleted, but no content was present - no op",
			initialGroupSnapshots: newGroupSnapshotArray(
				"snap-1-1", "snapuid1-1", map[string]string{
					"app.kubernetes.io/name": "postgresql",
				},
				"", classGold, "groupsnapcontent-snapuid1-1", &False, nil, nil, false, false, &timeNowMetav1,
			),
			expectedGroupSnapshots: newGroupSnapshotArray(
				"snap-1-1", "snapuid1-1", map[string]string{
					"app.kubernetes.io/name": "postgresql",
				},
				"", classGold, "groupsnapcontent-snapuid1-1", &False, nil, nil, false, false, &timeNowMetav1,
			),
			initialGroupContents:  nil,
			expectedGroupContents: nil,
			initialClaims: withClaimLabels(
				newClaimCoupleArray("claim1-1", "pvc-uid6-1", "1Gi", "volume6-1", v1.ClaimBound, &classGold),
				map[string]string{
					"app.kubernetes.io/name": "postgresql",
				}),
			initialVolumes: newVolumeCoupleArray("volume6-1", "pv-uid6-1", "pv-handle6-1", "1Gi", "pvc-uid6-1", "claim1-1", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classGold),
			errors:         noerrors,
			test:           testSyncGroupSnapshot,
			expectSuccess:  true,
		},
		{
			name: "6-2 - dynamic group snapshot have been deleted, retention policy set to retain - set 'being deleted' annotation",
			initialGroupSnapshots: withGroupSnapshotFinalizers(
				newGroupSnapshotArray(
					"snap-1-1", "snapuid1-1", map[string]string{
						"app.kubernetes.io/name": "postgresql",
					},
					"", classGold, "groupsnapcontent-snapuid1-1", &False, nil, nil, false, false, &timeNowMetav1,
				),
				utils.VolumeGroupSnapshotBoundFinalizer,
			),
			expectedGroupSnapshots: newGroupSnapshotArray(
				"snap-1-1", "snapuid1-1", map[string]string{
					"app.kubernetes.io/name": "postgresql",
				},
				"", classGold, "groupsnapcontent-snapuid1-1", &False, nil, nil, false, false, &timeNowMetav1,
			),
			initialGroupContents: newGroupSnapshotContentArray(
				"groupsnapcontent-snapuid1-1", "snapuid1-1", "snap-1-1", "snapshot-handle", classGold, []string{
					"1-pv-handle6-1",
					"2-pv-handle6-1",
				}, "", crdv1.VolumeSnapshotContentRetain, nil, false, false,
			),
			expectedGroupContents: withGroupContentAnnotations(
				newGroupSnapshotContentArray(
					"groupsnapcontent-snapuid1-1", "snapuid1-1", "snap-1-1", "snapshot-handle", classGold, []string{
						"1-pv-handle6-1",
						"2-pv-handle6-1",
					}, "", crdv1.VolumeSnapshotContentRetain, nil, false, false,
				),
				map[string]string{
					utils.AnnVolumeGroupSnapshotBeingDeleted: "yes",
				},
			),
			initialClaims: withClaimLabels(
				newClaimCoupleArray("claim1-1", "pvc-uid6-1", "1Gi", "volume6-1", v1.ClaimBound, &classGold),
				map[string]string{
					"app.kubernetes.io/name": "postgresql",
				}),
			initialVolumes: newVolumeCoupleArray("volume6-1", "pv-uid6-1", "pv-handle6-1", "1Gi", "pvc-uid6-1", "claim1-1", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classGold),
			errors:         noerrors,
			test:           testSyncGroupSnapshot,
			expectSuccess:  true,
		},
		{
			name: "6-3 - dynamic group snapshot have been deleted, retention policy set to delete - volume snapshot content will be deleted",
			initialGroupSnapshots: withGroupSnapshotFinalizers(
				newGroupSnapshotArray(
					"snap-1-1", "snapuid1-1", map[string]string{
						"app.kubernetes.io/name": "postgresql",
					},
					"", classGold, "groupsnapcontent-snapuid1-1", &False, nil, nil, false, false, &timeNowMetav1,
				),
				utils.VolumeGroupSnapshotBoundFinalizer,
			),
			expectedGroupSnapshots: withGroupSnapshotFinalizers(
				newGroupSnapshotArray(
					"snap-1-1", "snapuid1-1", map[string]string{
						"app.kubernetes.io/name": "postgresql",
					},
					"", classGold, "groupsnapcontent-snapuid1-1", &False, nil, nil, false, false, &timeNowMetav1,
				),
				utils.VolumeGroupSnapshotBoundFinalizer,
			),
			initialGroupContents: newGroupSnapshotContentArray(
				"groupsnapcontent-snapuid1-1", "snapuid1-1", "snap-1-1", "snapshot-handle", classGold, []string{
					"1-pv-handle6-1",
					"2-pv-handle6-1",
				}, "", crdv1.VolumeSnapshotContentDelete, nil, false, false,
			),
			expectedGroupContents: nil,
			initialClaims: withClaimLabels(
				newClaimCoupleArray("claim1-1", "pvc-uid6-1", "1Gi", "volume6-1", v1.ClaimBound, &classGold),
				map[string]string{
					"app.kubernetes.io/name": "postgresql",
				}),
			initialVolumes: newVolumeCoupleArray("volume6-1", "pv-uid6-1", "pv-handle6-1", "1Gi", "pvc-uid6-1", "claim1-1", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classGold),
			errors:         noerrors,
			test:           testSyncGroupSnapshot,
			expectSuccess:  true,
		},
	}
	runSyncTests(t, tests, nil, groupSnapshotClasses)
}