package tunnel

import (
//	"fmt"
	"sync"
)

func addToEncryptionBuffer(outboundQueue chan *Packet, encryptionQueue chan *Packet, pktent *Packet) {
	outboundQueue <- pktent
	encryptionQueue <- pktent
}

func (tunnel *Tunnel) RoutineReadFromTUN(index int, max_enc int) {
	pool := make([]Packet, 15000, 15000)
	for i := 0; i < len(pool); i += 1 {
		pool[i].buffer = make([]byte, 2000, 2000)
		pool[i].Mutex = sync.Mutex{}
		pool[i].Lock()
	}
	var pos int = 0
	var enc int = 0
	for {
		pkt := pool[pos % len(pool)]
		size, _ := tunnel.tun.tunnel.Read(index, pkt.buffer[:])
		pkt.packet = pkt.buffer[:size]
		addToEncryptionBuffer(tunnel.queue.outbound[index], tunnel.queue.encryption[index][enc % max_enc], &pkt)
		pos += 1
		enc += 1
	}
}

func (tunnel *Tunnel) RoutineEncryption(queue int, enc int) {
	for {
		pkt, _ := <-tunnel.queue.encryption[queue][enc]
		// encrypt packet
		pkt.Unlock()
	}
}

func (tunnel *Tunnel) RoutineWriteToUDP(index int) {
	for {
		pkt, _ := <-tunnel.queue.outbound[index]
		pkt.Lock()
		tunnel.Send(index, pkt.packet)
	}
}
