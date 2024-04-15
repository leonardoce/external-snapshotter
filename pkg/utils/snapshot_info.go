package utils

import (
	"encoding/json"
	"fmt"
)

// SnapshotInfo contains basic information about a volume being snapshotted
type SnapshotInfo struct {
	VolumeHandle string `json:"volumeHandle"`
	PVName       string `json:"pvName"`
	PVCName      string `json:"pvcName"`
}

// SnapshotInfoList contains basic information about a set of volumes being snapshotted
type SnapshotInfoList []SnapshotInfo

// ToJSON serizalizes to JSON a set of SnapshotInfo
func (data SnapshotInfoList) ToJSON() (string, error) {
	result, err := json.Marshal(data)
	if err != nil {
		err = fmt.Errorf("while serializing SnapshotInfoList: %w", err)
	}
	return string(result), err
}

// SnapshotInfoFromJSON deserializes from JSON a set of snapshot info
func SnapshotInfoFromJSON(content string) (SnapshotInfoList, error) {
	var result SnapshotInfoList

	err := json.Unmarshal([]byte(content), &result)
	if err != nil {
		err = fmt.Errorf("while de-serializing SnapshotInfoList: %w", err)
	}

	return result, err
}

// GetFromVolumeHandle gets the entry from the list corresponding to a certain
// volume handle. Returns an empty SnapshotInfo if there is no such entry
func (data SnapshotInfoList) GetFromVolumeHandle(volumeHandle string) SnapshotInfo {
	for i := range data {
		if data[i].VolumeHandle == volumeHandle {
			return data[i]
		}
	}

	return SnapshotInfo{}
}

// GetFromPVName gets the entry from the list corresponding to a certain
// PV name. Returns an empty SnapshotInfo if there is no such entry
func (data SnapshotInfoList) GetFromPVName(pvName string) SnapshotInfo {
	for i := range data {
		if data[i].PVName == pvName {
			return data[i]
		}
	}

	return SnapshotInfo{}
}
