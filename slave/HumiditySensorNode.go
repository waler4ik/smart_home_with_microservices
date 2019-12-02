package main

import (
	"context"
	"fmt"
	"time"
	"os"
	"encoding/json"

	micro "github.com/micro/go-micro"
	proto "protofiles"
  dht "github.com/d2r2/go-dht"
)

//Configuration is ..
type Configuration struct {
	Location string
	ReadPeriod time.Duration
}

func collect(collector proto.HumidityCollectorService, deathSignal chan bool, config *Configuration) {
	duration := config.ReadPeriod * time.Minute
	for {
		temperature, humidity, _, err := dht.ReadDHTxxWithRetry(dht.DHT22, 4, false, 10)
		if err != nil {
			fmt.Println(err)
			deathSignal <- true
			return
		}

		rsp, err := collector.Collect(context.TODO(), &proto.CollectRequest{NodeName: config.Location, Humidity: humidity, Temperature: temperature})

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(rsp)
		time.Sleep(duration)
	}
}

func main() {
	service := micro.NewService(micro.Name("Humidity_Value_Creator"))
	service.Init()

	collector := proto.NewHumidityCollectorService("Humidity_Value_Collector", service.Client())

	deathSignal := make(chan bool)

	file, err := os.Open("/home/pi/config.json")
	if err != nil {
		fmt.Println(err)
	}
	decoder := json.NewDecoder(file)
	configuration := Configuration{Location: "Keller", ReadPeriod: 1}
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Println(err)
	}

  go collect(collector, deathSignal, &configuration)

	fmt.Printf("I'm done with work: %v", <-deathSignal)
}
