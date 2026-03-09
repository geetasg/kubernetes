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

package structured

import (
	resourceapi "k8s.io/api/resource/v1"
	"k8s.io/utils/ptr"
)

// SliceIndex provides O(1) lookup of ResourceSlices by node name.
// Built once per scheduling cycle in PreFilter, used per-node in Filter
// to avoid scanning all slices in GatherPools.
type SliceIndex struct {
	byNodeName map[string][]*resourceapi.ResourceSlice
	global     []*resourceapi.ResourceSlice
	all        []*resourceapi.ResourceSlice
}

// NewSliceIndex partitions slices into per-node and global buckets.
func NewSliceIndex(slices []*resourceapi.ResourceSlice) *SliceIndex {
	idx := &SliceIndex{
		byNodeName: make(map[string][]*resourceapi.ResourceSlice, len(slices)),
		all:        slices,
	}
	for _, slice := range slices {
		if nodeName := ptr.Deref(slice.Spec.NodeName, ""); nodeName != "" {
			idx.byNodeName[nodeName] = append(idx.byNodeName[nodeName], slice)
		} else {
			idx.global = append(idx.global, slice)
		}
	}
	return idx
}

// SlicesForNode returns slices relevant to a node: those targeting it by
// name plus any global slices requiring per-node evaluation.
func (idx *SliceIndex) SlicesForNode(nodeName string) []*resourceapi.ResourceSlice {
	nodeSlices := idx.byNodeName[nodeName]
	if len(idx.global) == 0 {
		return nodeSlices
	}
	result := make([]*resourceapi.ResourceSlice, 0, len(nodeSlices)+len(idx.global))
	result = append(result, nodeSlices...)
	result = append(result, idx.global...)
	return result
}

// All returns the original unindexed slice list.
func (idx *SliceIndex) All() []*resourceapi.ResourceSlice {
	return idx.all
}
