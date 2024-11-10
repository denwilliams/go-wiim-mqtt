package wiim

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/denwilliams/go-wiim-mqtt/internal/logging"
)

var transport = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

var client = &http.Client{Transport: transport}

type Device struct {
	host             string
	LastStatusEx     *StatusEx
	LastPlayerStatus *PlaybackStatus
	NextUpdateTime   int64
	UpdateInterval   uint
}

func NewDevice(host string) *Device {
	return &Device{
		host:           host,
		UpdateInterval: 5,
	}
}

func (c *Device) runCmd(name string) (io.ReadCloser, error) {
	url := "https://" + c.host + "/httpapi.asp?command=" + name
	logging.Info("Request: %s", url)

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (c *Device) runVoidCmd(name string) error {
	body, err := c.runCmd(name)
	if err != nil {
		return err
	}
	defer body.Close()

	return nil
}

func (c *Device) runStrCmd(name string) (string, error) {
	body, err := c.runCmd(name)
	if err != nil {
		return "", err
	}
	defer body.Close()

	bytes, err := io.ReadAll(body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(bytes), nil
}

func (c *Device) runIntCmd(name string) (int, error) {
	body, err := c.runCmd(name)
	if err != nil {
		return 0, err
	}
	defer body.Close()

	bytes, err := io.ReadAll(body)
	if err != nil {
		log.Fatalln(err)
	}

	str := string(bytes)

	return strconv.Atoi(str)
}

func (c *Device) runDecodedCmd(name string, respBody any) error {
	body, err := c.runCmd(name)
	if err != nil {
		return err
	}
	defer body.Close()

	if err := json.NewDecoder(body).Decode(respBody); err != nil {
		return err
	}

	return nil
}

func (c *Device) runStatusCmd(name string) (*StatusResponse, error) {
	body, err := c.runCmd(name)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var resp StatusResponse
	if err := json.NewDecoder(body).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Device) GetStatusEx() (*StatusEx, error) {
	var respBody StatusEx

	if err := c.runDecodedCmd("getStatusEx", &respBody); err != nil {
		return nil, err
	}

	c.LastStatusEx = &respBody
	return &respBody, nil
}

func (c *Device) WLANGetConnectState() (WLANGetConnectState, error) {
	return c.runStrCmd("wlanGetConnectState")
}

// Device

func (c *Device) GetPlayerStatus() (*PlaybackStatus, error) {
	var respBody PlaybackStatus
	if err := c.runDecodedCmd("getPlayerStatus", &respBody); err != nil {
		return nil, err
	}

	c.LastPlayerStatus = &respBody
	return &respBody, nil
}

func (c *Device) GetUpdatedPlayerStatus() (*PlaybackStatus, error) {
	lastStatus := c.LastPlayerStatus
	newStatus, err := c.GetPlayerStatus()
	if err != nil {
		return nil, err
	}

	diff := newStatus.GetDiff(lastStatus)

	return diff, nil
}

// Playback

func (c *Device) Pause() error {
	return c.runVoidCmd("setPlayerCmd:pause")
}

func (c *Device) Resume() error {
	return c.runVoidCmd("setPlayerCmd:resume")
}

func (c *Device) TogglePausePlay() error {
	return c.runVoidCmd("setPlayerCmd:onepause")
}

func (c *Device) Previous() error {
	return c.runVoidCmd("setPlayerCmd:prev")
}

func (c *Device) Next() error {
	return c.runVoidCmd("setPlayerCmd:next")
}

func (c *Device) Seek(seconds int) error {
	return c.runVoidCmd("setPlayerCmd:seek:" + fmt.Sprint(seconds))
}

func (c *Device) Stop() error {
	return c.runVoidCmd("setPlayerCmd:stop")
}

func (c *Device) SetLoopMode(mode PlayerLoopMode) error {
	return c.runVoidCmd("setPlayerCmd:loopmode:" + fmt.Sprint(mode))
}

func (c *Device) PlayPreset(preset int) (string, error) {
	return c.runStrCmd("MCUKeyShortClick:" + fmt.Sprint(preset))
}

// Volume

func (c *Device) SetVolume(level int) error {
	return c.runVoidCmd("setPlayerCmd:vol:" + fmt.Sprint(level))
}

func (c *Device) Mute(muted bool) error {
	return c.runVoidCmd("setPlayerCmd:vol:" + boolStr(muted))
}

// EQ

func (c *Device) EQOn() (*StatusResponse, error) {
	return c.runStatusCmd("setPlayerCmd:EQOn")
}

func (c *Device) EQOff() (*StatusResponse, error) {
	return c.runStatusCmd("setPlayerCmd:EQOff")
}

func (c *Device) EQGetStat() (*EQStatResponse, error) {
	var respBody EQStatResponse

	if err := c.runDecodedCmd("setPlayerCmd:EQGetStat", &respBody); err != nil {
		return nil, err
	}

	return &respBody, nil
}

func (c *Device) EQGetList() (*[]string, error) {
	var respBody []string
	if err := c.runDecodedCmd("setPlayerCmd:EQGetList", &respBody); err != nil {
		return nil, err
	}

	return &respBody, nil
}

func (c *Device) EQLoad(name string) (*StatusResponse, error) {
	return c.runStatusCmd("setPlayerCmd:EQLoad:" + name)
}

// Power

func (c *Device) Reboot() (*StatusResponse, error) {
	return c.runStatusCmd("setPlayerCmd:reboot")
}

func (c *Device) Shutdown(sec ShutdownSec) (*StatusResponse, error) {
	return c.runStatusCmd("setPlayerCmd:setShutdown:" + fmt.Sprint(sec))
}

func (c *Device) GetShutdownTimer() (int, error) {
	return c.runIntCmd("setPlayerCmd:getShutdown")
}

// Alarm

// TODO: alarm

// Source

func (c *Device) SetSource(source string) (*StatusResponse, error) {
	return c.runStatusCmd("setPlayerCmd:switchmode:" + source)
}

// Groups

func (c *Device) JoinGroup(masterIp string) (string, error) {
	return c.runStrCmd("ConnectMasterAp:JoinGroupMaster:eth" + masterIp)
}

func (c *Device) Ungroup() (string, error) {
	return c.runStrCmd("multiroom:Ungroup")
}

func boolStr(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

type StatusString string

const (
	StatusOK     StatusString = "OK"
	StatusFailed              = "Failed"
)

type StatusResponse struct {
	// OK or Failed
	Status string `json:"status"`
}

type ShutdownSec = int

const (
	ShutdownImmediately ShutdownSec = 0
	ShutdownCancel                  = -1
	ShutdownIn10s                   = 10
	ShutdownIn60s                   = 60
	ShutdownIn1h                    = 60 * 60
)

type PlaybackSource = string

const (
	// Line in / aux
	SourceLineIn    PlaybackSource = "line-in"
	SourceBluetooth                = "bluetooth"
	SourceOptical                  = "optical"
	SourceUDisk                    = "udisk"
	SourceWiFi                     = "wifi"
)
