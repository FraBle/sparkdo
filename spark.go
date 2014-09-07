package main

import (
	// standard library
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	// external packages
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/spark"
)

func monitorDroplet(dropletId, sparkDeviceId, sparkAccessToken, digitalOceanAccessToken string, client *http.Client) {
	gbot := gobot.NewGobot()
	sparkCore := spark.NewSparkCoreAdaptor("spark", sparkDeviceId, sparkAccessToken)
	redLed := gpio.NewLedDriver(sparkCore, "led", "D0")
	greenLed := gpio.NewLedDriver(sparkCore, "led", "D1")
	work := func() {
		gobot.Every(5*time.Second, func() {
			dropletReq, err := http.NewRequest("GET", "https://api.digitalocean.com/v2/droplets/"+dropletId, nil)
			if err != nil {
				log.Printf("Error %v", err)
				return
			}
			dropletReq.Header.Add("Authorization", "Bearer "+digitalOceanAccessToken)
			dropletResp, err := client.Do(dropletReq)
			if err != nil {
				log.Printf("Error %v", err)
				return
			}
			defer dropletResp.Body.Close()
			data, err := ioutil.ReadAll(dropletResp.Body)
			if err != nil {
				log.Printf("Error %v", err)
				return
			}
			var droplet SingleDroplet
			err = json.Unmarshal(data, &droplet)
			if err != nil {
				log.Printf("Error %v", err)
				return
			}
			if droplet.Droplet.Status == "active" {
				greenLed.On()
				redLed.Off()
			} else {
				greenLed.Off()
				redLed.On()
			}
		})
	}
	robot := gobot.NewRobot("spark",
		[]gobot.Connection{sparkCore},
		[]gobot.Device{redLed, greenLed},
		work,
	)
	gbot.AddRobot(robot)
	gbot.Start()
}
