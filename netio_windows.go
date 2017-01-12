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
	_, _, err := syscall.Syscall(ngetIPForwardTable2, 1, uintptr(unsafe.Pointer(table)), 0, 0)
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
		family := table.Table[i].DestinationPrefix.Prefix.family
		saData := table.Table[i].DestinationPrefix.Prefix.data
		gwData := table.Table[i].NextHop.data
		length := int(table.Table[i].DestinationPrefix.PrefixLength)
		metric := int(table.Table[i].Metric)
		var ip []byte
		var gateway []byte
		var mask net.IPMask
		if family == nAF_INET {
			ip = make([]byte, net.IPv4len)
			copy(ip, saData[2:])
			mask = net.CIDRMask(length, net.IPv4len*8)
			gateway = make([]byte, net.IPv4len)
			copy(gateway, gwData[2:])
		} else if family == nAF_INET6 {
			ip = make([]byte, net.IPv6len)
			copy(ip, saData[6:])
			mask = net.CIDRMask(length, net.IPv6len*8)
			gateway = make([]byte, net.IPv6len)
			copy(gateway, gwData[6:])
		} else {
			continue
		}
		ipnet := &net.IPNet{IP: ip, Mask: mask}
		ifce, _ := net.InterfaceByIndex(int(table.Table[i].InterfaceIndex))
		isDefault := false
		if length == 0 {
			isDefault = true
		}
		route := Route{ipnet, ifce, gateway, isDefault, metric}
		routes = append(routes, route)
	}
	return routes, nil
}
