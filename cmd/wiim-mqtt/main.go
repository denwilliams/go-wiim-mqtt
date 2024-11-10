package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/denwilliams/go-wiim-mqtt/internal/logging"
	"github.com/denwilliams/go-wiim-mqtt/internal/mqtt"
	"github.com/denwilliams/go-wiim-mqtt/internal/web"
	"github.com/denwilliams/go-wiim-mqtt/internal/wiim"
	"github.com/joho/godotenv"
)

func init() {
	logging.Init(nil, logging.DefaultFlags)
	logging.Info("Loading .env file")
	err := godotenv.Load(".env")

	if err != nil {
		logging.Warn("Unable to load .env")
	}
}

func main() {
	mu, err := url.Parse(os.Getenv("MQTT_URI"))
	if err != nil {
		logging.Error("Error parsing URL %s", err)
	}
	baseTopic := os.Getenv("MQTT_TOPIC_PREFIX")
	subscribeTopic := os.ExpandEnv("$MQTT_TOPIC_PREFIX/set/#")
	portStr := os.Getenv("PORT")
	serverPort, _ := strconv.Atoi(portStr)
	if err != nil {
		logging.Error("Error parsing PORT %s", err)
	}
	ipAddresses := strings.Split(os.Getenv("WIIM_IPS"), ",")

	mc := mqtt.NewMQTTClient(mu, baseTopic, subscribeTopic)
	wc := wiim.NewMuxClient()
	se := mqtt.NewMqttStatusEmitter(mc)
	for _, ip := range ipAddresses {
		wc.AddDevice(ip, wiim.NewDevice(ip))
	}
	mc.Connect(wc)
	defer mc.Disconnect()

	go statusLoop(wc, se)

	if serverPort > 0 {
		go startServer(serverPort)
	}

	logging.Info("Ready")

	waitForExit()

	logging.Info("Terminating")
}

func waitForExit() {
	// Set up a channel to receive OS signals so we can gracefully exit
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	logging.Info("Exit signal received")
}

func statusLoop(wc *wiim.MuxClient, se *mqtt.MqttStatusEmitter) {
	// Set up a channel to receive OS signals so we can gracefully exit
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// logging.Info("Performing initial discovery")
	// // It can take a few runs to discover all the lights
	// // Keep going until we find no new lights for a few runs
	// emptyRuns := 0
	// for emptyRuns < 10 {
	// 	found := lc.DiscoverWithTimeout(15 * time.Second)
	// 	if found == 0 {
	// 		emptyRuns++
	// 	} else {
	// 		emptyRuns = 0
	// 	}
	// }
	// logging.Info("Finished initial light discovery, will continue to discover every 10 minutes")

	// We want to continually call the Discover method at an interval
	// to pick up on new lights that come online
	tick := time.Tick(5 * time.Second)

	updateStatuses(wc, se)

	for {
		select {
		case <-tick:
			updateStatuses(wc, se)
		case <-signalChan:
			// Stop the loop when an interrupt signal is received
			logging.Info("Background discovery loop interrupted, exiting")
			return
		}
	}
}

// Max returns the larger of x or y.
func min(x, y uint) uint {
	if x < y {
		return x
	}
	return y
}

func updateStatuses(wc *wiim.MuxClient, se *mqtt.MqttStatusEmitter) {
	for device := range wc.GetDevicesForUpdate() {
		if device.LastPlayerStatus == nil {
			updateStatus(device, se)
		}
		didChange := updatePlayerStatus(device, se)
		if didChange {
			device.UpdateInterval = 5
		} else {
			// back off polling frequency
			device.UpdateInterval = min(60, device.UpdateInterval+5)
		}
		device.NextUpdateTime = time.Now().Add(time.Duration(device.UpdateInterval) * time.Second).Unix()
	}
}

func updateStatus(device *wiim.Device, se *mqtt.MqttStatusEmitter) {
	status, err := device.GetStatusEx()
	if err != nil {
		logging.Error("Error getting status: %s", err)
		return
	}
	logging.Info("Emitting status: %s %d %d", status.DeviceName, status.Group, status.PowerMode)
	se.EmitStatus(context.TODO(), status.DeviceName, "StatusEx", status)
}

func updatePlayerStatus(device *wiim.Device, se *mqtt.MqttStatusEmitter) bool {
	last := device.LastPlayerStatus
	status, err := device.GetPlayerStatus()
	if err != nil {
		logging.Error("Error getting status: %s", err)
		return false
	}

	ctx := context.TODO()

	if last == nil || last.Type != status.Type {
		typeName := "master"
		if status.Type == 1 {
			typeName = "slave"
		}
		logging.Info("Type changed: %s %d %s", device.LastStatusEx.DeviceName, status.Type, typeName)
		se.EmitStatus(ctx, device.LastStatusEx.DeviceName, "PlayerStatusType", typeName)
	}
	if last == nil || last.Mode != status.Mode {
		modeName, ok := wiim.PlayerStatusModeName[status.Mode]
		if ok {
			logging.Info("Mode changed: %s %d %s", device.LastStatusEx.DeviceName, status.Mode, modeName)
			se.EmitStatus(ctx, device.LastStatusEx.DeviceName, "PlayerStatusMode", modeName)
		} else {
			logging.Warn("Unknown Mode: %s %d", device.LastStatusEx.DeviceName, status.Mode)
		}
	}
	if last == nil || last.Status != status.Status {
		logging.Info("Status changed: %s %s", device.LastStatusEx.DeviceName, status.Status)
		se.EmitStatus(ctx, device.LastStatusEx.DeviceName, "PlayerStatusStatus", status.Status)
		se.EmitStatus(ctx, device.LastStatusEx.DeviceName, "PlayerStatusPlaying", status.Status == "play")
	}
	if last == nil || last.Vol != status.Vol {
		logging.Info("Volume changed: %s %d", device.LastStatusEx.DeviceName, status.Vol)
		se.EmitStatus(ctx, device.LastStatusEx.DeviceName, "PlayerStatusVolume", status.Vol)
	}
	if last == nil || last.Mute != status.Mute {
		logging.Info("Mute changed: %s %t", device.LastStatusEx.DeviceName, status.Mute)
		se.EmitStatus(ctx, device.LastStatusEx.DeviceName, "PlayerStatusMute", status.Mute)
	}

	// check if last is different to status
	if last != nil && last.Type == status.Type && last.Mode == status.Mode && last.Status == status.Status && last.Vol == status.Vol && last.CurPos == status.CurPos {
		return false
	}

	logging.Info("Emitting player status: %s type=%d mode=%d status=%s vol=%d curpos=%d", device.LastStatusEx.DeviceName, status.Type, status.Mode, status.Status, status.Vol, status.CurPos)
	se.EmitStatus(context.TODO(), device.LastStatusEx.DeviceName, "PlayerStatus", status)
	return true
}

func startServer(port int) {
	logging.Info("Creating HTTP server")
	handler := web.CreateHandler()
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}
	logging.Info("Starting HTTP server on port %d", port)
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("error running http server: %s\n", err)
		}
		log.Fatal(err)
	}
}
