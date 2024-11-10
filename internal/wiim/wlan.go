package wiim

type WLANGetConnectState = string

const (
	PROCESS  WLANGetConnectState = "PROCESS"
	PAIRFAIL                     = "PAIRFAIL"
	FAIL                         = "FAIL"
	OK                           = "OK"
)
