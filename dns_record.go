package main

import (
	"errors"
	"fmt"
	"net"
)

type DnsRecord interface {
	isDnsRecord()
	Write(*BytePacketBuffer)(int, error)
}

type UnknownRecord struct {
	Domain  string
	QType   uint16
	DataLen uint16
	TTL     uint32
}

func (u *UnknownRecord) isDnsRecord() {}

func (u *UnknownRecord) Write(buffer *BytePacketBuffer) (int, error) {
	fmt.Println("Skipping record")
	return 0, errors.New("Unsupported DNS record type")
}

type ARecord struct {
	Domain string
	Addr   net.IP
	TTL    uint32
}

func (a *ARecord) isDnsRecord() {}

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