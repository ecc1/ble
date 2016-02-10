Package ble provides functions to discover, connect, pair,
and communicate with Bluetooth Low Energy peripheral devices.

Documentation: <https://godoc.org/github.com/ecc1/ble>

This implementation uses the BlueZ D-Bus interface, rather than sockets.
It is similar to <https://github.com/adafruit/Adafruit_Python_BluefruitLE>

The apps directory contains some simple example programs.

Some older Linux kernels, like the one on the Intel Edison, may not
properly support the SetDiscoveryFilter method.  The ble package can
be built with the "nofilter" tag to work around this.
