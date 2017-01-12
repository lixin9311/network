package network

import (
	"fmt"
	"net"
)

// Route is the struct of an entry in route table
type Route struct {
	*net.IPNet
	Iface   *net.Interface
	Gateway net.IP
	Default bool
	Metric  int
}

func (r Route) String() string {
	var gw string
	if r.Gateway.IsUnspecified() {
		gw = "Local-link"
	} else {
		gw = fmt.Sprintf("%s", r.Gateway)
	}
	s := fmt.Sprintf("%s via %s id %d gateway %s metric %d", r.IPNet, r.Iface.Name, r.Iface.Index, gw, r.Metric)
	return s
}

// GetRoutes returns the route table of all networks on the system
func GetRoutes() ([]Route, error) {
	return getRoutes()
}
