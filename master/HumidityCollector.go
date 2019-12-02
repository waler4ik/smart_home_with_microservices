package main

import "context"
import "fmt"
import proto "protofiles"
import micro "github.com/micro/go-micro"

//Collector is ...
type Collector struct {

}

//Collect is a callback function
func (collector *Collector) Collect(ctx context.Context, request *proto.CollectRequest, response *proto.CollectResponse) error {
  fmt.Printf("Node: %s, Humidity: %f, Tmp: %f\n", request.GetNodeName(), request.GetHumidity(), request.GetTemperature())
  response.StatusCode = 0
  response.StatusMessage = "Nice"
  return nil
}

func main() {
  service := micro.NewService(micro.Name("Humidity_Value_Collector"))
  service.Init()
  proto.RegisterHumidityCollectorHandler(service.Server(), new(Collector))
  if err := service.Run(); err != nil {
    fmt.Println(err)
  }
}
