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
