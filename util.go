package iwd

import "github.com/godbus/dbus/v5"

func castOrDefault[Type any](v dbus.Variant, fallback Type) Type {
	cast, ok := v.Value().(Type)
	if !ok {
		cast = fallback
	}
	return cast
}
