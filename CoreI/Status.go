package CoreI

import (
	core "github.com/v2fly/v2ray-core/v4"
)

type Status struct {
	IsRunning       bool
	VpnSupportnodup bool
	PackageName     string

	Vpoint core.Server
}

func (v *Status) GetDataDir() string {
	return v.getDataDir()
}

func (v *Status) getDataDir() string {
	var datadir = "/data/data/org.kkdev.v2raygo/"
	if v.PackageName != "" {
		datadir = "/data/data/" + v.PackageName + "/"
	}
	return datadir
}
