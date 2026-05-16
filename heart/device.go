package heart

import (
	"errors"
	"fmt"
	"time"

	"github.com/avast/retry-go/v5"
	"tinygo.org/x/bluetooth"
)

type Device struct {
	Addr          bluetooth.Address
	Name          string
	Maker         []bluetooth.ManufacturerDataElement
	HaveHeartRate bool //
	RSSI          int16
}

func ScanDeviceWithTimeout(timeout time.Duration) ([]Device, error) {
	var devices = make(map[string]bluetooth.ScanResult)
	var failedChan = make(chan error, 1)
	go func() {
		err := adapter.Scan(func(a *bluetooth.Adapter, result bluetooth.ScanResult) {
			devices[result.Address.String()] = result
		})
		failedChan <- err
	}()
	select {
	case err := <-failedChan:
		return nil, err
	case <-time.After(timeout):
	}
	adapter.StopScan()
	rsp := make([]Device, 0, len(devices))
	for _, v := range devices {
		rsp = append(rsp, Device{
			Addr:          v.Address,
			Name:          v.LocalName(),
			Maker:         v.ManufacturerData(),
			RSSI:          v.RSSI,
			HaveHeartRate: v.HasServiceUUID(bluetooth.ServiceUUIDHeartRate),
		})
	}
	return rsp, nil
}

var device *bluetooth.Device
var deviceService *bluetooth.DeviceService

func ConnectDevice(addr bluetooth.Address) error {
	var s []bluetooth.DeviceService
	time.Sleep(1000 * time.Millisecond)
	var breakErr error = nil
	err := retry.New(
		retry.Attempts(3),
		retry.Delay(100*time.Millisecond),
	).Do(func() error {
		d, err := adapter.Connect(addr, bluetooth.ConnectionParams{})
		if err != nil {
			breakErr = fmt.Errorf("connect device %s error: %s", addr, err)
			return nil
		}
		device = &d
		srv, err := device.DiscoverServices([]bluetooth.UUID{bluetooth.ServiceUUIDHeartRate})
		if err != nil {
			device.Disconnect()
			return fmt.Errorf("discover service error: %s", err)
		}
		s = srv
		return nil
	})
	if breakErr != nil {
		return breakErr
	}
	if err != nil {
		return fmt.Errorf("discover service error: %s", err)
	}
	if len(s) == 0 {
		return errors.New("service not found")
	}
	deviceService = &(s[0])
	err = startRevcHeartRate()
	if err != nil {
		return err
	}
	return nil
}

func IsConnected() bool {
	if device == nil {
		return false
	}
	ok, err := device.Connected()
	if err != nil {
		return false
	}
	return ok
}

func DisconnectDevice() error {
	if device == nil {
		return nil
	}
	defer func() {
		device = nil
		deviceService = nil
	}()
	return device.Disconnect()
}
