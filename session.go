package nvml

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Session struct {
	active bool
}

func NewSession() (*Session, error) {
	return &Session{active: true}, nvmlInit()
}

func (s *Session) Close() {
	if !s.active {
		log.Fatal("Already closed.")
	}
	s.active = false
	err := nvmlShutdown()
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Session) DeviceCount() (int, error) {
	if !s.active {
		return 0, fmt.Errorf("Already closed.")
	}
	return nvmlDeviceGetCount()
}

func (s *Session) GetDevice(idx int) (*Device, error) {
	if !s.active {
		return nil, fmt.Errorf("Already closed.")
	}
	dev, err := nvmlDeviceGetHandleByIndex(idx)
	if err != nil {
		return nil, err
	}
	return &Device{handle: dev}, nil
}

func (s *Session) GetAllDevices() ([]Device, error) {
	if !s.active {
		return nil, fmt.Errorf("Already closed.")
	}
	count, err := s.DeviceCount()
	if err != nil {
		return nil, err
	}
	var ret []Device
	for i := 0; i < count; i++ {
		dev, err := s.GetDevice(i)
		if err != nil {
			return nil, err
		}
		ret = append(ret, *dev)
	}
	return ret, nil
}

type Device struct {
	handle deviceHandle
}

type MemoryInfo struct {
	Free  uint64
	Used  uint64
	Total uint64
}

func (d *Device) MemoryInfo() (MemoryInfo, error) {
}
