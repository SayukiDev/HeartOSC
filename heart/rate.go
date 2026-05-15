package heart

import (
	"errors"
	"sync/atomic"

	"tinygo.org/x/bluetooth"
)

var rate = atomic.Int32{}

func startRevcHeartRate() error {
	if device == nil {
		return errors.New("device is not connected")
	}
	ok, err := device.Connected()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("device is not connected")
	}
	chars, err := deviceService.DiscoverCharacteristics([]bluetooth.UUID{bluetooth.CharacteristicUUIDHeartRateMeasurement})
	if err != nil {
		return err
	}
	err = chars[0].EnableNotifications(func(buf []byte) {
		if buf[0]&0x01 == 0 {
			rate.Store(int32(buf[1]))
		} else {
			rate.Store(int32(buf[1]) | (int32(buf[2]) << 8))
		}
	})
	if err != nil {
		return err
	}
	return nil
}

func GetHeartRate() int32 {
	return rate.Load()
}
