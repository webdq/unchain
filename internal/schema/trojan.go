package schema

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

// ProtoTrojan is the structure of trojan protocol
// https://github.com/trojan-gfw/trojan/blob/master/docs/protocol.md
type ProtoTrojan struct {
	sha224password []byte //trojan password
	DstProtocol    string //tcp or udp
	dstHost        string
	dstHostType    string //ipv6 or ipv4,domain
	dstPort        uint16
	Version        byte
	payload        []byte
}

const (
	byteCR = '\r'
	byteLF = '\n'
)

func parseTrojanHeader(buffer []byte) (*ProtoTrojan, error) {
	if len(buffer) < 56 {
		return nil, errors.New("invalid data")
	}
	bytes.Split(buffer, []byte{byteCR, byteLF})

	crLfIndex := 56
	if buffer[56] != byteCR || buffer[57] != byteLF {
		return nil, errors.New("invalid header format (missing CR LF)")
	}
	p := &ProtoTrojan{
		sha224password: buffer[:crLfIndex],
	}

	socks5DataBuffer := buffer[crLfIndex+2:]
	if len(socks5DataBuffer) < 6 {
		return nil, errors.New("invalid SOCKS5 request data")
	}

	cmd := socks5DataBuffer[0]
	if cmd == 0x01 { //connect
		p.DstProtocol = "tcp"
	} else if cmd == 0x03 { //udp
		p.DstProtocol = "udp"
		//todo:: udp
	} else {
		return nil, errors.New("unsupported command, only TCP (CONNECT) is allowed")
	}

	atype := socks5DataBuffer[1]
	var addressLength int
	addressIndex := 2
	switch atype {
	case 1:
		addressLength = 4
		ip := net.IP(socks5DataBuffer[addressIndex : addressIndex+addressLength])
		p.dstHost = ip.String()
		p.dstHostType = "ipv4"
		p.dstPort = binary.BigEndian.Uint16(socks5DataBuffer[addressIndex+addressLength : addressIndex+addressLength+2])
		p.payload = socks5DataBuffer[addressIndex+addressLength+4:]
	case 3: //domain
		addressLength = int(socks5DataBuffer[addressIndex])
		addressIndex++
		p.dstHostType = "domain"
		p.dstHost = string(socks5DataBuffer[addressIndex : addressIndex+addressLength])
		p.dstPort = binary.BigEndian.Uint16(socks5DataBuffer[addressIndex+addressLength : addressIndex+addressLength+2])
		p.payload = socks5DataBuffer[addressIndex+addressLength+4:]
	case 4:
		addressLength = 16
		ip := net.IP(socks5DataBuffer[addressIndex : addressIndex+addressLength])
		p.dstPort = binary.BigEndian.Uint16(socks5DataBuffer[addressIndex+addressLength : addressIndex+addressLength+2])
		p.dstHost = ip.String()
		p.dstHostType = "ipv6"
		p.payload = socks5DataBuffer[addressIndex+addressLength+4:]
	default:
		return nil, fmt.Errorf("invalid addressType is %d", atype)
	}
	return p, nil
}
