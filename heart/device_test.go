package heart

import (
	"testing"
	"time"
)

func TestStartScanDevice(t *testing.T) {
	Start()
	scs, err := ScanDeviceWithTimeout(20 * time.Second)
	t.Log(scs, err)
	Close()
}
