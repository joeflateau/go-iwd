package iwd

import (
	"github.com/godbus/dbus/v5"
)

const (
	objectIwd              = "net.connman.iwd"
	objectIwdPath          = "/net/connman/iwd"
	iwdAgentManager        = "net.connman.iwd.AgentManager"
	iwdAdapter             = "net.connman.iwd.Adapter"
	iwdDevice              = "net.connman.iwd.Device"
	iwdSimpleConfiguration = "net.connman.iwd.SimpleConfiguation"
	iwdNetwork             = "net.connman.iwd.Network"
)

// Iwd is a struct over all major iwd components
type Iwd struct {
	Agents        []Agent
	Adapters      []Adapter
	KnownNetworks []KnownNetwork
	Networks      []Network
	Stations      []Station
	Devices       []Device
}

// New parses the net.connman.iwd object index and initializes an iwd object
func New(conn *dbus.Conn) (*Iwd, error) {
	var objects map[dbus.ObjectPath]map[string]map[string]dbus.Variant
	objectManager := conn.Object(objectIwd, "/")
	err := objectManager.Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0).Store(&objects)

	if err != nil {
		return nil, err
	}

	iwd := Iwd{
		Agents:        make([]Agent, 0),
		Adapters:      make([]Adapter, 0),
		KnownNetworks: make([]KnownNetwork, 0),
		Networks:      make([]Network, 0),
		Stations:      make([]Station, 0),
		Devices:       make([]Device, 0),
	}

	for k, v := range objects {
		for resource, obj := range v {
			switch resource {
			case objectAdapter:
				iwd.Adapters = append(iwd.Adapters, Adapter{
					Path:           k,
					Model:          castOrDefault(obj["Model"], ""),
					Name:           castOrDefault(obj["Name"], ""),
					Powered:        obj["Powered"].Value().(bool),
					SupportedModes: obj["SupportedModes"].Value().([]string),
					Vendor:         castOrDefault(obj["Vendor"], ""),
				})
			case objectKnownNetwork:
				iwd.KnownNetworks = append(iwd.KnownNetworks, KnownNetwork{
					Path:              k,
					AutoConnect:       obj["AutoConnect"].Value().(bool),
					Hidden:            obj["Hidden"].Value().(bool),
					LastConnectedTime: obj["LastConnectedTime"].Value().(string),
					Name:              castOrDefault(obj["Name"], ""),
					Type:              obj["Type"].Value().(string),
				})
			case objectNetwork:
				iwd.Networks = append(iwd.Networks, Network{
					Path:      k,
					Connected: obj["Connected"].Value().(bool),
					Device:    obj["Device"].Value().(dbus.ObjectPath),
					Name:      castOrDefault(obj["Name"], ""),
					Type:      obj["Type"].Value().(string),
				})
			case objectStation:
				iwd.Stations = append(iwd.Stations, Station{
					Path:             k,
					ConnectedNetwork: castOrDefault[*dbus.ObjectPath](obj["ConnectedNetwork"], nil),
					Scanning:         obj["Scanning"].Value().(bool),
					State:            obj["State"].Value().(string),
				})
			case objectDevice:
				iwd.Devices = append(iwd.Devices, Device{
					Path:    k,
					Adapter: obj["Adapter"].Value().(dbus.ObjectPath),
					Address: obj["Address"].Value().(string),
					Mode:    obj["Mode"].Value().(string),
					Name:    castOrDefault(obj["Name"], ""),
					Powered: obj["Powered"].Value().(bool),
				})
			}
		}
	}
	return &iwd, nil
}
