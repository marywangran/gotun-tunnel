package main

import (
//	"runtime"
	"os"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"tuntap/tunnel"
	"tuntap/tun"
)

func init() {
//	runtime.LockOSThread()
//	runtime.GOMAXPROCS(48)
}

func main() {
	var client bool = false
	var queues int = 4

	if len(os.Args) == 3 {
		if os.Args[1] == "client" {
			client = true
		}
		queues, _ = strconv.Atoi(os.Args[2])
	}
	go func() {
		log.Println(http.ListenAndServe("10.198.54.67:6061", nil))
	}()

	tun := func() (tun.Device) {
		return tun.CreateTUN("wg2", 1500, queues)
	} ()

	instance := tunnel.NewInstance(tun, client, queues)
	instance.WG.Wait()
}
