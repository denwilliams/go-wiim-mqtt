package wiim

type PlaybackStatus struct {
	Type PlaybackStatusType `json:"type,string"`
	Ch   PlaybackStatusCh   `json:"ch,string"`
	Mode PlaybackStatusMode `json:"mode,string"`
	Loop int                `json:"loop,string"`
	// The preset number of the Equalizer
	Eq        int         `json:"eq,string"`
	Status    string      `json:"status"`
	CurPos    int         `json:"curpos,string"`
	OffsetPts int         `json:"offset_pts,string"`
	TotLen    int         `json:"totlen,string"`
	AlarmFlag int         `json:"alarmflag,string"`
	PliCount  int         `json:"plicount,string"`
	PliCurr   int         `json:"plicurr,string"`
	Vol       int         `json:"vol,string"`
	VolEdit   int         `json:"voledit,string"`
	Mute      JsonBoolean `json:"mute"`
}

func (ps *PlaybackStatus) GetDiff(comparison *PlaybackStatus) *PlaybackStatus {
	result := PlaybackStatus{}

	dirty := false
	if ps.Type != comparison.Type {
		dirty = true
		result.Type = ps.Type
	}
	if ps.Ch != comparison.Ch {
		dirty = true
		result.Ch = ps.Ch
	}
	if ps.Mode != comparison.Mode {
		dirty = true
		result.Mode = ps.Mode
	}
	if ps.Loop != comparison.Loop {
		dirty = true
		result.Loop = ps.Loop
	}
	if ps.Eq != comparison.Eq {
		dirty = true
		result.Eq = ps.Eq
	}
	if ps.Status != comparison.Status {
		dirty = true
		result.Status = ps.Status
	}
	if ps.CurPos != comparison.CurPos {
		dirty = true
		result.CurPos = ps.CurPos
	}
	if ps.OffsetPts != comparison.OffsetPts {
		dirty = true
		result.OffsetPts = ps.OffsetPts
	}
	if ps.TotLen != comparison.TotLen {
		dirty = true
		result.TotLen = ps.TotLen
	}
	if ps.AlarmFlag != comparison.AlarmFlag {
		dirty = true
		result.AlarmFlag = ps.AlarmFlag
	}
	if ps.PliCount != comparison.PliCount {
		dirty = true
		result.PliCount = ps.PliCount
	}
	if ps.PliCurr != comparison.PliCurr {
		dirty = true
		result.PliCurr = ps.PliCurr
	}
	if ps.VolEdit != comparison.VolEdit {
		dirty = true
		result.VolEdit = ps.VolEdit
	}
	if ps.Mute != comparison.Mute {
		dirty = true
		result.Mute = ps.Mute
	}

	if !dirty {
		return nil
	}
	return &result
}

type PlaybackStatusType = int

const (
	Master PlaybackStatusType = iota // 0 = master speaker
	Slave                            // 1 = slave speaker in a group
)

type PlaybackStatusCh = int

const (
	Stereo PlaybackStatusCh = iota // 0 = stereo
	Left                           // 1 = left
	Right                          // 2 = right
)

type PlaybackStatusMode = int

const (
	ModeNone       PlaybackStatusMode = 0 // none
	ModeAirplay                       = 1 // AirPlay
	ModeDLNA                          = 2 // 3rd party DLNA
	ModeChromecast                    = 5 // Undocumented
	// this is the mode for Wifi sources like Soundcloud, Plex, etc
	ModePlaylistDefault = 10 // Wiimu playlist - default wiimu mode
	ModePlaylistUSB     = 11 // Wiimu playlist - USB disk playlist
	ModePlaylistTF      = 16 // Wiimu playlist - TF
	ModeSpotifyConnect  = 31 // Spotify Connect
	ModeTidalConnect    = 32 // Tidal Connect
	ModeAuxIn           = 40 // Aux-In
	ModeBluetooth       = 41 // BT
	ModeExternalStorage = 42 // external storage
	ModeOpticalIn       = 43 // Optical-In
	ModeMirror          = 50 // Mirror
	ModeVoiceMail       = 60 // Voice mail
	ModeSlave           = 99 // Slave
)

var PlayerStatusModeName = map[PlaybackStatusMode]string{
	ModeNone:            "none",
	ModeAirplay:         "airplay",
	ModeDLNA:            "dlna",
	ModeChromecast:      "chromecast",
	ModePlaylistDefault: "default",
	ModePlaylistUSB:     "usb",
	ModePlaylistTF:      "tf",
	ModeSpotifyConnect:  "spotify_connect",
	ModeTidalConnect:    "tidal_connect",
	ModeAuxIn:           "aux_in",
	ModeBluetooth:       "bluetooth",
	ModeExternalStorage: "external_storage",
	ModeOpticalIn:       "optical_in",
	ModeMirror:          "mirror",
	ModeVoiceMail:       "voice_mail",
	ModeSlave:           "slave",
}

type PlayerLoopMode = int

const (
	LoopNone     PlayerLoopMode = 0  // Sequence, no loop
	LoopSingle   PlayerLoopMode = 1  // Single loop
	LoopShuffle  PlayerLoopMode = 2  // Shuffle loop
	LoopSequence PlayerLoopMode = -1 // Sequence loop
)

type EQStat = string

const (
	EQStatOn  EQStat = "On"
	EQStatOff        = "Off"
)

type EQStatResponse struct {
	EQStat EQStat `json:"EQStat"`
}
