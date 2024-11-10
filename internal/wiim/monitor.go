package wiim

import "fmt"

type DeviceMonitor struct {
	muxer *MuxClient
}

func NewDeviceMonitor(muxer *MuxClient) *DeviceMonitor {
	return &DeviceMonitor{muxer: muxer}
}

func (dm *DeviceMonitor) Poll(name string) error {
	d := dm.muxer.GetDevice(name)
	s, err := d.GetPlayerStatus()
	if err != nil {
		return err
	}
	fmt.Printf("Player Status: %v\n", s)

	return nil
}

func (dm *DeviceMonitor) PollAll() error {
	for name, d := range *dm.muxer {
		s, err := d.GetUpdatedPlayerStatus()
		if err != nil {
			return err
		}
		if s == nil {
			return nil
		}
		fmt.Printf("Player Status updated: %s %v\n", name, s)
	}
	return nil
}
