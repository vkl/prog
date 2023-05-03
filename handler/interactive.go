package handler

import (
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"go.bug.st/serial"
)

func usage(w io.Writer) {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, completer.Tree("    "))
}

// Function constructor - constructs new function for listing given directory
func listFiles(path string) func(string) []string {
	return func(line string) []string {
		names := make([]string, 0)
		files, _ := ioutil.ReadDir(path)
		for _, f := range files {
			names = append(names, f.Name())
		}
		return names
	}
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem("read",
		readline.PcItem("addr"),
		readline.PcItemDynamic(listFiles("./")),
	),
	readline.PcItem("write",
		readline.PcItemDynamic(listFiles("./")),
		readline.PcItem("addr"),
	),
	readline.PcItem("check",
		readline.PcItemDynamic(listFiles("./")),
	),
	readline.PcItem("show",
		readline.PcItemDynamic(listFiles("./")),
	),
	readline.PcItem("ping"),
	readline.PcItem("serial"),
	readline.PcItem("set serial"),
	readline.PcItem("size"),
	readline.PcItem("set size",
		readline.PcItem("value")),
	readline.PcItem("quit"),
	readline.PcItem("help"),
)

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func Interactive(eeprom EEPROMHandler) {
	l, err := readline.NewEx(&readline.Config{
		Prompt:          "Prog> ",
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()
	//l.CaptureExitSignal()

	//log.SetOutput(l.Stderr())
	for {

		line, err := l.Readline()

		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		switch {
		case line == "debug":
			currfiles := listFiles(".")
			for _, file := range currfiles("") {
				fmt.Println(file)
			}
		case line == "help":
			usage(l.Stderr())
		case line == "quit":
			goto exit
		case line == "ping":
			eeprom.Ping()
		case line == "serial":
			fmt.Println(eeprom.Uart.port)
		case line == "set serial":
			ports, err := serial.GetPortsList()
			if err != nil {
				fmt.Printf("error: %v\n", err)
				break
			}
			if len(ports) == 0 {
				fmt.Println("No serial ports found!")
			}
			fmt.Println("Ports:")
			for i, port := range ports {
				fmt.Printf("\t%d %v\n", i, port)
			}
			l.SetPrompt(" ?> ")
			for {
				s, err := l.Readline()
				if err != nil {
					fmt.Println(err)
					continue
				}
				i, err := strconv.ParseInt(s, 0, 8)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if i < 0 || int(i) > len(ports) {
					fmt.Println("wrong choice")
					continue
				}
				eeprom.Uart.port = ports[i]
				l.SetPrompt("Prog> ")
				break
			}

		case strings.HasPrefix(line, "read addr"):
			line := strings.TrimSpace(line[9:])
			if len(line) == 0 {
				fmt.Println("addr?")
				break
			}
			addr, err := strconv.ParseInt(line, 0, 64)
			if err != nil {
				fmt.Printf("Addr %s is wrong\n", line)
				break
			}
			eeprom.ReadAddr(addr)
		case strings.HasPrefix(line, "write addr"):
			data := strings.Split(
				strings.TrimSpace(line[10:]), " ")
			if len(data) < 2 {
				fmt.Println("addr or/and byte?")
				break
			}
			addr, err := strconv.ParseInt(data[0], 0, 64)
			if err != nil {
				fmt.Printf("Addr %s is wrong\n", data[0])
				break
			}
			b, err := strconv.ParseInt(data[1], 0, 64)
			if err != nil {
				fmt.Printf("Byte %s is wrong\n", data[1])
				break
			}
			if b > 255 || b < 0 {
				fmt.Printf("Byte %d must be between 0-255\n", b)
				break
			}
			eeprom.WriteAddr(addr, b)
		case strings.HasPrefix(line, "write"):
			fileName := strings.TrimSpace(line[5:])
			if len(fileName) == 0 {
				fmt.Println("filename?")
				break
			}
			eeprom.Write(fileName)
		case strings.HasPrefix(line, "read"):
			if len(line) <= 4 {
				fmt.Println("filename?")
				break
			}
			fileName := strings.TrimSpace(line[4:])
			if len(fileName) == 0 {
				fmt.Println("filename?")
				break
			}
			eeprom.Read(fileName)
		case strings.HasPrefix(line, "check"):
			if len(line) <= 5 {
				fmt.Println("filename?")
				break
			}
			fileName := strings.TrimSpace(line[5:])
			if len(fileName) == 0 {
				fmt.Println("filename?")
				break
			}
			eeprom.Check(fileName)
		case line == "size":
			fmt.Println(eeprom.Size)
		case strings.HasPrefix(line, "set size"):
			sz_s := []int{S1k, S8k, S16k, S32k, S64k}
			for i, sz := range sz_s {
				fmt.Printf("\t%d %v\n", i, sz)
			}
			l.SetPrompt(" ?> ")
			for {
				s, err := l.Readline()
				if err != nil {
					fmt.Println(err)
					continue
				}
				i, err := strconv.ParseInt(s, 0, 8)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if i < 0 || int(i) > len(sz_s) {
					fmt.Println("wrong choice")
					continue
				}
				eeprom.Size = sz_s[i]
				l.SetPrompt("Prog> ")
				break
			}

		case strings.HasPrefix(line, "show"):
			if len(line) <= 4 {
				fmt.Println("filename?")
				break
			}
			fileName := strings.TrimSpace(line[4:])
			if len(fileName) == 0 {
				fmt.Println("filename?")
				break
			}
			f, err := os.Open(fileName)
			defer f.Close()
			if err != nil {
				fmt.Println(err)
				f.Close()
				break
			}
			fstat, err := f.Stat()
			if err != nil {
				fmt.Println(err)
				f.Close()
				break
			}
			data := make([]byte, fstat.Size())
			for {
				n, err := f.Read(data)
				if err != nil {
					fmt.Println(err)
					break
				}
				if n == 0 {
					fmt.Println("0")
					break
				}
			}
			hexdumper := hex.Dumper(os.Stdout)
			_, err = hexdumper.Write(data)
			if err != nil {
				fmt.Println(err)
			}
			hexdumper.Close()
			f.Close()
		default:
			//log.Println("you said:", strconv.Quote(line))
		}
	}
exit:
}
