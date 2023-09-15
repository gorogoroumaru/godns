package main

import (
	"errors"
)

type BytePacketBuffer struct {
    buf [512]byte
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

func (b *BytePacketBuffer) Read() (byte, error) {
    if b.pos >= 512 {
        return 0, errors.New("End of buffer")
    }
    res := b.buf[b.pos]
    b.pos++
    return res, nil
}

func (b *BytePacketBuffer) Get(pos uint16) (byte, error) {
    if pos >= 512 {
        return 0, errors.New("End of buffer")
    }
    return b.buf[pos], nil
}

func (b *BytePacketBuffer) GetRange(start, length uint16) ([]byte, error) {
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
