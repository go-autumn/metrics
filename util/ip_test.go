package util

import "testing"

func TestIntranetIP2(t *testing.T) {
	IntranetIP()
}

func TestIsIntranet(t *testing.T) {
	IsIntranet("172.134.0.1")
	IsIntranet("172.134.1")
	IsIntranet("172.你.1.0")
	IsIntranet("172.17.1.0")
}
