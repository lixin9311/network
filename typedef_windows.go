// +build windows

package network

import (
	"net"
	"syscall"
)

type nNetLuid uint64
type nNetIFIndex uint32

type nAddressFamily uint16

const (
	nAF_UNSPEC nAddressFamily = syscall.AF_UNSPEC
	nAF_INET   nAddressFamily = syscall.AF_INET
	nAF_INET6  nAddressFamily = syscall.AF_INET6
)

type nNLRouteProtocol uint16

const (
	nRouteProtocolOther nNLRouteProtocol = iota + 1
	nRouteProtocolLocal
	nRouteProtocolNetMgmt
	nRouteProtocolIcmp
	nRouteProtocolEgp
	nRouteProtocolGgp
	nRouteProtocolHello
	nRouteProtocolRip
	nRouteProtocolIsIs
	nRouteProtocolEsIs
	nRouteProtocolCisco
	nRouteProtocolBbn
	nRouteProtocolOspf
	nRouteProtocolBgp
	nRouteProtocolNTAutoStatic   nNLRouteProtocol = 10002
	nRouteProtocolNTStatic       nNLRouteProtocol = 10006
	nRouteProtocolNTStaticNonDod nNLRouteProtocol = 10007
)

type nNLRouteOrigin uint16

const (
	nNlroManual nNLRouteOrigin = iota
	nNlroWellKnown
	nNlroDHCP
	nNlroRouterAdvertisement
	nNlro6to4
)

type nRawSockaddrInet4 struct {
	family nAddressFamily
	port   uint16
	addr   net.IP
	zero   [8]byte
}

type nRawSockaddrInet6 struct {
	family   nAddressFamily
	port     uint16
	flowinfo uint32
	addr     [16]byte
	scopeID  uint32
}

type nRawSockAddrAny struct {
	family nAddressFamily
	data   [26]byte
}

type nIPAddressPrefix struct {
	Prefix       nRawSockAddrAny
	PrefixLength uint8
	_            [2]uint8 //padding
}

type nMIBIpforwardRow2 struct {
	InterfaceLuid        nNetLuid
	InterfaceIndex       nNetIFIndex
	DestinationPrefix    nIPAddressPrefix
	NextHop              nRawSockAddrAny
	SitePrefixLength     byte
	_                    [2]byte // padding
	ValidLifetime        uint32
	PreferredLifetime    uint32
	Metric               uint32
	Protocol             nNLRouteProtocol
	Loopback             bool
	AutoconfigureAddress bool
	Publish              bool
	Immortal             bool
	Age                  uint32
	Origin               nNLRouteOrigin
}

type nMIBIpforwardTable2 struct {
	NumEntries uint32
	Table      [128]nMIBIpforwardRow2
}
