package network

import (
	"fmt"
	"net"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	errInvalidAddr = fmt.Errorf("Provided address is invalid")
)

func InterfaceAliasToGUID(alias string) (guid *windows.GUID, err error) {
	var luid nNetLuid
	guid = &windows.GUID{}
	err = convertInterfaceAliasToLUID(windows.StringToUTF16Ptr(alias), &luid)
	if err != nil {
		return
	}
	err = convertInterfaceLUIDToGUID(&luid, guid)
	return
}

func InterfaceGUIDToAlias(guid *windows.GUID) (alias string, err error) {
	var luid nNetLuid
	strBuf := make([]uint16, 256+1)
	err = convertInterfaceGUIDToLUID(guid, &luid)
	if err != nil {
		return
	}
	err = convertInterfaceLUIDToAlias(&luid, &strBuf[0], 256+1)
	if err != nil {
		return
	}
	alias = windows.UTF16ToString(strBuf)
	return
}

func InterfaceGUIDToIndex(guid *windows.GUID) (index int, err error) {
	var luid nNetLuid
	var ifaceIndex nNetIFIndex
	err = convertInterfaceGUIDToLUID(guid, &luid)
	if err != nil {
		return
	}
	err = convertInterfaceLUIDToIndex(&luid, &ifaceIndex)
	index = int(ifaceIndex)
	return
}

func InterfaceIndexToGUID(index int) (guid *windows.GUID, err error) {
	var luid nNetLuid
	guid = &windows.GUID{}
	ifaceIndex := nNetIFIndex(index)
	err = convertInterfaceIndexToLUID(ifaceIndex, &luid)
	if err != nil {
		return
	}
	err = convertInterfaceLUIDToGUID(&luid, guid)
	return
}

func getRoute(ip net.IP) (*Route, error) {
	var ifIndex nNetIFIndex
	dstAddr := NewSockAddrFromIP(ip)
	if dstAddr == nil {
		return nil, errInvalidAddr
	}
	if err := getBestInterfaceEx(dstAddr, &ifIndex); err != nil {
		return nil, err
	}
	row := newMIBIPForwardRow2()
	src := nRawSockAddrAny{}
	if err := getBestRoute2(nil, ifIndex, nil, dstAddr, 0, row, &src); err != nil {
		return nil, err
	}
	return row.ToRoute(), nil
}

func getRoutes() ([]Route, error) {
	table := &nMIBIPForwardTable2{}
	if err := getIPForwardTable2(nAFUnspecified, &table); err != nil {
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
	row := newMIBIPForwardRow2()
	// TODO: parse route to *nMIBIPForwardTable2
	err := createIPForwardEntry2(row)
	if err != nil {
		// TODO: parse err
		return err
	}
	return nil
}

func setRoute(route *Route) error {
	// TODO
	return fmt.Errorf("Not implemented yet")
	row := newMIBIPForwardRow2()
	// TODO: parse route to *nMIBIPForwardTable2
	err := setIPForwardEntry2(row)
	if err != nil {
		// TODO: parse err
		return err
	}
	return nil
}

func delRoute(route *Route) error {
	return fmt.Errorf("Not implemented yet")
	// TODO
	row := newMIBIPForwardRow2()
	// TODO: parse route to *nMIBIPForwardTable2
	err := deleteIPForwardEntry2(row)
	if err != nil {
		// TODO: parse err
		return err
	}
	return nil
}
