package common

import (
	"encoding/json"
	"fmt"
	"net"
)

type Result struct {
	Value int
	Message string
}

const (
	MSG_SUCCESS = 200
	MSG_FAILURE
)

func SendResponse(value int, message string, conn net.Conn) error {
	e := json.NewEncoder(conn)

	r := &Result{Value: value, Message: message}
	if err := e.Encode(r); err != nil {
		return err
	}

	return nil
}

func RecvResponse(conn net.Conn) (*Result, error) {
	d := json.NewDecoder(conn)

	var r Result
	if err := d.Decode(&r); err != nil {
		fmt.Println("RecvResponse failed")
		return nil, err
	}

	return &r, nil
}
