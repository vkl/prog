package handler

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	S1k  = 64
	S8k  = 8192
	S16k = 16384
	S32k = 32768
	S64k = 65536
)

type EEPROMHandler struct {
	Uart UartHandler
	Size int
}

func (eeprom *EEPROMHandler) ReadData() []byte {
	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, "read 0x%04X\r", eeprom.Size)
	if err := eeprom.Uart.SendMsg(buf.Bytes()); err != nil {
		fmt.Printf("error: %v\n", err)
		return nil
	}
	return eeprom.Uart.ReceiveData(eeprom.Size)
}

func (eeprom *EEPROMHandler) Read(fileName string) error {
	data := eeprom.ReadData()
	if len(data) == 0 {
		return fmt.Errorf("No data")
	}
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	defer f.Close()
	if err != nil {
		return err
	}
	fmt.Printf("Write to file '%s' size %d\n", fileName, len(data))
	n, err := f.Write(data)
	if err != nil {
		return err
	}
	fmt.Printf("Written %d\n", n)
	return nil
}

func (eeprom *EEPROMHandler) Write(fileName string) error {
	var data []byte
	f, err := os.Open(fileName)
	defer f.Close()
	if err != nil {
		return err
	}
	data, err = ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	fmt.Printf("Write from file '%s' size %d\n", fileName, len(data))
	if err := eeprom.Uart.SendMsg([]byte("write\r")); err != nil {
		return fmt.Errorf("error: %v\n", err)
	}
	size := len(data)
	pages := size / PAGE_SZ
	buf := bytes.Buffer{}
	fmt.Printf("Write %d pages by %db\n", pages, PAGE_SZ)
	fmt.Fprintf(&buf, "0x%04X\r", pages)
	if err := eeprom.Uart.SendMsg(buf.Bytes()); err != nil {
		return fmt.Errorf("error: %v\n", err)
	}
	eeprom.Uart.SendData(data)
	return nil
}

func (eeprom *EEPROMHandler) Check(fileName string) error {
	f, err := os.Open(fileName)
	defer f.Close()
	if err != nil {
		return err
	}
	data := make([]byte, eeprom.Size)
	_, err = f.Read(data)
	if err != nil {
		return nil
	}
	md5_file := fmt.Sprintf("%x", md5.Sum(data))
	md5_eeprom := fmt.Sprintf("%x", md5.Sum(eeprom.ReadData()))
	fmt.Printf("file %s\n", md5_file)
	fmt.Printf("eeprom %s\n", md5_eeprom)
	if strings.Compare(md5_file, md5_eeprom) == 0 {
		fmt.Println("Verification OK")
	} else {
		fmt.Println("Verification BAD")
	}
	return nil
}

func (eeprom *EEPROMHandler) ReadAddr(addr int64) {
	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, "read addr 0x%04X\r", addr)
	if err := eeprom.Uart.SendMsg(buf.Bytes()); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

func (eeprom *EEPROMHandler) WriteAddr(addr int64, b int64) {
	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, "write addr 0x%04X 0x%02X\r", addr, b)
	if err := eeprom.Uart.SendMsg(buf.Bytes()); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

func (eeprom *EEPROMHandler) Ping() {
	if err := eeprom.Uart.SendMsg([]byte("ping\r")); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
