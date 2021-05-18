package tunnel

import (
	"runtime"
	"sync"
	"tuntap/tun"
)

type Packet struct {
	sync.Mutex
	buffer   []byte
	packet   []byte
}

type Tunnel struct {
	isClient  bool

	WG sync.WaitGroup
	state struct {
		sync.Mutex
	}

	net struct {
		socket	*UDPScoket
		port    int
		addr    [4]byte
	}

	queue struct {
		inbound []chan *Packet
		outbound []chan *Packet
		encryption [][]chan *Packet
		decryption [][]chan *Packet
	}

	tun struct {
		tunnel tun.Device
		queues    int
	}
}

func NewInstance(tunTunnel tun.Device, isClient bool, queues int) *Tunnel {
	tunnel := new(Tunnel)

	tunnel.isClient = isClient

	tunnel.tun.queues = queues
	tunnel.tun.tunnel = tunTunnel
	tunnel.net.port = 12346
	tunnel.net.addr = [4]byte{10, 198, 54, 67}

	if tunnel.isClient {
		tunnel.net.socket = CreateUDPScoket(tunnel.net.port, tunnel.net.addr, tunnel.tun.queues, 1)
	} else {
		tunnel.net.socket = CreateUDPScoket(tunnel.net.port, tunnel.net.addr, tunnel.tun.queues, 0)
	}

	tunnel.queue.outbound = make([]chan *Packet, queues)
	tunnel.queue.inbound = make([]chan *Packet, queues)

	enc := runtime.NumCPU()/queues
	tunnel.queue.encryption = make([][]chan *Packet, queues)
	tunnel.queue.decryption = make([][]chan *Packet, queues)

	for i := 0; i < queues; i += 1 {
		tunnel.queue.outbound[i] = make(chan *Packet, 15000)
		tunnel.queue.inbound[i] = make(chan *Packet, 15000)
		tunnel.queue.encryption[i] = make([]chan *Packet, enc)
		tunnel.queue.decryption[i] = make([]chan *Packet, enc)
		for j := 0; j < enc; j += 1 {
			tunnel.queue.encryption[i][j] = make(chan *Packet, 8000)
			tunnel.queue.decryption[i][j] = make(chan *Packet, 8000)
			go tunnel.RoutineDecryption(i, j)
			go tunnel.RoutineEncryption(i, j)
		}
		go tunnel.RoutineReadFromUDP(i, enc)
		go tunnel.RoutineWriteToTUN(i)
		go tunnel.RoutineReadFromTUN(i, enc)
		go tunnel.RoutineWriteToUDP(i)
	}
	tunnel.WG.Add(1)

	return tunnel
}
