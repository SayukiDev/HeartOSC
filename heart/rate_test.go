package heart

import (
	"testing"
	"time"
)

func TestGetHeartRate(t *testing.T) {
	err := Start()
	if err != nil {
		t.Error(err)
		return
	}
	defer Close()
	_, err = ScanAndConnectDeviceWithTimeout("FA:BE:C5:43:DC:1F", 20*time.Second)
	if err != nil {
		t.Error(err)
		return
	}
	err = startRevcHeartRate()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(IsConnected())
}
