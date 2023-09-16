package main

import (
	"errors"
	"strings"
)

type BytePacketBuffer struct {
    buf [512]uint8
    pos uint16
}

func NewBytePacketBuffer() *BytePacketBuffer {
    return &BytePacketBuffer{}
}

func (b *BytePacketBuffer) Pos() uint16 {
    return b.pos
}

func (b *BytePacketBuffer) Step(steps uint16) error {
    b.pos += steps
    return nil
}

func (b *BytePacketBuffer) Seek(pos uint16) error {
    b.pos = pos
    return nil
}

func (b *BytePacketBuffer) Read() (uint8, error) {
    if b.pos >= 512 {
        return 0, errors.New("End of buffer")
    }
    res := b.buf[b.pos]
    b.pos++
    return res, nil
}

func (b *BytePacketBuffer) Write(val uint8) (error) {
    if b.pos >= 512 {
        return errors.New("End of buffer")
    }
    b.buf[b.pos] = val
    b.pos++
    return nil
}

func (b *BytePacketBuffer) Set(pos uint8, val uint8) (error) {
    b.buf[pos] = val
    return nil
}


func (b *BytePacketBuffer) SetU16(pos uint8, val uint16) (error) {
    b.Set(pos, uint8(val >> 8))
    b.Set(pos + 1, uint8(val & 0xff))
    return nil
}

func (b *BytePacketBuffer) Get(pos uint16) (uint8, error) {
    if pos >= 512 {
        return 0, errors.New("End of buffer")
    }
    return b.buf[pos], nil
}

func (b *BytePacketBuffer) GetRange(start, length uint16) ([]uint8, error) {
    if start+length >= 512 {
        return nil, errors.New("End of buffer")
    }
    return b.buf[start : start+length], nil
}

func (b *BytePacketBuffer) ReadU16() (uint16, error) {
    high, err := b.Read()
    if err != nil {
        return 0, err
    }
    low, err := b.Read()
    if err != nil {
        return 0, err
    }
    return uint16(high)<<8 | uint16(low), nil
}

func (b *BytePacketBuffer) WriteU16(val uint16) (error) {
    if b.pos >= 512 {
        return errors.New("End of buffer")
    }
    b.Write(uint8(val >> 8))
    b.Write(uint8(val & 0xff))
    return nil
}

func (b *BytePacketBuffer) ReadU32() (uint32, error) {
    b1, err := b.Read()
    if err != nil {
        return 0, err
    }
    b2, err := b.Read()
    if err != nil {
        return 0, err
    }
    b3, err := b.Read()
    if err != nil {
        return 0, err
    }
    b4, err := b.Read()
    if err != nil {
        return 0, err
    }
    return (uint32(b1) << 24) | (uint32(b2) << 16) | (uint32(b3) << 8) | uint32(b4), nil
}

func (b *BytePacketBuffer) WriteU32(val uint32) (error) {
    if b.pos >= 512 {
        return errors.New("End of buffer")
    }
    b.Write(uint8((val >> 24) & 0xff))
    b.Write(uint8((val >> 16) & 0xff))
    b.Write(uint8((val >> 8) & 0xff))
    b.Write(uint8(val & 0xff))
    return nil
}

func (b *BytePacketBuffer) ReadQName(outstr *string) error {
    pos := b.Pos()
    jumped := false
    maxJumps := 5
    jumpsPerformed := 0
    delim := ""

    for {
        if jumpsPerformed > maxJumps {
            return errors.New("Limit of 5 jumps exceeded")
        }

        lenByte, err := b.Get(pos)
        if err != nil {
            return err
        }

        if lenByte&0xC0 == 0xC0 {
            if !jumped {
                b.Seek(pos + 2)
            }

            b2, err := b.Get(pos + 1)
            if err != nil {
                return err
            }
			// 0xC0がフラグになっていてそれ以外のビットがoffsetになっている
            offset := ((uint16(lenByte) ^ 0xC0) << 8) | uint16(b2)
            pos = uint16(offset)

            jumped = true
            jumpsPerformed++
            continue
        } else {
            pos++
            if lenByte == 0 {
                break
            }

            *outstr += delim

            strBuffer, err := b.GetRange(pos, uint16(lenByte))
            if err != nil {
                return err
            }
            *outstr += string(strBuffer)
            delim = "."
            pos += uint16(lenByte)
        }
    }

    if !jumped {
        b.Seek(pos)
    }

    return nil
}

func (buffer *BytePacketBuffer) WriteQName(qname *string) error {
	labels := strings.Split(*qname, ".")

	for _, label := range labels {
		length := len(label)
		if length > 0x3F {
			return errors.New("Single label exceeds 63 characters of length")
		}

		if err := buffer.Write(byte(length)); err != nil {
			return err
		}

		for _, b := range []byte(label) {
			if err := buffer.Write(b); err != nil {
				return err
			}
		}
	}

	if err := buffer.Write(0); err != nil {
		return err
	}

	return nil
}
