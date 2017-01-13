// +build windows

package network

import (
	"fmt"
	"net"
	"syscall"
	"unsafe"
)

var (
	ngetIPForwardTable2 uintptr
	nFreeMibTable       uintptr
)

func getProcAddr(lib syscall.Handle, name string) uintptr {
	addr, err := syscall.GetProcAddress(lib, name)
	if err != nil {
		panic(name + " " + err.Error())
	}
	return addr
}

func getIPForwardTable2(family nAddressFamily, table **nMIBIpforwardTable2) error {
	r, _, err := syscall.Syscall(ngetIPForwardTable2, 2, uintptr(family), uintptr(unsafe.Pointer(table)), 0)
	if r == 0 {
		return nil
	}
	return err
}

func freeMibTable(table *nMIBIpforwardTable2) error {
	_, _, err := syscall.Syscall(nFreeMibTable, 1, uintptr(unsafe.Pointer(table)), 0, 0)
	return err
}

func getRoutes() ([]Route, error) {
	iph, err := syscall.LoadLibrary("Iphlpapi.dll")
	if err != nil {
		panic("LoadLibrary " + err.Error())
	}
	defer syscall.FreeLibrary(iph)
	ngetIPForwardTable2 = getProcAddr(iph, "GetIpForwardTable2")
	nFreeMibTable = getProcAddr(iph, "FreeMibTable")
	table := &nMIBIpforwardTable2{}
	if err := getIPForwardTable2(nAF_UNSPEC, &table); err != nil {
		fmt.Println(err)
	}
	defer freeMibTable(table)
	routes := make([]Route, 0, int(table.NumEntries))
	for i := 0; i < int(table.NumEntries); i++ {
		if route := table.Table[i].ToRoute(); route != nil {
			routes = append(routes, *route)
		}
	}
	return routes, nil
}
