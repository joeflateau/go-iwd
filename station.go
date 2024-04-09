package iwd

import (
	"github.com/godbus/dbus/v5"
)

const (
	objectStation                 = "net.connman.iwd.Station"
	callStationScan               = "net.connman.iwd.Station.Scan"
	callStationGetOrderedNetworks = "net.connman.iwd.Station.GetOrderedNetworks"
	callStationDisconnect         = "net.connman.iwd.Station.Disconnect"
)

// Station refers to net.connman.iwd.Station
type Station struct {
	iwd              *Iwd
	Path             dbus.ObjectPath
	ConnectedNetwork *dbus.ObjectPath
	Scanning         bool
	State            string
}

type OrderedNetwork struct {
	Network
	SignalStrength int16
}

// Scan scans for wireless networks
func (s *Station) Scan(conn *dbus.Conn) error {
	obj := conn.Object(objectIwd, s.Path)
	call := obj.Call(callStationScan, 0)
	if call.Err != nil {
		return call.Err
	}
	return nil
}

func (s *Station) GetOrderedNetworks(conn *dbus.Conn) ([]*OrderedNetwork, error) {
	obj := conn.Object(objectIwd, s.Path)

	var objects [][]dbus.Variant
	err := obj.Call(callStationGetOrderedNetworks, 0).Store(&objects)
	if err != nil {
		return nil, err
	}
	orderedNetworks := make([]*OrderedNetwork, 0, len(objects))
	for _, obj := range objects {
		networkObject := obj[0].Value().(dbus.ObjectPath)
		signalStrength := obj[1].Value().(int16)

		for _, network := range s.iwd.Networks {
			if network.Path == networkObject {
				orderedNetworks = append(orderedNetworks, &OrderedNetwork{
					Network:        network,
					SignalStrength: signalStrength,
				})
				break
			}
		}
	}
	return orderedNetworks, nil
}

func (s *Station) Disconnect(conn *dbus.Conn) error {
	obj := conn.Object(objectIwd, s.Path)
	call := obj.Call(callStationDisconnect, 0)
	if call.Err != nil {
		return call.Err
	}
	return nil
}
