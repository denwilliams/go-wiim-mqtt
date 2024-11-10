package wiim

type StatusExGroup = int

const (
	GroupMaster StatusExGroup = iota // 0
	GroupSlave                       // 1
)

type MACAddress string

type StatusExSecurityCapabilities struct {
	Ver    string `json:"ver"`
	AESVer string `json:"aes_ver"`
}

type StatusEx struct {
	Language string `json:"language"`
	// Name of the device
	SSID     string      `json:"ssid"`
	HideSSID JsonBoolean `json:"hideSSID"`
	// firmware version
	Firmware string `json:"firmware"`
	Build    string `json:"build"`
	Project  string `json:"project"`
	PrivPrj  string `json:"priv_prj"`
	// data the firmware is released
	Release string `json:"Release"`
	// Reserved
	FWReleaseVersion string `json:"FW_Release_version"`
	// 0 means it's a master speaker, 1 means a slave speaker in a group
	Group StatusExGroup `json:"group,string"`
	// LinkPlay's MRM SDK version, version 4.2 or above won't work with any version below 4.2
	WMRMVersion string `json:"wmrm_version"`
	// Reserved
	Expired string `json:"expired"`
	// Is it connected to Internet
	Internet JsonBoolean `json:"internet"`
	UUID     MACAddress  `json:"uuid"`
	MAC      MACAddress  `json:"MAC"`
	BT_MAC   MACAddress  `json:"BT_MAC"`
	// The MAC address of the AP that the device is connected to
	AP_MAC string `json:"AP_MAC"`
	// eg 2022:08:09"
	Date                  string                       `json:"date"`
	Time                  string                       `json:"time"`
	Netstat               string                       `json:"netstat"`
	ESSID                 string                       `json:"essid"`
	apcli0                string                       `json:"apcli0"`
	Eth0                  string                       `json:"eth0"`
	ETHMAC                string                       `json:"ETH_MAC"`
	Hardware              string                       `json:"hardware"`
	VersionUpdate         string                       `json:"VersionUpdate"`
	NewVer                string                       `json:"NewVer"`
	mcu_ver               string                       `json:"mcu_ver"`
	mcu_ver_new           string                       `json:"mcu_ver_new"`
	update_check_count    string                       `json:"update_check_count"`
	ra0                   string                       `json:"ra0"`
	temp_uuid             string                       `json:"temp_uuid"`
	cap1                  string                       `json:"cap1"`
	Capability            string                       `json:"capability"`
	Languages             string                       `json:"languages"`
	prompt_status         int                          `json:"prompt_status,string"`
	alexa_ver             string                       `json:"alexa_ver"`
	alexa_beta_enable     string                       `json:"alexa_beta_enable"`
	alexa_force_beta_cfg  string                       `json:"alexa_force_beta_cfg"`
	DSPVer                string                       `json:"dsp_ver"`
	StreamsAll            string                       `json:"streams_all"`
	Streams               string                       `json:"streams"`
	Region                string                       `json:"region"`
	VolumeControl         string                       `json:"volume_control"`
	External              string                       `json:"external"`
	PresetKey             int                          `json:"preset_key,string"`
	plm_support           string                       `json:"plm_support"`
	lbc_support           string                       `json:"lbc_support"`
	WifiChannel           int                          `json:"WifiChannel,string"`
	RSSI                  int                          `json:"RSSI,string"`
	BSSID                 string                       `json:"BSSID"`
	WLANFrequency         string                       `json:"wlanFrequency"`
	WLANDataRate          string                       `json:"wlanDataRate"`
	Battery               string                       `json:"battery"`
	BatteryPercent        string                       `json:"battery_percent"`
	SecureMode            string                       `json:"securemode"`
	OTAInterfaceVer       string                       `json:"ota_interface_ver"`
	UPNPVersion           string                       `json:"upnp_version"`
	UPNPUUID              string                       `json:"upnp_uuid"`
	UARTPassPort          string                       `json:"uart_pass_port"`
	CommunicationPort     string                       `json:"communication_port"`
	WebFirmwareUpdateHide string                       `json:"web_firmware_update_hide"`
	TidalVersion          string                       `json:"tidal_version"`
	ServiceVersion        string                       `json:"service_version"`
	EQSupport             string                       `json:"EQ_support"`
	HiFiSRCVersion        string                       `json:"HiFiSRC_version"`
	PowerMode             int                          `json:"power_mode,string"`
	Security              string                       `json:"security"`
	SecurityVersion       string                       `json:"security_version"`
	SecurityCapabilities  StatusExSecurityCapabilities `json:"security_capabilities"`
	PublicHTTPSVersion    string                       `json:"public_https_version"`
	PrivacyMode           string                       `json:"privacy_mode"`
	// The device name
	DeviceName string `json:"DeviceName"`
	// The group name of the device is belonged to
	GroupName string `json:"GroupName"`
}
