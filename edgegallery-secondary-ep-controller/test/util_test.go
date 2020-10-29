/*
 *  Copyright 2020 Huawei Technologies Co., Ltd.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package test

import (
	"edgegallery-secondary-ep-controller/watcher"
	"errors"
	"github.com/agiledragon/gomonkey"
	"testing"

	"github.com/intel/multus-cni/types"
	"github.com/smartystreets/goconvey/convey"
)

func TestIsInNetworkSelectionElementsArray(t *testing.T) {
	convey.Convey("Testing network selection array", t, func() {
		isOk := watcher.IsInNetworkSelectionElementsArray("", nil)
		if isOk {
			t.Error("TestCase failed")
		}
	})
}
func TestIsInNetworkSelectionElementsArraySuccess(t *testing.T) {
	convey.Convey("Testing network selection array success", t, func() {
		var networkSelections []*types.NetworkSelectionElement
		networkSelectionElement := &types.NetworkSelectionElement{
			Namespace: "default",
			Name:      "mp1",
		}
		networkSelections = append(networkSelections, networkSelectionElement)
		isOk := watcher.IsInNetworkSelectionElementsArray("mp1", networkSelections)
		if !isOk {
			t.Error("TestCase failed")
		}
	})
}

func TestIsInNetworkSelectionElementsArrayFailure(t *testing.T) {
	convey.Convey("Testing network selection array success", t, func() {
		var networkSelections []*types.NetworkSelectionElement
		patch1 := gomonkey.ApplyFunc(watcher.ResolveNetworkAnnotation, func(string, string) (string, string, error) {
			return "", "", errors.New("some error")
		})
		networkSelectionElement := &types.NetworkSelectionElement{
			Namespace: "default",
			Name:      "mp1",
		}
		networkSelections = append(networkSelections, networkSelectionElement)
		isOk := watcher.IsInNetworkSelectionElementsArray("mp1", networkSelections)
		if isOk {
			t.Error("TestCase failed")
		}
		defer patch1.Reset()
	})
}
