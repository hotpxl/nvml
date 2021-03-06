package nvml

// #cgo CFLAGS: -I/usr/local/cuda/include
// #cgo LDFLAGS: -lnvidia-ml
// #include <nvml.h>
import "C"
import (
	"fmt"
	"strconv"
)

type deviceHandle = C.nvmlDevice_t

func handleError(ret C.nvmlReturn_t) error {
	if ret == C.NVML_SUCCESS {
		return nil
	}
	err := C.GoString(C.nvmlErrorString(ret))
	return fmt.Errorf("NVML error: %s.", strconv.QuoteToASCII(err))
}

func nvmlInit() error {
	return handleError(C.nvmlInit())
}

func nvmlShutdown() error {
	return handleError(C.nvmlShutdown())
}

func nvmlDeviceGetCount() (int, error) {
	var n C.uint
	ret := C.nvmlDeviceGetCount(&n)
	return int(n), handleError(ret)
}

func nvmlDeviceGetHandleByIndex(idx int) (deviceHandle, error) {
	var dev deviceHandle
	ret := C.nvmlDeviceGetHandleByIndex(C.uint(idx), &dev)
	return dev, handleError(ret)
}

func nvmlDeviceGetMemoryInfo(h deviceHandle) (MemoryInfo, error) {
	var mem C.nvmlMemory_t
	ret := C.nvmlDeviceGetMemoryInfo(h, &mem)
	return MemoryInfo{Free: uint64(mem.free), Used: uint64(mem.used), Total: uint64(mem.total)}, handleError(ret)
}

func nvmlDeviceGetComputeRunningProcesses(h deviceHandle) ([]ProcessInfo, error) {
	count := C.uint(64)
	var processes [64]C.nvmlProcessInfo_t
	ret := C.nvmlDeviceGetComputeRunningProcesses(h, &count, &processes[0])
	err := handleError(ret)
	if err != nil {
		return nil, err
	}
	var result []ProcessInfo
	for i := 0; i < int(count); i++ {
		result = append(result, ProcessInfo{
			PID:        int32(processes[i].pid),
			UsedMemory: uint64(processes[i].usedGpuMemory),
		})
	}
	return result, nil
}
