package main

import (
	"errors"
	"fmt"
	"net"
)

type DnsRecord interface {
	getType() int
	Write(*BytePacketBuffer)(int, error)
}

type UnknownRecord struct {
	Domain  string
	QType   uint16
	DataLen uint16
	TTL     uint32
}

func (u *UnknownRecord) getType() int {
	return Unknown
}

func (u *UnknownRecord) Write(buffer *BytePacketBuffer) (int, error) {
	fmt.Println("Skipping record")
	return 0, errors.New("Unsupported DNS record type")
}

type ARecord struct {
	Domain string
	Addr   net.IP
	TTL    uint32
}

func (a *ARecord) getType() int {
	return A
}

func (a *ARecord) Write(buffer *BytePacketBuffer) (int, error) {
		startPos := buffer.pos

		if err := buffer.WriteQName(&a.Domain); err != nil {
			return 0, err
		}
		if err := buffer.WriteU16(A); err != nil {
			return 0, err
		}
		if err := buffer.WriteU16(1); err != nil {
			return 0, err
		}
		if err := buffer.WriteU32(a.TTL); err != nil {
			return 0, err
		}
		if err := buffer.WriteU16(4); err != nil {
			return 0, err
		}

		octets := a.Addr.To4()
		if octets == nil {
			return 0, errors.New("Invalid IPv4 address")
		}

		if err := buffer.Write(octets[0]); err != nil {
			return 0, err
		}
		if err := buffer.Write(octets[1]); err != nil {
			return 0, err
		}
		if err := buffer.Write(octets[2]); err != nil {
			return 0, err
		}
		if err := buffer.Write(octets[3]); err != nil {
			return 0, err
		}

	return int(buffer.pos - startPos), nil
}

type NSRecord struct {
	Domain string
	Host   string
	TTL    uint32
}

func (ns *NSRecord) getType() int {
	return NS
}

func (ns *NSRecord) Write(buffer *BytePacketBuffer) (int, error) {
	startPos := buffer.pos

	if err := buffer.WriteQName(&ns.Domain); err != nil {
		return 0, err
	}
	if err := buffer.WriteU16(NS); err != nil {
		return 0, err
	}
	if err := buffer.WriteU16(1); err != nil {
		return 0, err
	}
	if err := buffer.WriteU32(ns.TTL); err != nil {
		return 0, err
	}

	pos := buffer.pos
	if err := buffer.WriteU16(0); err != nil {
		return 0, err
	}

	buffer.WriteQName(&ns.Host)

	size := buffer.pos - (pos + 2)
	buffer.SetU16(uint8(pos), uint16(size))

	return int(buffer.pos - startPos), nil
}


type CNAMERecord struct {
	Domain string
	Host   string
	TTL    uint32
}

func (cname *CNAMERecord) getType() int {
	return CNAME
}

func (cname *CNAMERecord) Write(buffer *BytePacketBuffer) (int, error) {
	startPos := buffer.pos

	if err := buffer.WriteQName(&cname.Domain); err != nil {
		return 0, err
	}
	if err := buffer.WriteU16(CNAME); err != nil {
		return 0, err
	}
	if err := buffer.WriteU16(1); err != nil {
		return 0, err
	}
	if err := buffer.WriteU32(cname.TTL); err != nil {
		return 0, err
	}

	pos := buffer.pos
	if err := buffer.WriteU16(0); err != nil {
		return 0, err
	}

	buffer.WriteQName(&cname.Host)

	size := buffer.pos - (pos + 2)
	buffer.SetU16(uint8(pos), uint16(size))

	return int(buffer.pos - startPos), nil
}


type MXRecord struct {
	Domain string
	Priority   uint16
	Host	string
	TTL    uint32
}

func (mx *MXRecord) getType() int {
	return MX
}

func (mx *MXRecord) Write(buffer *BytePacketBuffer) (int, error) {
	startPos := buffer.pos

	if err := buffer.WriteQName(&mx.Domain); err != nil {
		return 0, err
	}
	if err := buffer.WriteU16(MX); err != nil {
		return 0, err
	}
	if err := buffer.WriteU16(1); err != nil {
		return 0, err
	}
	if err := buffer.WriteU32(mx.TTL); err != nil {
		return 0, err
	}

	pos := buffer.pos
	if err := buffer.WriteU16(0); err != nil {
		return 0, err
	}

	buffer.WriteU16(mx.Priority)
	buffer.WriteQName(&mx.Host)

	size := buffer.pos - (pos + 2)
	buffer.SetU16(uint8(pos), uint16(size))

	return int(buffer.pos - startPos), nil
}

type AAAARecord struct {
	Domain string
	Addr   net.IP
	TTL    uint32
}

func (a4 *AAAARecord) getType() int {
	return AAAA
}

func (a4 *AAAARecord) Write(buffer *BytePacketBuffer) (int, error) {
	startPos := buffer.pos

	if err := buffer.WriteQName(&a4.Domain); err != nil {
		return 0, err
	}
	if err := buffer.WriteU16(AAAA); err != nil {
		return 0, err
	}
	if err := buffer.WriteU16(1); err != nil {
		return 0, err
	}
	if err := buffer.WriteU32(a4.TTL); err != nil {
		return 0, err
	}
	if err := buffer.WriteU16(16); err != nil {
		return 0, err
	}

	for _, octet := range a4.Addr {
		buffer.WriteU16(uint16(octet))
	}

	return int(buffer.pos - startPos), nil
}

func ReadDnsRecord(buffer *BytePacketBuffer) (DnsRecord, error) {
	var domain string
	if err := buffer.ReadQName(&domain); err != nil {
		return nil, err
	}

	qTypeNum, err := buffer.ReadU16()
	if err != nil {
		return nil, err
	}
	qType := QueryTypeFromNum(qTypeNum)

	if _, err := buffer.ReadU16(); err != nil {
		return nil, err
	}

	ttl, err := buffer.ReadU32()
	if err != nil {
		return nil, err
	}

	dataLen, err := buffer.ReadU16()
	if err != nil {
		return nil, err
	}

	switch qType.query_type {
	case A:
		rawAddr, err := buffer.ReadU32()
		if err != nil {
			return nil, err
		}
		addr := net.IPv4(
			byte(rawAddr>>24),
			byte(rawAddr>>16),
			byte(rawAddr>>8),
			byte(rawAddr),
		)
		return &ARecord{
			Domain: domain,
			Addr:   addr,
			TTL:    ttl,
		}, nil
	case AAAA:
		rawAddr1, err := buffer.ReadU32()
		if err != nil {
			return nil, err
		}
		rawAddr2, err := buffer.ReadU32()
		if err != nil {
			return nil, err
		}
		rawAddr3, err := buffer.ReadU32()
		if err != nil {
			return nil, err
		}
		rawAddr4, err := buffer.ReadU32()
		if err != nil {
			return nil, err
		}
		addr := net.IP{
			byte(rawAddr1 >> 16),
			byte(rawAddr1),
			byte(rawAddr2 >> 16),
			byte(rawAddr2),
			byte(rawAddr3 >> 16),
			byte(rawAddr3),
			byte(rawAddr4 >> 16),
			byte(rawAddr4),
		}
		return &AAAARecord{
			Domain: domain,
			Addr:   addr,
			TTL:    ttl,
		}, nil
	case NS:
		var ns string
		buffer.ReadQName(&ns)

		return &NSRecord{
			Domain: domain,
			Host: ns,
			TTL: ttl,
		}, nil
	case CNAME:
		var cname string
		buffer.ReadQName(&cname)

		return &CNAMERecord{
			Domain: domain,
			Host: cname,
			TTL: ttl,
		}, nil
	case MX:
		priority, err := buffer.ReadU16()
		if err != nil {
			return nil, err
		}
		var mx string
		buffer.ReadQName(&mx)

		return &MXRecord{
			Domain: domain,
			Priority: priority,
			Host: mx,
			TTL: ttl,
		}, nil
	default:
		if err := buffer.Step(uint16(dataLen)); err != nil {
			return nil, err
		}
		return &UnknownRecord{
			Domain:  domain,
			QType:   qTypeNum,
			DataLen: dataLen,
			TTL:     ttl,
		}, nil
	}
}