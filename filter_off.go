// +build nofilter

package ble

// Discvoery filtering doesn't work on Intel Edison
// running Debian stretch and kernel 3.10.17-yocto-standard-r2.
func (adapter *blob) SetDiscoveryFilter(uuids ...string) error {
	return nil
}
