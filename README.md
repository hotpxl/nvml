# NVIDIA Management Library (NVML) Go Bindings

[![GoDoc](https://godoc.org/github.com/hotpxl/nvml?status.svg)](https://godoc.org/github.com/hotpxl/nvml)

There are multiple NVML Go bindings lying around GitHub. But they
either are unmaintained or require configuring compiler flags. This
package uses cgo and aims to be usable without any configuration.

## Example

Following is an easy example that displays processes information on all devices.

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