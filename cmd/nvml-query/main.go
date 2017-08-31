package main

import (
	"fmt"

	"github.com/hotpxl/nvml"
)

func main() {
	s, err := nvml.NewSession()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	devices, err := s.GetAllDevices()
	if err != nil {
		panic(err)
	}
	for _, d := range devices {
		p, err := d.Processes()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", p)
	}
}
