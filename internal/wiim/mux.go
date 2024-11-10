package wiim

import (
	"encoding/json"
	"fmt"
	"iter"
	"maps"
	"strconv"
	"time"
)

type MuxClient map[string]*Device

func NewMuxClient() *MuxClient {
	mc := make(MuxClient)
	return &mc
}

func (mc *MuxClient) AddDevice(name string, device *Device) {
	(*mc)[name] = device
}

func (mc *MuxClient) GetDevice(name string) *Device {
	return (*mc)[name]
}

func (mc *MuxClient) GetDevices() iter.Seq[*Device] {
	// convert map to array
	// return slices.Collect(maps.Values(mc))
	return maps.Values(*mc)
}

func (mc *MuxClient) GetDevicesForUpdate() iter.Seq[*Device] {
	return func(yield func(*Device) bool) {
		now := time.Now().Unix()
		for _, v := range *mc {
			if v.NextUpdateTime > now {
				continue
			}
			if !yield(v) {
				return
			}
		}
	}

}

// Command Handler
func (mc *MuxClient) HandleCommand(name string, cmd string, arg1 *string, arg2 *string, payload *[]byte) error {
	c := (*mc)[name]

	switch cmd {
	case "pause":
		return c.Pause()
	case "resume":
		return c.Resume()
	case "togglePausePlay":
		return c.TogglePausePlay()
	case "previous":
		return c.Previous()
	case "next":
		return c.Next()
	case "seek":
		seconds, err := strconv.Atoi(*arg1)
		if err != nil {
			return err
		}
		return c.Seek(seconds)
	case "stop":
		return c.Stop()
	case "setLoopMode":
		mode, err := strconv.Atoi(*arg1)
		if err != nil {
			return err
		}
		return c.SetLoopMode(PlayerLoopMode(mode))
	case "setVolume":
		level, err := strconv.Atoi(*arg1)
		if err != nil {
			return err
		}
		return c.SetVolume(level)
	case "mute":
		muted, err := strconv.ParseBool(*arg1)
		if err != nil {
			return err
		}
		return c.Mute(muted)
	case "eqOn":
		_, err := c.EQOn()
		return err
	case "eqOff":
		_, err := c.EQOff()
		return err
	case "eqGetStat":
		resp, err := c.EQGetStat()
		if err != nil {
			return err
		}
		bytes, err := json.Marshal(resp)
		if err != nil {
			return err
		}
		*payload = bytes
		return nil
	case "eqGetList":
		resp, err := c.EQGetList()
		if err != nil {
			return err
		}
		bytes, err := json.Marshal(resp)
		if err != nil {
			return err
		}
		*payload = bytes
		return nil
	case "eqLoad":
		_, err := c.EQLoad(*arg1)
		return err
	case "reboot":
		_, err := c.Reboot()
		return err
	case "shutdown":
		sec, err := strconv.Atoi(*arg1)
		if err != nil {
			return err
		}
		_, err = c.Shutdown(ShutdownSec(sec))
		return err
	case "getShutdownTimer":
		sec, err := c.GetShutdownTimer()
		if err != nil {
			return err
		}
		bytes, err := json.Marshal(sec)
		if err != nil {
			return err
		}
		*payload = bytes
		return nil
	case "joinGroup":
		_, err := c.JoinGroup(*arg1)
		return err
	case "ungroup":
		_, err := c.Ungroup()
		return err
	case "setSource":
		_, err := c.SetSource(*arg1)
		return err
	default:
		return fmt.Errorf("Unknown command %s", name)
	}
}
