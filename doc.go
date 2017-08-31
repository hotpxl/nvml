/*
Package nvml is the NVIDIA Management Library (NVML) bindings for Go.

Following is an easy example that displays processes information on
all devices.

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

Visit https://github.com/hotpxl/nvml for more information.
*/
package nvml
