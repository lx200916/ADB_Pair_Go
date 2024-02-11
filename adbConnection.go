package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	AdbStatusOkay = "OKAY"
	AdbStatusFail = "FAIL"
)

type AdbConnection struct {
	Addr string
	conn *net.TCPConn
}

func (c *AdbConnection) Init(addr string) {
	//if c.conn != nil {
	//	c.conn.Close()
	//}
	if addr == "" {
		log.Fatal("adb: no address specified")
	}
	if strings.HasPrefix(addr, "tcp:") {
		addr = addr[4:]
		addr = strings.Trim(addr, " ")
		resolvedAddr, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			log.Fatal("adb: failed to resolve ", addr, ": ", err)
		}
		c.Addr = addr
		conn, err := net.DialTCP("tcp", nil, resolvedAddr)
		if err != nil {
			log.Fatal("adb: failed to connect to ", addr, ": ", err)
		}
		c.conn = conn
	} else {
		log.Fatal("adb: unsupported address ", addr)
	}

}
func (c *AdbConnection) Close() error {
	return (*c.conn).Close()
}
func (c *AdbConnection) writeString(s string) {
	if c.conn == nil {
		log.Fatal("adb: connection is closed")
	}
	c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	byteContent := fmt.Sprintf("%04x%s", len(s), s)
	_, err := c.conn.Write([]byte(byteContent))
	if err != nil {
		log.Fatal("adb: failed to write to ", c.Addr, ": ", err)
	}
}
func (c *AdbConnection) readExactly(n int) []byte {
	if c.conn == nil {
		log.Fatal("adb: connection is closed")
	}
	c.conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	buf := make([]byte, n)
	_, err := c.conn.Read(buf)
	if err != nil {
		log.Fatal("adb: failed to read from ", c.Addr, ": ",
			err)
	}
	return buf
}
func (c *AdbConnection) readString() string {
	if c.conn == nil {
		log.Fatal("adb: connection is closed")
	}
	c.conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	buf := make([]byte, 4)
	_, err := c.conn.Read(buf)
	if err != nil {
		log.Fatal("adb: failed to read from ", c.Addr, ": ", err)
	}
	length, err := strconv.ParseInt(string(buf), 16, 32)
	if err != nil {
		log.Fatal("adb: failed to convert length from ", c.Addr, ": ", err)
	}
	buf = make([]byte, length)
	_, err = c.conn.Read(buf)
	if err != nil {
		log.Fatal("adb: failed to read from ", c.Addr, ": ", err)
	}
	return string(buf)
}
func (c *AdbConnection) readStatus() (string, error) {
	status := string(c.readExactly(4))
	if status == AdbStatusFail {
		message := "adb protocol failed: " + c.readString()
		return AdbStatusFail, errors.New(message)
	}
	if status != AdbStatusOkay {
		log.Fatal("adb: unexpected status ", status)
	}
	return AdbStatusOkay, nil
}
