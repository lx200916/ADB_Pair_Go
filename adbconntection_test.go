package main

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestAdbConnection_Init(t *testing.T) {
	var c AdbConnection
	c.Init("tcp:127.0.0.1:5037")
}
func TestAdbConnection_Pair(t *testing.T) {
	var c AdbConnection
	c.Init("tcp:127.0.0.1:5037")
	fmt.Println("Connected")
	c.writeString("host:pair:123456:192.168.0.0:1111") //Dummy data
	fmt.Println("Sent")

	status, err := c.readStatus()

	status = strings.ToLower(status)

	if err != nil {
		t.Errorf("Error: %s", err)
		t.Fail()
	}
	if !regexp.MustCompile(`^okay`).MatchString(status) {
		t.Errorf("Expected OKAY, got %s", status)
		t.Fail()
	}

}
