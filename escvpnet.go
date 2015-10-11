package escvpnetgo

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

const HeaderIdentifierString string = "ESC/VP.net"

const (
	TYPE_NULL uint8 = iota
	TYPE_HELLO
	TYPE_PASSWORD
	TYPE_CONNECT
)

const (
	STATUS_OK              uint8 = 32
	STATUS_BAD_REQUEST           = 64
	STATUS_UNAUTHORIZED          = 65
	STATUS_FORBIDDEN             = 67
	STATUS_NOT_ALLOWED           = 69
	STATUS_UNAVAILABLE           = 83
	STATUS_INVALID_VERSION       = 85
)

const (
	PACKET_ID_NULL uint8 = iota
	PACKET_ID_PASSWORD
	PACKET_ID_NEW_PASSWORD
	PACKET_ID_PROJECTOR_NAME
	PACKET_ID_IM
	PACKET_ID_COMMAND
)

var HeaderIdentifier = ([]byte)(HeaderIdentifierString)

type Header struct {
	Identifier []byte
	Version    uint8
	Type       uint8
	Reserved   uint16
	Status     uint8
	PacketsNum uint8
}

type HeaderError struct {
	status uint8
}

type Packet struct {
	Identifier uint8
	Attribute  uint8
	Str        []byte
}

type ESCPVPNET struct {
	Target string
	conn   net.Conn
	reader *bufio.Reader
}

func (he *HeaderError) Error() string {
	return fmt.Sprintf("INVALID CONNECTION STATUS: %d", he.status)
}

func NewESCVPNET(target string) (*ESCPVPNET, error) {

	tcpAddr, err := net.ResolveTCPAddr("tcp", target)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	client := &ESCPVPNET{target, conn, nil}

	err = client.connect()

	if err != nil {
		return nil, err
	}

	client.reader = bufio.NewReader(conn)

	return client, nil
}

func (c *ESCPVPNET) connect() error {
	_, err := c.conn.Write(NewHeader(TYPE_CONNECT, 0))

	if err != nil {
		return err
	}

	headerBytes := make([]byte, 16)

	_, err = io.ReadFull(c.conn, headerBytes)

	if err != nil {
		return err
	}

	header, err := BytesToHeader(headerBytes)

	if err != nil {
		return err
	}

	if header.Status != STATUS_OK {
		return &HeaderError{header.Status}
	}

	return nil
}

func (c *ESCPVPNET) Close() error {
	return c.conn.Close()
}

func (c *ESCPVPNET) Execute(command string) (string, error) {
	_, err := c.conn.Write(([]byte)(command + "\r"))

	if err != nil {
		return "", err
	}

	str, err := c.reader.ReadString('\r')

	if err != nil {
		return "", err
	}

	str = strings.TrimLeft(strings.TrimSpace(str), ":")

	if str == "ERR" {
		return "", errors.New("ERROR EXECUTING COMMAND")
	}

	return str, nil
}

func NewHeader(headerType uint8, count uint8) []byte {
	header := &Header{HeaderIdentifier, 16, headerType, 0, 0, count}
	return HeaderToBytes(header)
}

func NewPacket(identifier uint8, attribute uint8, command string) []byte {
	packet := &Packet{identifier, attribute, ([]byte)(command)}
	return PacketToBytes(packet)
}

func HeaderToBytes(header *Header) []byte {
	headerBytes := make([]byte, 16)

	for i := 0; i < 10; i++ {
		headerBytes[i] = header.Identifier[i]
	}

	headerBytes[10] = header.Version
	headerBytes[11] = header.Type

	// Reserved
	headerBytes[12] = 0
	headerBytes[13] = 0

	headerBytes[14] = header.Status
	headerBytes[15] = header.PacketsNum

	return headerBytes
}

func PacketToBytes(packet *Packet) []byte {
	packetBytes := make([]byte, 18)
	packetBytes[0] = packet.Identifier
	packetBytes[1] = packet.Attribute

	for i := 0; i < 16 && i < len(packet.Str); i++ {
		packetBytes[i+2] = packet.Str[i]
	}

	return packetBytes
}

func BytesToHeader(headerBytes []byte) (*Header, error) {
	if len(headerBytes) != 16 {
		return nil, errors.New("Invalid length")
	}

	header := new(Header)
	header.Identifier = headerBytes[0:10]

	if bytes.Compare(header.Identifier, HeaderIdentifier) != 0 {
		return nil, errors.New("Invalid Identifier " + (string)(header.Identifier))
	}

	header.Version = headerBytes[10]
	header.Type = headerBytes[11]
	header.Status = headerBytes[14]
	header.PacketsNum = headerBytes[15]

	return header, nil
}

func BytesToPacket(packetBytes []byte) (*Packet, error) {
	if len(packetBytes) != 18 {
		return nil, errors.New("Invalid length")
	}

	packet := new(Packet)

	packet.Identifier = packetBytes[0]
	packet.Attribute = packetBytes[1]
	packet.Str = packetBytes[2:18]

	return packet, nil
}
