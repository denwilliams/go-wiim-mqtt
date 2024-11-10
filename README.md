# go-wiim-mqtt

Early version, limited error handling, etc.

## Configuration

`MQTT_URI` - URI to the MQTT broker

`MQTT_TOPIC_PREFIX` - topic prefix for the MQTT messages, eg `wiim`

`PORT` - port to expose HTTP server for metrics on

`WIIM_IPS` - comma separated list of IP addresses for the WiiM devices

## WiiM

https://developer.arylic.com/httpapi/

https://www.wiimhome.com/pdf/HTTP%20API%20for%20WiiM%20Mini.pdf

https://github.com/AndersFluur/LinkPlayApi/blob/master/api.md

## MQTT

When changes are detected on the WiiM devices, a message is published to the MQTT broker under `/status/{name of device}/*`.

eg:

/status/Living Room/StatusEx contains the fields from getStatusEx
/status/Living Room/PlayerStatus contains the fields from getPlayerStatus
/status/Living Room/PlayerStatusVolume contains only the volume level when changed
/status/Living Room/PlayerStatusMute contains only the mute status when changed
/status/Living Room/PlayerStatusType contains only the type when changed (`master`/`slave`)
/status/Living Room/PlayerStatusMode contains only the mode when changed (`default`/`spotify_connect`,`chromecast`,`airplay`,etc)
/status/Living Room/PlayerStatusStatus contains only the status when changed (`play`/`stop`/`pause`)
/status/Living Room/PlayerStatusPlaying contains true when when status is `play`, else false

## Prometheus Metrics

GET from /metrics on the HTTP server

Not much there now, will add more later