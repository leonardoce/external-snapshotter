package utils

import (
	"reflect"
	"testing"
)

var (
	firstInfoListEntry = SnapshotInfo{
		VolumeHandle: "b98b9100-fbe2-11ee-b405-b2139ff66f78",
		PVName:       "pvc-f1718d88-e548-480b-bee8-cbfc47faaf59",
		PVCName:      "cluster-example-1",
	}

	secondInfoListEntry = SnapshotInfo{
		VolumeHandle: "b98c3e51-fbe2-11ee-b405-b2139ff66f788",
		PVName:       "pvc-c59c9c0f-159a-43d6-9c60-61ccaf03158c",
		PVCName:      "cluster-example-1-wal",
	}

	infoList = SnapshotInfoList{firstInfoListEntry, secondInfoListEntry}
)

func TestMarshalUnmarshal(t *testing.T) {
	jsonContent, err := infoList.ToJSON()
	if err != nil {
		t.Fatalf("JSON serialization failed: %v", err)
	}

	unmarshalledInfoList, err := SnapshotInfoFromJSON(jsonContent)
	if err != nil {
		t.Fatalf("JSON deserialization failed: %v", err)
	}

	if !reflect.DeepEqual(infoList, unmarshalledInfoList) {
		t.Fatalf("unexpected info loss in serialization/deserialization: %v %v", infoList, unmarshalledInfoList)
	}
}

func TestEmptyMarshalUnmarshal(t *testing.T) {
	var emptyInfoList SnapshotInfoList = nil
	jsonContent, err := emptyInfoList.ToJSON()
	if err != nil {
		t.Fatalf("JSON serialization failed: %v", err)
	}

	unmarshalledInfoList, err := SnapshotInfoFromJSON(jsonContent)
	if err != nil {
		t.Fatalf("JSON deserialization failed: %v", err)
	}

	if !reflect.DeepEqual(emptyInfoList, unmarshalledInfoList) {
		t.Fatalf("unexpected info loss in serialization/deserialization: %v %v", infoList, unmarshalledInfoList)
	}
}

func TestGetFromVolumeHandle(t *testing.T) {
	testcases := []struct {
		volumeHandle string
		result       SnapshotInfo
		infoList     SnapshotInfoList
	}{
		{
			volumeHandle: firstInfoListEntry.VolumeHandle,
			result:       firstInfoListEntry,
			infoList:     infoList,
		},
		{
			volumeHandle: secondInfoListEntry.VolumeHandle,
			result:       secondInfoListEntry,
			infoList:     infoList,
		},
		{
			volumeHandle: "<non-existing>",
			result:       SnapshotInfo{},
			infoList:     infoList,
		},
		{
			volumeHandle: "<non-existing>",
			result:       SnapshotInfo{},
			infoList:     nil,
		},
	}

	for _, tc := range testcases {
		t.Logf("looking for %s:", tc.volumeHandle)
		result := tc.infoList.GetFromVolumeHandle(tc.volumeHandle)
		if !reflect.DeepEqual(result, tc.result) {
			t.Fatalf("unexpected GetFromVolumeHandle result: %v %v", result, tc.result)
		}
	}
}

func TestGetFromPVName(t *testing.T) {
	testcases := []struct {
		pvName   string
		result   SnapshotInfo
		infoList SnapshotInfoList
	}{
		{
			pvName:   firstInfoListEntry.PVName,
			result:   firstInfoListEntry,
			infoList: infoList,
		},
		{
			pvName:   secondInfoListEntry.PVName,
			result:   secondInfoListEntry,
			infoList: infoList,
		},
		{
			pvName:   "<non-existing>",
			result:   SnapshotInfo{},
			infoList: infoList,
		},
		{
			pvName:   "<non-existing>",
			result:   SnapshotInfo{},
			infoList: nil,
		},
	}

	for _, tc := range testcases {
		t.Logf("looking for %s:", tc.pvName)
		result := tc.infoList.GetFromPVName(tc.pvName)
		if !reflect.DeepEqual(result, tc.result) {
			t.Fatalf("unexpected GetFromPVName result: %v %v", result, tc.result)
		}
	}
}
