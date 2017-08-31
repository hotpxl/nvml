package nvml

import (
	"fmt"

	"github.com/shirou/gopsutil/process"
	log "github.com/sirupsen/logrus"
)

// Session manages the initialization and shutdown of NVML context.
type Session struct {
	active bool
}

// NewSession creates a new session.
func NewSession() (*Session, error) {
	return &Session{active: true}, nvmlInit()
}

// Close frees existing session and underlying NVML context.
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

// DeviceCount returns the number of devices.
func (s *Session) DeviceCount() (int, error) {
	if !s.active {
		return 0, fmt.Errorf("Already closed.")
	}
	return nvmlDeviceGetCount()
}

// GetDevice returns a specific device given its index.
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

// GetAllDevices returns all devices accessible.
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

// Device represents a single device.
type Device struct {
	handle deviceHandle
}

// MemoryInfo holds memory consumption information for a device.
type MemoryInfo struct {
	Free  uint64
	Used  uint64
	Total uint64
}

// ProcessInfo holds process information on a device.
type ProcessInfo struct {
	PID        int32  `json:"pid"`
	UsedMemory uint64 `json:"usedMemory"`
	Username   string `json:"username"`
}

// MemoryInfo returns memory consumption information from a device.
func (d *Device) MemoryInfo() (MemoryInfo, error) {
	return nvmlDeviceGetMemoryInfo(d.handle)
}

// Processes returns processes running on a device.
func (d *Device) Processes() ([]ProcessInfo, error) {
	processes, err := nvmlDeviceGetComputeRunningProcesses(d.handle)
	if err != nil {
		return nil, err
	}
	for idx, p := range processes {
		pp, err := process.NewProcess(p.PID)
		if err != nil {
			return nil, err
		}
		username, err := pp.Username()
		if err != nil {
			return nil, err
		}
		processes[idx].Username = username
	}
	return processes, nil
}
