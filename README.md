# network - network tools library for go #

The network provides a simple set of network tools for go.
Currently we are only supporting GetRoutes on Windows,
but we are planning to do:

- [x] GetRoute(WIN32)
- [ ] Set/Change/DelRoute(WIN32)
- [ ] Get/SetAddr(WIN32)
- [ ] LinkSetUP/DOWN(WIN32)
- [ ] The same interface using Netlink on Linux

It's only a side work for me, so don't expect fast development or any time schedule and guarantee.

Currently my priority is the implement on WIN32, because there are existing
tools can help you to do the same job on Linux with Netlink.
Like:

1. [docker/libcontainer/netlink](https://github.com/docker/libcontainer/tree/master/netlink)
2. [vishvananda/netlink](https://github.com/vishvananda/netlink)

If you are interested in this project, any contribution will be appreciated.



## Examples ##

See [demo directory](https://github.com/lixin9311/network/tree/master/demo).

Get and print all entries in route table:

```go
package main

import (
	"fmt"

	"github.com/lixin9311/network"
)

func main() {
	r, _ := network.GetRoutes()
	for _, v := range r {
		fmt.Println("route:", v)
	}
}
```

With `route` tool from the system:

```dos
C:\mygo\src\github.com\lixin9311\network\demo>route print
===========================================================================
Interface List
 11...30 5a 3a 76 09 60 ......Intel(R) I210 Gigabit Network Connection
 14...00 ff 77 e6 aa 05 ......TAP-Windows Adapter V9
  1...........................Software Loopback Interface 1
===========================================================================

IPv4 Route Table
===========================================================================
Active Routes:
Network Destination        Netmask          Gateway       Interface  Metric
          0.0.0.0          0.0.0.0      192.168.1.1      192.168.1.6     30
        127.0.0.0        255.0.0.0         On-link         127.0.0.1    331
        127.0.0.1  255.255.255.255         On-link         127.0.0.1    331
  127.255.255.255  255.255.255.255         On-link         127.0.0.1    331
      192.168.1.0    255.255.255.0         On-link       192.168.1.6    271
      192.168.1.6  255.255.255.255         On-link       192.168.1.6    271
    192.168.1.255  255.255.255.255         On-link       192.168.1.6    271
        224.0.0.0        240.0.0.0         On-link         127.0.0.1    331
        224.0.0.0        240.0.0.0         On-link       192.168.1.6    271
  255.255.255.255  255.255.255.255         On-link         127.0.0.1    331
  255.255.255.255  255.255.255.255         On-link       192.168.1.6    271
===========================================================================
Persistent Routes:
  None

IPv6 Route Table
===========================================================================
Active Routes:
 If Metric Network Destination      Gateway
 11    281 ::/0                     fe80::1
  1    331 ::1/128                  On-link
 11    281 240d:1a:60e:7f00::/64    On-link
 11    281 240d:1a:60e:7f00:5d3e:c004:7bb6:5bc4/128
                                    On-link
 11    281 240d:1a:60e:7f00:ccfc:9179:5:418c/128
                                    On-link
 11    281 fe80::/64                On-link
 11    281 fe80::ccfc:9179:5:418c/128
                                    On-link
  1    331 ff00::/8                 On-link
 11    281 ff00::/8                 On-link
===========================================================================
Persistent Routes:
  None
```

You may get output like this.

```dos
C:\mygo\src\github.com\lixin9311\network\demo>go run GetRoute.go
route: 0.0.0.0/0 via Ethernet id 11 gateway 192.168.1.1 metric 15
route: 127.0.0.0/8 via Loopback Pseudo-Interface 1 id 1 gateway 0.0.0.0 metric 256
route: 127.0.0.1/32 via Loopback Pseudo-Interface 1 id 1 gateway 0.0.0.0 metric 256
route: 127.255.255.255/32 via Loopback Pseudo-Interface 1 id 1 gateway 0.0.0.0 metric 256
route: 192.168.1.0/24 via Ethernet id 11 gateway 0.0.0.0 metric 256
route: 192.168.1.6/32 via Ethernet id 11 gateway 0.0.0.0 metric 256
route: 192.168.1.255/32 via Ethernet id 11 gateway 0.0.0.0 metric 256
route: 224.0.0.0/4 via Loopback Pseudo-Interface 1 id 1 gateway 0.0.0.0 metric 256
route: 224.0.0.0/4 via Ethernet 3 id 14 gateway 0.0.0.0 metric 256
route: 224.0.0.0/4 via Ethernet id 11 gateway 0.0.0.0 metric 256
route: 255.255.255.255/32 via Loopback Pseudo-Interface 1 id 1 gateway 0.0.0.0 metric 256
route: 255.255.255.255/32 via Ethernet 3 id 14 gateway 0.0.0.0 metric 256
route: 255.255.255.255/32 via Ethernet id 11 gateway 0.0.0.0 metric 256
route: ::/0 via Ethernet id 11 gateway fe80::1 metric 256
route: ::1/128 via Loopback Pseudo-Interface 1 id 1 gateway Local-link metric 256
route: 240d:1a:60e:7f00::/64 via Ethernet id 11 gateway Local-link metric 256
route: 240d:1a:60e:7f00:5d3e:c004:7bb6:5bc4/128 via Ethernet id 11 gateway Local-link metric 256
route: 240d:1a:60e:7f00:ccfc:9179:5:418c/128 via Ethernet id 11 gateway Local-link metric 256
route: fe80::/64 via Ethernet id 11 gateway Local-link metric 256
route: fe80::ccfc:9179:5:418c/128 via Ethernet id 11 gateway Local-link metric 256
route: ff00::/8 via Loopback Pseudo-Interface 1 id 1 gateway Local-link metric 256
route: ff00::/8 via Ethernet id 11 gateway Local-link metric 256
```