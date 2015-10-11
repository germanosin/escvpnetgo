package escvpnetgo

import (
	"fmt"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Log("Running Create Test\n")
	c, err := NewESCVPNET("10.150.136.32:3629")

	if err != nil {
		t.Fatalf("Unable to connect server %s\n", err.Error())
	}

	if c == nil {
		t.Fatalf("Unable to connect server client is nil\n")
	}

}

func TestClose(t *testing.T) {
	t.Log("Running Close Test\n")
	c, err := NewESCVPNET("10.150.136.32:3629")

	if err != nil {
		t.Fatalf("Unable to connect server %s\n", err.Error())
	}

	if c == nil {
		t.Fatalf("Unable to connect server client is nil\n")
	}

	err = c.Close()

	if err != nil {
		t.Fatalf("Unable to close connection %s\n", err.Error())
	}
}

func TestExecute(t *testing.T) {

	t.Log("Running Close Test\n")

	c, err := NewESCVPNET("10.150.136.32:3629")

	if err != nil {
		t.Fatalf("Unable to connect server %s\n", err.Error())
	}

	if c == nil {
		t.Fatalf("Unable to connect server client is nil\n")
	}

	result, err := c.Execute("LAMP?")

	if err != nil {
		t.Fatalf("Unable to execute command %s\n", err.Error())
	}

	fmt.Printf("Execute result: %s\n", result)
}
