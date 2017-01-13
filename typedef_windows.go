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
	addr   [4]byte
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

func ParseIP(ip net.IP) *nRawSockAddrAny {
	ret := &nRawSockAddrAny{}
	if ip4 := ip.To4(); ip4 != nil {
		ret.family = nAF_INET
		copy(ret.data[2:], ip4)
		return ret
	}
	if ip6 := ip.To16(); ip6 != nil {
		ret.family = nAF_INET6
		copy(ret.data[6:], ip6)
		return ret
	}
	return nil
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

func (row *nMIBIpforwardRow2) ToRoute() *Route {
	family := row.DestinationPrefix.Prefix.family
	saData := row.DestinationPrefix.Prefix.data
	gwData := row.NextHop.data
	length := int(row.DestinationPrefix.PrefixLength)
	metric := int(row.Metric)
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
		return nil
	}
	ipnet := &net.IPNet{IP: ip, Mask: mask}
	ifce, _ := net.InterfaceByIndex(int(row.InterfaceIndex))
	isDefault := false
	if length == 0 {
		isDefault = true
	}
	return &Route{ipnet, ifce, gateway, isDefault, metric}
}

func newMIBIpforwardRow2() *nMIBIpforwardRow2 {
	return &nMIBIpforwardRow2{
		ValidLifetime:     0xffffffff,
		PreferredLifetime: 0xffffffff,
		Protocol:          nRouteProtocolNetMgmt,
		Immortal:          true,
		Origin:            nNlroManual,
		Metric:            0xffffffff,
	}
}

type nMIBIpforwardTable2 struct {
	NumEntries uint32
	Table      [128]nMIBIpforwardRow2
}
