// goでdnsサーバー
// https://github.com/EmilHernvall/dnsguide
// qname圧縮に関する仕様の解説
// http://park12.wakwak.com/~eslab/pcmemo/dns/dns5.html#condense

package main

import (
	"fmt"
	"os"

	"github.com/k0kubun/pp/v3"
)

func main() {
	file, err := os.Open("response_packet.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	buffer := NewBytePacketBuffer()
	_, err = file.Read(buffer.buf[:])
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	packet, err := ReadDnsPacket(buffer)
	if err != nil {
		fmt.Println("Error reading DNS packet:", err)
		return
	}

	pp.Print(packet.Header)

	for _, q := range packet.Questions {
		pp.Print(q)
	}
	for _, rec := range packet.Answers {
		pp.Print(rec)
	}
	for _, rec := range packet.Authorities {
		pp.Print(rec)
	}
	for _, rec := range packet.Resources {
		pp.Print(rec)
	}
}
