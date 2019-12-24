package network

import (
	"net"
	"syscall"
	"unsafe"
)

//sys   getBestInterfaceEx(dstAddr *nRawSockAddrAny, ifIndex *nNetIFIndex) (ret error) = iphlpapi.GetBestInterfaceEx
//sys   getBestRoute2(ifLuid *nNetLuid, ifIndex nNetIFIndex, srcAddr *nRawSockAddrAny, dstAddr *nRawSockAddrAny, sortOpt uint32, bestRoute *nMIBIPForwardRow2, bestSrcAddr *nRawSockAddrAny) (ret error) = iphlpapi.GetBestRoute2
//sys   getIPForwardTable2(family nAddressFamily, table **nMIBIPForwardTable2) (ret error) = iphlpapi.GetIpForwardTable2
//sys   createIPForwardEntry2(row *nMIBIPForwardRow2) (ret error) = iphlpapi.CreateIpForwardEntry2
//sys   setIPForwardEntry2(row *nMIBIPForwardRow2) (ret error) = iphlpapi.SetIpForwardEntry2
//sys   deleteIPForwardEntry2(row *nMIBIPForwardRow2) (ret error) = iphlpapi.DeleteIpForwardEntry2
//sys   freeMibTable(table nMIBTable) = iphlpapi.FreeMibTable
//sys	convertInterfaceLUIDToGUID(interfaceLUID *nNetLuid, interfaceGUID *windows.GUID) (ret error) = iphlpapi.ConvertInterfaceLuidToGuid
//sys	convertInterfaceAliasToLUID(interfaceAlias *uint16, interfaceLUID *nNetLuid) (ret error) = iphlpapi.ConvertInterfaceAliasToLuid
//sys	convertInterfaceLUIDToAlias(interfaceLUID *nNetLuid, interfaceAlias *uint16, size uintptr) (ret error) = iphlpapi.ConvertInterfaceLuidToAlias
//sys   convertInterfaceGUIDToLUID(interfaceGUID *windows.GUID, interfaceLUID *nNetLuid) (ret error) = iphlpapi.ConvertInterfaceGuidToLuid
//sys   convertInterfaceLUIDToIndex(interfaceLUID *nNetLuid, interfaceIndex *nNetIFIndex) (ret error) = iphlpapi.ConvertInterfaceLuidToIndex
//sys   convertInterfaceIndexToLUID(interfaceIndex nNetIFIndex, interfaceLUID *nNetLuid) (ret error) = iphlpapi.ConvertInterfaceIndexToLuid
//sys   getIPInterfaceTable(family nAddressFamily, table **nMIBIPInterfaceTable) (ret error) = iphlpapi.GetIpInterfaceTable
//sys   setIPInterfaceEntry(row *nMIBIPInterfaceRow) (ret error) = iphlpapi.SetIpInterfaceEntry
//sys   getUnicastIPAddressTable(family nAddressFamily, table **nMIBUnicastIPAddressTable) (ret error) = iphlpapi.GetUnicastIpAddressTable
//sys   setUnicastIPAddressEntry(row *nMIBUnicastIPAddressRow) (ret error) = iphlpapi.SetUnicastIpAddressEntry
//sys   createUnicastIPAddressEntry(row *nMIBUnicastIPAddressRow) (ret error) = iphlpapi.CreateUnicastIpAddressEntry
//sys   getUnicastIPAddressEntry(row *nMIBUnicastIPAddressRow) (ret error) = iphlpapi.GetIpInterfaceEntry
//sys   initializeUnicastIPAddressEntry(row *nMIBUnicastIPAddressRow) = iphlpapi.InitializeUnicastIpAddressEntry
//sys   deleteUnicastIPAddressEntry(row *nMIBUnicastIPAddressRow) = iphlpapi.DeleteUnicastIpAddressEntry

type nMIBTable interface {
	unsafePointer() uintptr
}

type nNetLuid uint64
type nNetIFIndex uint32

type nAddressFamily uint16

const (
	nAFUnspecified nAddressFamily = syscall.AF_UNSPEC
	nAFIPv4        nAddressFamily = syscall.AF_INET
	nAFIPv6        nAddressFamily = syscall.AF_INET6
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

func (s *nRawSockAddrAny) SetIP(ip net.IP) {
	for i := range s.data {
		s.data[i] = 0
	}
	if ip4 := ip.To4(); ip4 != nil {
		s.family = nAFIPv4
		copy(s.data[2:], ip4)
	} else if ip6 := ip.To16(); ip6 != nil {
		s.family = nAFIPv6
		copy(s.data[6:], ip6)
	}
}

func (s *nRawSockAddrAny) ToIP() (ip net.IP) {
	if s.family == nAFIPv4 {
		ip = net.IP(make([]byte, net.IPv4len))
		copy(ip, s.data[2:])
	} else if s.family == nAFIPv6 {
		ip = net.IP(make([]byte, net.IPv6len))
		copy(ip, s.data[6:])
	}
	return
}

func NewSockAddrFromIP(ip net.IP) *nRawSockAddrAny {
	ret := &nRawSockAddrAny{}
	if ip4 := ip.To4(); ip4 != nil {
		ret.family = nAFIPv4
		copy(ret.data[2:], ip4)
		return ret
	}
	if ip6 := ip.To16(); ip6 != nil {
		ret.family = nAFIPv6
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

type nMIBIPForwardRow2 struct {
	InterfaceLuid        nNetLuid
	InterfaceIndex       nNetIFIndex
	DestinationPrefix    nIPAddressPrefix
	NextHop              nRawSockAddrAny
	SitePrefixLength     byte
	_                    [3]byte // padding
	ValidLifetime        uint32
	PreferredLifetime    uint32
	Metric               uint32
	Protocol             nNLRouteProtocol
	_                    uint16
	Loopback             bool
	AutoconfigureAddress bool
	Publish              bool
	Immortal             bool
	Age                  uint32
	Origin               nNLRouteOrigin
	_                    uint16
}

func (row *nMIBIPForwardRow2) ToRoute() *Route {
	family := row.DestinationPrefix.Prefix.family
	saData := row.DestinationPrefix.Prefix.data
	gwData := row.NextHop.data
	length := int(row.DestinationPrefix.PrefixLength)
	metric := int(row.Metric)
	var ip []byte
	var gateway []byte
	var mask net.IPMask
	if family == nAFIPv4 {
		ip = make([]byte, net.IPv4len)
		copy(ip, saData[2:])
		mask = net.CIDRMask(length, net.IPv4len*8)
		gateway = make([]byte, net.IPv4len)
		copy(gateway, gwData[2:])
	} else if family == nAFIPv6 {
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

func newMIBIPForwardRow2() *nMIBIPForwardRow2 {
	return &nMIBIPForwardRow2{
		ValidLifetime:     0xffffffff,
		PreferredLifetime: 0xffffffff,
		Protocol:          nRouteProtocolNetMgmt,
		Immortal:          true,
		Origin:            nNlroManual,
		Metric:            0xffffffff,
	}
}

type nMIBIPForwardTable2 struct {
	NumEntries uint32
	Table      [128]nMIBIPForwardRow2
}

func (table *nMIBIPForwardTable2) unsafePointer() uintptr {
	return uintptr(unsafe.Pointer(table))
}

type nNLPrefixOrigin uint16

const (
	nIPPrefixOriginOther nNLPrefixOrigin = iota
	nIPPrefixOriginManual
	nIPPrefixOriginWellKnown
	nIPPrefixOriginDhcp
	nIPPrefixOriginRouterAdvertisement
	nIPPrefixOriginUnchanged
)

type nNLSuffixOrigin uint16

const (
	nNlsoOther nNLSuffixOrigin = iota
	nNlsoManual
	nNlsoWellKnown
	nNlsoDhcp
	nNlsoLinkLayerAddress
	nNlsoRandom
	nIPSuffixOriginOther
	nIPSuffixOriginManual
	nIPSuffixOriginWellKnown
	nIPSuffixOriginDhcp
	nIPSuffixOriginLinkLayerAddress
	nIPSuffixOriginRandom
	nIPSuffixOriginUnchanged
)

type nNLDadState uint16

const (
	nNldsInvalid nNLDadState = iota
	nNldsTentative
	nNldsDuplicate
	nNldsDeprecated
	nNldsPreferred
	nIPDadStateInvalid
	nIPDadStateTentative
	nIPDadStateDuplicate
	nIPDadStateDeprecated
	nIPDadStatePreferred
)

type nMIBUnicastIPAddressRow struct {
	Address            nRawSockAddrAny
	InterfaceLuid      nNetLuid
	InterfaceIndex     nNetIFIndex
	PrefixOrigin       nNLPrefixOrigin
	_                  uint16
	SuffixOrigin       nNLSuffixOrigin
	_                  uint16
	ValidLifetime      uint32
	PreferredLifetime  uint32
	OnLinkPrefixLength uint8
	SkipAsSource       bool
	_                  uint16
	DadState           nNLDadState
	ScopeID            uint32
	CreationTimeStamp  uint64
}

type nMIBUnicastIPAddressTable struct {
	NumEntries uint32
	Table      [128]nMIBUnicastIPAddressRow
}

func (table *nMIBUnicastIPAddressTable) unsafePointer() uintptr {
	return uintptr(unsafe.Pointer(table))
}

// SCOPE_ID
// typedef struct {
// 	union {
// 	  struct {
// 		ULONG  Zone : 28;
// 		ULONG  Level : 4;
// 	  };
// 	  ULONG  Value;
// 	};
//   } SCOPE_ID, *PSCOPE_ID;

type nNLRouterDiscoveryBehavior uint16

const (
	nRouterDiscoveryDisabled nNLRouterDiscoveryBehavior = iota
	nRouterDiscoveryEnabled
	nRouterDiscoveryDhcp
	nRouterDiscoveryUnchanged
)

type nNLLinkLocalAddressBehavior uint16

const (
	nLinkLocalAlwaysOff nNLLinkLocalAddressBehavior = iota
	nLinkLocalDelayed
	nLinkLocalAlwaysOn
	nLinkLocalUnchanged
)

type nNLInterfaceOffloadRod uint8

const (
	nNlChecksumSupported = 1 << iota
	nNlOptionsSupported
	nTlDatagramChecksumSupported
	nTlStreamChecksumSupported
	nTlStreamOptionsSupported
	nFastPathCompatible
	nTlLargeSendOffloadSupported
	nTlGiantSendOffloadSupported
)

var nNLInterfaceOffloadRodNames = []string{
	"ChecksumSupported",
	"OptionsSupported",
	"DatagramChecksumSupported",
	"StreamChecksumSupported",
	"StreamOptionsSupported",
	"FastPathCompatible",
	"LargeSendOffloadSupported",
	"GiantSendOffloadSupported",
}

func (rod nNLInterfaceOffloadRod) String() string {
	s := ""
	for i, name := range nNLInterfaceOffloadRodNames {
		if rod&(1<<uint(i)) != 0 {
			if s != "" {
				s += "|"
			}
			s += name
		}
	}
	if s == "" {
		s = "0"
	}
	return s
}

type nMIBIPInterfaceRow struct {
	Family                               nAddressFamily
	_                                    uint16
	InterfaceLuid                        nNetLuid
	InterfaceIndex                       nNetIFIndex
	MaxReassemblySize                    uint32
	InterfaceIdentifier                  uint64
	MinRouterAdvertisementInterval       uint32
	MaxRouterAdvertisementInterval       uint32
	AdvertisingEnabled                   bool
	ForwardingEnabled                    bool
	WeakHostSend                         bool
	WeakHostReceive                      bool
	UseAutomaticMetric                   bool
	UseNeighborUnreachabilityDetection   bool
	ManagedAddressConfigurationSupported bool
	OtherStatefulConfigurationSupported  bool
	AdvertiseDefaultRoute                bool
	_                                    [3]byte // padding
	RouterDiscoveryBehavior              nNLRouterDiscoveryBehavior
	_                                    uint16
	DadTransmits                         uint32
	BaseReachableTime                    uint32
	RetransmitTime                       uint32
	PathMtuDiscoveryTimeout              uint32
	LinkLocalAddressBehavior             nNLLinkLocalAddressBehavior
	_                                    uint16
	LinkLocalAddressTimeout              uint32
	ZoneIndices                          [16]uint32
	SitePrefixLength                     uint32
	Metric                               uint32
	NlMtu                                uint32
	Connected                            bool
	SupportsWakeUpPatterns               bool
	SupportsNeighborDiscovery            bool
	SupportsRouterDiscovery              bool
	ReachableTime                        uint32
	TransmitOffload                      nNLInterfaceOffloadRod
	ReceiveOffload                       nNLInterfaceOffloadRod
	DisableDefaultRoutes                 bool
	_                                    uint8
}

type nMIBIPInterfaceTable struct {
	NumEntries uint32
	Table      [128]nMIBIPInterfaceRow
}

func (table *nMIBIPInterfaceTable) unsafePointer() uintptr {
	return uintptr(unsafe.Pointer(table))
}
