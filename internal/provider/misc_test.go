package provider

import (
	"testing"
)

func TestSimpleResourceLifecycles(t *testing.T) {
	cases := []lifecycle{
		{
			resourceType: "nfs_share", address: "truenas_nfs_share.media",
			create: `resource "truenas_nfs_share" "media" {
  path = "/mnt/tank/data/media"
}`,
			update: `resource "truenas_nfs_share" "media" {
  path     = "/mnt/tank/data/media"
  ro       = true
  networks = ["192.0.2.0/24"]
}`,
		},
		{
			resourceType: "static_route", address: "truenas_static_route.vpn",
			create: `resource "truenas_static_route" "vpn" {
  destination = "198.51.100.0/24"
  gateway     = "192.0.2.1"
}`,
			update: `resource "truenas_static_route" "vpn" {
  destination = "198.51.100.0/24"
  gateway     = "192.0.2.1"
  description = "VPN macvlan"
}`,
		},
		{
			resourceType: "tunable", address: "truenas_tunable.arc",
			create: `resource "truenas_tunable" "arc" {
  type  = "ZFS"
  var   = "zfs_arc_max"
  value = "68719476736"
}`,
			update: `resource "truenas_tunable" "arc" {
  type    = "ZFS"
  var     = "zfs_arc_max"
  value   = "68719476736"
  comment = "64 GiB ARC ceiling"
}`,
		},
		{
			resourceType: "iscsi_target", address: "truenas_iscsi_target.vm1",
			create: `resource "truenas_iscsi_target" "vm1" {
  name = "vm1-data"
}`,
			update: `resource "truenas_iscsi_target" "vm1" {
  name  = "vm1-data"
  alias = "Windows test data"
}`,
		},
		{
			resourceType: "iscsi_extent", address: "truenas_iscsi_extent.vm1",
			onWrite: func(row map[string]any) {
				if row["serial"] == nil {
					row["serial"] = "a1b2c3d4"
				}
			},
			create: `resource "truenas_iscsi_extent" "vm1" {
  name = "vm1-data"
  type = "DISK"
  disk = "zvol/tank/vm1-data"
}`,
			update: `resource "truenas_iscsi_extent" "vm1" {
  name    = "vm1-data"
  type    = "DISK"
  disk    = "zvol/tank/vm1-data"
  comment = "G: drive"
}`,
		},
		{
			resourceType: "iscsi_auth", address: "truenas_iscsi_auth.vm1",
			create: `resource "truenas_iscsi_auth" "vm1" {
  tag    = 1
  user   = "vm1"
  secret = "at-least-8-chars"
}`,
			update: `resource "truenas_iscsi_auth" "vm1" {
  tag                = 1
  user               = "vm1"
  secret             = "rotated-secret"
  secret_wo_version  = 1
}`,
		},
	}
	for _, lc := range cases {
		t.Run(lc.resourceType, func(t *testing.T) { lc.run(t) })
	}
}
