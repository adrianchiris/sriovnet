//go:build integration
// +build integration

/*
Copyright 2023 NVIDIA CORPORATION & AFFILIATES

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

/*
NOTE:
====
This test although not usable as is, since netdev names and configuration differs from one setup to the next
is useful for testing your changes on a real setup.

for new functionality in sriovnet package add a new integration test case.
to run an existing test modify the netdev/VF in the set to fit your setup and execute the test.

Build and run integration test:
==============================
# go test --tags integration -v -run <TestName>
*/

package sriovnet

import (
	"fmt"
	"testing"
)

func TestEnableSriov(t *testing.T) {

	err := EnableSriov("ib0")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDisableSriov(t *testing.T) {
	err := DisableSriov("ib0")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetPfHandle(t *testing.T) {
	err1 := EnableSriov("ib0")
	if err1 != nil {
		t.Fatal(err1)
	}

	handle, err2 := GetPfNetdevHandle("ib0")
	if err2 != nil {
		t.Fatal(err2)
	}
	for _, vf := range handle.List {
		fmt.Printf("vf = %v\n", vf)
	}
}

func TestConfigVfs(t *testing.T) {
	err1 := EnableSriov("ens2f0")
	if err1 != nil {
		t.Fatal(err1)
	}

	handle, err2 := GetPfNetdevHandle("ens2f0")
	if err2 != nil {
		t.Fatal(err2)
	}
	err3 := ConfigVfs(handle, false)
	if err3 != nil {
		t.Fatal(err3)
	}
	for _, vf := range handle.List {
		fmt.Printf("after config vf = %v\n", vf)
	}
}

func TestIsSriovEnabled(t *testing.T) {
	status := IsSriovEnabled("ens2f0")

	fmt.Printf("sriov status = %v", status)
}

func TestAllocFreeVf(t *testing.T) {
	var vfList [10]*VfObj

	err1 := EnableSriov("ib0")
	if err1 != nil {
		t.Fatal(err1)
	}

	handle, err2 := GetPfNetdevHandle("ib0")
	if err2 != nil {
		t.Fatal(err2)
	}
	err3 := ConfigVfs(handle, false)
	if err3 != nil {
		t.Fatal(err3)
	}
	for i := 0; i < 10; i++ {
		vfList[i], _ = AllocateVf(handle)
	}
	for _, vf := range handle.List {
		fmt.Printf("after allocation vf = %v\n", vf)
	}
	for i := 0; i < 10; i++ {
		if vfList[i] == nil {
			continue
		}
		FreeVf(handle, vfList[i])
	}
	for _, vf := range handle.List {
		fmt.Printf("after free vf = %v\n", vf)
	}
}

func TestFreeByName(t *testing.T) {
	var vfList [10]*VfObj

	err1 := EnableSriov("ib0")
	if err1 != nil {
		t.Fatal(err1)
	}

	handle, err2 := GetPfNetdevHandle("ib0")
	if err2 != nil {
		t.Fatal(err2)
	}
	err3 := ConfigVfs(handle, false)
	if err3 != nil {
		t.Fatal(err3)
	}
	for i := 0; i < 10; i++ {
		vfList[i], _ = AllocateVf(handle)
	}
	for _, vf := range handle.List {
		fmt.Printf("after allocation vf = %v\n", vf)
	}
	for i := 0; i < 10; i++ {
		if vfList[i] == nil {
			continue
		}
		err4 := FreeVfByNetdevName(handle, vfList[i].Index)
		if err4 != nil {
			t.Fatal(err4)
		}
	}
	for _, vf := range handle.List {
		fmt.Printf("after free vf = %v\n", vf)
	}
}

func TestAllocateVfByMac(t *testing.T) {
	var vfList [10]*VfObj
	var vfName [10]string

	err1 := EnableSriov("ens2f0")
	if err1 != nil {
		t.Fatal(err1)
	}

	handle, err2 := GetPfNetdevHandle("ens2f0")
	if err2 != nil {
		t.Fatal(err2)
	}
	err3 := ConfigVfs(handle, true)
	if err3 != nil {
		t.Fatal(err3)
	}
	for i := 0; i < 10; i++ {
		vfList[i], _ = AllocateVf(handle)
		if vfList[i] != nil {
			vfName[i] = GetVfNetdevName(handle, vfList[i])
		}
	}
	for _, vf := range handle.List {
		fmt.Printf("after allocation vf = %v\n", vf)
	}
	for i := 0; i < 10; i++ {
		if vfList[i] == nil {
			continue
		}
		FreeVf(handle, vfList[i])
	}
	for _, vf := range handle.List {
		fmt.Printf("after alloc vf = %v\n", vf)
	}
	for i := 0; i < 2; i++ {
		if vfName[i] == "" {
			continue
		}
		mac, _ := GetVfDefaultMacAddr(vfName[i])
		vfList[i], _ = AllocateVfByMacAddress(handle, mac)
	}
	for _, vf := range handle.List {
		fmt.Printf("after alloc vf = %v\n", vf)
	}
}

func TestGetVfPciDevList(t *testing.T) {

	list, _ := GetVfPciDevList("ens2f0")
	fmt.Println("list is: ", list)
	t.Fatal(nil)
}

func TestGetVfNetdevName(t *testing.T) {
	var vfList [10]*VfObj
	var vfName [10]string

	handle, err2 := GetPfNetdevHandle("ens2f0")
	if err2 != nil {
		t.Fatal(err2)
	}
	for i := 0; i < 10; i++ {
		vfList[i], _ = AllocateVf(handle)
		if vfList[i] != nil {
			vfName[i] = GetVfNetdevName(handle, vfList[i])
			t.Log("Allocated VF: ", vfList[i].Index, "Netdev: ", vfName[i])
		}
	}
}

func TestIntegrationGetPfPciFromVfPci(t *testing.T) {
	vf := "0000:05:00.6"
	pf, err := GetPfPciFromVfPci(vf)
	if err != nil {
		t.Log("GetPfPciFromVfPci", "VF PCI: ", vf, "Error: ", err)
		t.Fatal()
	}
	t.Log("VF: ", vf, "PF: ", pf)
}

func TestIntegrationGetVfRepresentorSmartNIC(t *testing.T) {
	pfID := "0"
	vfIdx := "2"
	t.Log("GetVfRepresentorDPU ", "PF ID: ", pfID, "VF Index: ", vfIdx)
	rep, err := GetVfRepresentorDPU(pfID, vfIdx)
	if err != nil {
		t.Log("GetVfRepresentorDPU ", "Error: ", err)
		t.Fatal()
	}
	t.Log("VF Representor: ", rep)
}

func TestIntegrationGetSfRepresentorSmartNIC(t *testing.T) {
	pfID := "0"
	sfIdx := "1"
	t.Log("GetSfRepresentorDPU ", "PF ID: ", pfID, "SF Index: ", sfIdx)
	rep, err := GetSfRepresentorDPU(pfID, sfIdx)
	if err != nil {
		t.Log("GetSfRepresentorDPU ", "Error: ", err)
		t.Fatal()
	}
	t.Log("SF Representor: ", rep)
}

func TestIntegrationGetRepresentorPortFlavour(t *testing.T) {
	tcases := []struct {
		netdev          string
		expectedFlavour PortFlavour
		shouldFail      bool
	}{
		{netdev: "p0", expectedFlavour: PORT_FLAVOUR_PHYSICAL},
		{netdev: "pf0hpf", expectedFlavour: PORT_FLAVOUR_PCI_PF},
		{netdev: "pf0vf4", expectedFlavour: PORT_FLAVOUR_PCI_VF},
		{netdev: "fooBar", expectedFlavour: PORT_FLAVOUR_UNKNOWN, shouldFail: true},
	}

	for _, tcase := range tcases {
		flava, err := GetRepresentorPortFlavour(tcase.netdev)
		if tcase.shouldFail == true && err == nil {
			t.Fatal("Expected failure but no error occured")
		}
		if err != nil {
			t.Log("GetRepresentorPortFlavour got error", err)
		}
		if flava != tcase.expectedFlavour {
			t.Fatal("Actual flavour does not match expected flavour", flava, "!=", tcase.expectedFlavour)
		}
		t.Log("GetRepresentorPortFlavour", "netdev: ", tcase.netdev, "flavour: ", flava)
	}
}

func TestIntegrationGetRepresentorPeerMacAddress(t *testing.T) {
	tcases := []struct {
		netdev      string
		expectedMac string
		shouldFail  bool
	}{
		{netdev: "pf0hpf", expectedMac: "0c:42:a1:de:cf:7c", shouldFail: false},
		{netdev: "p0", expectedMac: "", shouldFail: true},
		{netdev: "pf0vf4", expectedMac: "", shouldFail: true},
		{netdev: "fooBar", expectedMac: "", shouldFail: true},
	}

	for _, tcase := range tcases {
		mac, err := GetRepresentorPeerMacAddress(tcase.netdev)
		if tcase.shouldFail {
			if err == nil {
				t.Fatal("Expected failure but no error occured")
			}
			continue
		}
		if err != nil {
			t.Fatal("GetRepresentorPeerMacAddress failed with error: ", err)
		}
		if mac.String() != tcase.expectedMac {
			t.Fatal("Actual MAC does not match expected MAC", mac, "!=", tcase.expectedMac)
		}
		t.Log("GetRepresentorPeerMacAddress", "netdev: ", tcase.netdev, "Mac: ", mac)
	}
}

func TestIntegrationGetVfRepresentor(t *testing.T) {
	tcases := []struct {
		uplink     string
		vfIndex    int
		expected   string
		shouldFail bool
	}{
		{uplink: "enp3s0f0", vfIndex: 2, expected: "enp3s0f0_2", shouldFail: false},
		{uplink: "foobar", vfIndex: 2, expected: "", shouldFail: true},
		{uplink: "enp3s0", vfIndex: 44, expected: "", shouldFail: true},
	}

	for _, tcase := range tcases {
		rep, err := GetVfRepresentor(tcase.uplink, tcase.vfIndex)
		if tcase.shouldFail {
			if err == nil {
				t.Fatal("Expected failure but no error occured")
			}
			continue
		}
		if err != nil {
			t.Fatal("GetVfRepresentor failed with error: ", err)
		}
		if rep != tcase.expected {
			t.Fatal("Actual Representor does not match expected Representor", rep, "!=", tcase.expected)
		}
		t.Log("GetVfRepresentor", "uplink: ", tcase.uplink, " VF Index: ", tcase.vfIndex, " Rep: ", rep)
	}
}

func TestIntegrationGetSfRepresentor(t *testing.T) {
	tcases := []struct {
		uplink     string
		sfNum      int
		expected   string
		shouldFail bool
	}{
		{uplink: "p0", sfNum: 2, expected: "en3f0pf0sf2", shouldFail: false},
		{uplink: "foobar", sfNum: 2, expected: "", shouldFail: true},
		{uplink: "enp3s0", sfNum: 44, expected: "", shouldFail: true},
	}

	for _, tcase := range tcases {
		sfRep, err := GetSfRepresentor(tcase.uplink, tcase.sfNum)
		if tcase.shouldFail {
			if err == nil {
				t.Fatal("Expected failure but no error occured")
			}
			continue
		}
		if err != nil {
			t.Fatal("GetSfRepresentor failed with error: ", err)
		}
		if sfRep != tcase.expected {
			t.Fatal("Actual Representor does not match expected Representor", sfRep, "!=", tcase.expected)
		}
		t.Log("GetSfRepresentor", "uplink: ", tcase.uplink, " SF Index: ", tcase.sfNum, " Rep: ", sfRep)
	}
}

func TestIntegrationGetUplinkRepresentor(t *testing.T) {
	uplink := "enp3s0f0np0"
	pfPciAddress := "0000:03:00.0"
	vfPciAddress := "0000:03:00.2"

	rep, err := GetUplinkRepresentor(pfPciAddress)

	if err != nil {
		t.Fatal("GetUplinkRepresentor failed with error: ", err)
	}

	if rep != uplink {
		t.Fatal("Actual Representor does not match expected Representor", rep, "!=", uplink)
	}

	rep, err = GetUplinkRepresentor(vfPciAddress)

	if err != nil {
		t.Fatal("GetUplinkRepresentor failed with error: ", err)
	}

	if rep != uplink {
		t.Fatal("Actual Representor does not match expected Representor", rep, "!=", uplink)
	}
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestIntegrationGetAuxNetDevicesFromPci(t *testing.T) {
	tcases := []struct {
		addr       string
		expected   []string
		shouldFail bool
	}{
		{addr: "0000:3b:00.0", expected: []string{"mlx5_core.eth.0", "mlx5_core.rdma.0"}, shouldFail: false},
		{addr: "0000:01:00.0", expected: []string{}, shouldFail: false},
		{addr: "0000:00:00.0", expected: []string(nil), shouldFail: true},
		{addr: "c0fe:00:00.0", expected: []string(nil), shouldFail: true},
	}

	for _, tcase := range tcases {
		devs, err := GetAuxNetDevicesFromPci(tcase.addr)
		if tcase.shouldFail {
			if err == nil {
				t.Fatal("Expected failure but no error occured")
			}
			continue
		}
		if err != nil {
			t.Fatal("GetAuxNetDevicesFromPci failed with error: ", err)
		}
		if !stringSlicesEqual(devs, tcase.expected) {
			t.Fatal("Actual devices does not match expected devices list", devs, "!=", tcase.expected)
		}
		t.Log("GetAuxNetDevicesFromPci", "PCI address: ", tcase.addr, " Devices: ", devs)
	}
}
