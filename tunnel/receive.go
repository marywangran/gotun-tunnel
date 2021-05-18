package tunnel

import (
	"sync"
	//"fmt"
)

func addToDecryptionBuffer(inboundQueue chan *Packet, decryptionQueue chan *Packet, pktent *Packet) {
	inboundQueue <- pktent
	decryptionQueue <- pktent
}

func (tunnel *Tunnel) RoutineReadFromUDP(queue int, max_enc int) {
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
		size := tunnel.Receive(queue, pkt.buffer[:])
		if pkt.buffer[0] == 'H' {
			continue
		}
		pkt.packet = pkt.buffer[:size]
		addToDecryptionBuffer(tunnel.queue.inbound[queue], tunnel.queue.decryption[queue][enc % max_enc], &pkt)
		pos += 1
		enc += 1
	}
}

func (tunnel *Tunnel) RoutineDecryption(queue int, enc int) {
	for {
		pkt, _ := <-tunnel.queue.decryption[queue][enc]
		// decrypt packet
		pkt.Unlock()
	}
}

func (tunnel *Tunnel) RoutineWriteToTUN(index int) {
	for {
		pkt, _ := <-tunnel.queue.inbound[index]
		pkt.Lock()
		tunnel.tun.tunnel.Write(index, pkt.buffer[:len(pkt.packet)])
	}
}
