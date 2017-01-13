// +build windows

package network

import (
	"fmt"
	"net"
	"syscall"
	"unsafe"
)

var (
	errInvalidAddr = fmt.Errorf("Provided address is invalid")
)

var (
	nGetIPForwardTable2,
	nFreeMibTable,
	nCreateIPForwardEntry2,
	nSetIPForwardEntry2,
	nDeleteIPForwardEntry2,
	nGetBestInterfaceEx,
	nGetBestRoute2 uintptr
)

func getProcAddr(lib syscall.Handle, name string) uintptr {
	addr, err := syscall.GetProcAddress(lib, name)
	if err != nil {
		panic(name + " " + err.Error())
	}
	return addr
}

func getBestInterfaceEx(dstAddr *nRawSockAddrAny, ifIndex *nNetIFIndex) error {
	// TODO: test the return value
	r, _, err := syscall.Syscall(nGetBestInterfaceEx, 2, uintptr(unsafe.Pointer(dstAddr)), uintptr(unsafe.Pointer(ifIndex)), 0)
	if r == 0 {
		return nil
	}
	return err
}

func getBestRoute2(ifLuid *nNetLuid, ifIndex nNetIFIndex, srcAddr *nRawSockAddrAny, dstAddr *nRawSockAddrAny, sortOpt uint32, bestRoute *nMIBIpforwardRow2, bestSrcAddr *nRawSockAddrAny) error {
	// TODO: test the return value
	r, _, err := syscall.Syscall9(nGetBestRoute2, 7, uintptr(unsafe.Pointer(ifLuid)), uintptr(ifIndex),
		uintptr(unsafe.Pointer(srcAddr)), uintptr(unsafe.Pointer(dstAddr)), uintptr(sortOpt),
		uintptr(unsafe.Pointer(bestRoute)), uintptr(unsafe.Pointer(bestSrcAddr)), 0, 0)
	if r == 0 {
		return nil
	}
	return err
}

func getIPForwardTable2(family nAddressFamily, table **nMIBIpforwardTable2) error {
	r, _, err := syscall.Syscall(nGetIPForwardTable2, 2, uintptr(family), uintptr(unsafe.Pointer(table)), 0)
	if r == 0 {
		return nil
	}
	return err
}

func createIPForwardEntry2(row *nMIBIpforwardRow2) error {
	// TODO: test the return value
	r, _, err := syscall.Syscall(nCreateIPForwardEntry2, 1, uintptr(unsafe.Pointer(row)), 0, 0)
	if r == 0 {
		return nil
	}
	return err
}

func setIPForwardEntry2(row *nMIBIpforwardRow2) error {
	// TODO: test the return value
	r, _, err := syscall.Syscall(nSetIPForwardEntry2, 1, uintptr(unsafe.Pointer(row)), 0, 0)
	if r == 0 {
		return nil
	}
	return err
}

func deleteIPForwardEntry2(row *nMIBIpforwardRow2) error {
	// TODO: test the return value
	r, _, err := syscall.Syscall(nDeleteIPForwardEntry2, 1, uintptr(unsafe.Pointer(row)), 0, 0)
	if r == 0 {
		return nil
	}
	return err
}

func freeMibTable(table *nMIBIpforwardTable2) error {
	_, _, err := syscall.Syscall(nFreeMibTable, 1, uintptr(unsafe.Pointer(table)), 0, 0)
	return err
}

func getRoute(ip net.IP) (*Route, error) {
	// TODO
	// return fmt.Errorf("Not implemented yet")
	iph, err := syscall.LoadLibrary("Iphlpapi.dll")
	if err != nil {
		panic("LoadLibrary " + err.Error())
	}
	defer syscall.FreeLibrary(iph)
	nGetBestInterfaceEx = getProcAddr(iph, "GetBestInterfaceEx")
	nGetBestRoute2 = getProcAddr(iph, "GetBestRoute2")

	var ifIndex nNetIFIndex
	dstAddr := ParseIP(ip)
	if dstAddr == nil {
		return nil, errInvalidAddr
	}
	if err := getBestInterfaceEx(dstAddr, &ifIndex); err != nil {
		// TODO: parse err
		return nil, err
	}

	row := newMIBIpforwardRow2()
	src := nRawSockAddrAny{}
	if err := getBestRoute2(nil, ifIndex, nil, dstAddr, 0, row, &src); err != nil {
		// TODO: parse err
		return nil, err
	}
	return row.ToRoute(), nil
}

func getRoutes() ([]Route, error) {
	iph, err := syscall.LoadLibrary("Iphlpapi.dll")
	if err != nil {
		panic("LoadLibrary " + err.Error())
	}
	defer syscall.FreeLibrary(iph)
	nGetIPForwardTable2 = getProcAddr(iph, "GetIpForwardTable2")
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

func createRoute(route *Route) error {
	// TODO
	return fmt.Errorf("Not implemented yet")
	iph, err := syscall.LoadLibrary("Iphlpapi.dll")
	if err != nil {
		panic("LoadLibrary " + err.Error())
	}
	defer syscall.FreeLibrary(iph)
	nCreateIPForwardEntry2 = getProcAddr(iph, "CreateIpForwardEntry2")

	row := newMIBIpforwardRow2()
	// TODO: parse route to *nMIBIpforwardTable2
	err = createIPForwardEntry2(row)
	if err != nil {
		// TODO: parse err
		return err
	}
	return nil
}

func setRoute(route *Route) error {
	// TODO
	return fmt.Errorf("Not implemented yet")
	iph, err := syscall.LoadLibrary("Iphlpapi.dll")
	if err != nil {
		panic("LoadLibrary " + err.Error())
	}
	defer syscall.FreeLibrary(iph)
	nSetIPForwardEntry2 = getProcAddr(iph, "SetIpForwardEntry2")

	row := newMIBIpforwardRow2()
	// TODO: parse route to *nMIBIpforwardTable2
	err = setIPForwardEntry2(row)
	if err != nil {
		// TODO: parse err
		return err
	}
	return nil
}

func delRoute(route *Route) error {
	return fmt.Errorf("Not implemented yet")
	// TODO
	iph, err := syscall.LoadLibrary("Iphlpapi.dll")
	if err != nil {
		panic("LoadLibrary " + err.Error())
	}
	defer syscall.FreeLibrary(iph)
	nSetIPForwardEntry2 = getProcAddr(iph, "DeleteIpForwardEntry2")

	row := newMIBIpforwardRow2()
	// TODO: parse route to *nMIBIpforwardTable2
	err = deleteIPForwardEntry2(row)
	if err != nil {
		// TODO: parse err
		return err
	}
	return nil
}
