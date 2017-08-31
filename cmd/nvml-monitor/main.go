package main

import (
	"context"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/hotpxl/nvml"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("nvml-monitor", "Monitor NVML status and upload to etcd.")
	duration := app.Flag("duration", "Duration before statistics report.").Default("5s").Duration()
	endpoints := app.Flag("endpoints", "Etcd cluster endpoints to connect to.").Required().Strings()
	base := app.Flag("base", "Base path of etcd.").Default("/").String()
	kingpin.MustParse(app.Parse(os.Args[1:]))

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   *endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.WithField("err", err).Fatal("Failed to connect to etcd.")
	}
	defer cli.Close()

	session, err := nvml.NewSession()
	if err != nil {
		log.WithError("err", err).Fatal("Failed to create NVML session.")
	}
	defer session.Close()

	hostname, err := os.Hostname()
	if err != nil {
		log.WithError("err", err).Fatal("Failed to retrieve hostname.")
	}

	for {
		devices, err := session.GetAllDevices()
		if err != nil {
			log.WithError("err", err).Fatal("Failed to get devices.")
		}
		for idx, d := range devices {
			mem, err := d.MemoryInfo()
			if err != nil {
				log.WithError("err", err).Fatal("Failed to get memory information.")
			}
			_, err = cli.Put(context.Background(), path.Join(base, hostname, strconv.Itoa(idx), "mem", "free"), strconv.Itoa(mem.Free))
			if err != nil {
				log.WithError("err", err).Fatal("Failed to upload memory information.")
			}
			_, err = cli.Put(context.Background(), path.Join(base, hostname, strconv.Itoa(idx), "mem", "used"), strconv.Itoa(mem.Used))
			if err != nil {
				log.WithError("err", err).Fatal("Failed to upload memory information.")
			}
			_, err = cli.Put(context.Background(), path.Join(base, hostname, strconv.Itoa(idx), "mem", "total"), strconv.Itoa(mem.Total))
			if err != nil {
				log.WithError("err", err).Fatal("Failed to upload memory information.")
			}
			processes, err := d.Processes()
			if err != nil {
				log.WithError("err", err).Fatal("Failed to get processes.")
			}
			for _, p := range processes {
				_, err = cli.Put(context.Background(), path.Join(base, hostname, strconv.Itoa(idx), "proc", strconv.Itoa(p.PID), "used_memory"), strconv.Itoa(p.UsedMemory))
				if err != nil {
					log.WithError("err", err).Fatal("Failed to upload process information.")
				}
				_, err = cli.Put(context.Background(), path.Join(base, hostname, strconv.Itoa(idx), "proc", strconv.Itoa(p.PID), "username"), p.Username)
				if err != nil {
					log.WithError("err", err).Fatal("Failed to upload process information.")
				}
			}
		}
		time.Sleep(duration)
	}
}