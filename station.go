package iwd

import "github.com/godbus/dbus/v5"

const (
	objectStation   = "net.connman.iwd.Station"
	callStationScan = "net.connman.iwd.Station.Scan"
)

// Station refers to net.connman.iwd.Station
type Station struct {
	Path             dbus.ObjectPath
	ConnectedNetwork *dbus.ObjectPath
	Scanning         bool
	State            string
}

// Scan scans for wireless networks
func (s *Station) Scan(conn *dbus.Conn) error {
	obj := conn.Object(objectStation, "")
	call := obj.Call(callStationScan, 0)
	if call.Err != nil {
		return call.Err
	}
	return nil
}
