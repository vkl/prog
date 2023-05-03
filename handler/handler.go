package handler

import (
	"fmt"
	"log"
	"strings"
	"time"

	"go.bug.st/serial"
)

const (
	PAGE_SZ = 64
)

type UartHandler struct {
	mode   *serial.Mode
	serial serial.Port
	port   string
}

func (c *UartHandler) SendMsg(msg []byte) error {
	if c.serial == nil {
		if err := c.Init(); err != nil {
			return err
		}
	}
	fmt.Printf("Sending %s\n", string(msg))
	n, err := c.serial.Write(msg)
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("prog didn't response or serial port error")
	}
	buff := make([]byte, 16)
	resp := ""
	for {
		n, err := c.serial.Read(buff)
		if err != nil {
			return err
		}
		if n == 0 {
			break
		}
		resp += string(buff[:n])
		if strings.HasSuffix(resp, "\r") {
			break
		}
	}
	resp = strings.Trim(resp, "\r")
	fmt.Println("Receive " + resp)
	if !strings.HasSuffix(resp, "OK") {
		return fmt.Errorf("response is invalid" + resp)
	}
	return nil
}

func (c *UartHandler) SendPage(page []byte) {
	n, err := c.serial.Write(page)
	if err != nil {
		log.Fatal(err)
	}
	if n == 0 {
		log.Fatal("prog didn't response or serial port error")
	}
	buff := make([]byte, 10)
	resp := ""
	for {
		n, err := c.serial.Read(buff)
		if err != nil {
			log.Fatal(err)
			break
		}
		if n == 0 {
			break
		}
		resp += string(buff[:n])
		if strings.HasSuffix(resp, "\r") {
			break
		}
	}
}

func (c *UartHandler) SendData(data []byte) {
	fmt.Println("Send data")
	pages := len(data) / PAGE_SZ
	progress := InitProgress()
	for page := 0; page < pages; page++ {
		c.SendPage(data[(page * PAGE_SZ):(PAGE_SZ * (1 + page))])
		progress.Show(page, pages-1)
	}
}

func (c *UartHandler) ReceiveData(size int) []byte {
	progress := InitProgress()
	data := make([]byte, 0)
	n, err := c.serial.Write([]byte("\r"))
	if err != nil {
		log.Fatal(err)
	}
	if n == 0 {
		log.Fatal("prog didn't response or serial port error")
	}
	fmt.Println("Receive data")
	buff := make([]byte, 128)
	ns := 0
	for {
		n, err := c.serial.Read(buff)
		if n == 0 {
			break
		}
		ns += n
		if err != nil {
			log.Fatal(err)
			break
		}

		progress.Show(ns, size)

		if strings.HasSuffix(string(buff[:n]), "\r") {
			data = append(data, buff[:(n-1)]...)
			break
		} else {
			data = append(data, buff[:n]...)
		}
	}
	return data
}

func (c *UartHandler) Init() error {
	serial, err := serial.Open(c.port, c.mode)
	if err != nil {
		return fmt.Errorf("open port %s: error: %v", c.port, err)
	}
	serial.SetReadTimeout(3 * time.Second)
	c.serial = serial
	return nil
}

func InitUartHandler(port string) UartHandler {
	mode := &serial.Mode{
		BaudRate: 57600,
	}

	cmdHandler := UartHandler{
		mode: mode,
		port: port,
	}
	return cmdHandler
}
