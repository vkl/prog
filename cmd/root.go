package cmd

import (
	"fmt"
	"log"
	"vkl/prog/handler"

	"github.com/spf13/cobra"
	"go.bug.st/serial"
)

var (
	Port     string
	FileName string
	uart     handler.UartHandler

	rootCmd = &cobra.Command{
		Use:   "prog",
		Short: "EEPROM Programmer",
	}

	interactiveCommand = &cobra.Command{
		Use:   "inter",
		Short: "Interactive mode",
		Run: func(cmd *cobra.Command, args []string) {
			uart = handler.InitUartHandler(Port)
			eeprom := handler.EEPROMHandler{
				Uart: uart,
				Size: handler.S32k,
			}
			handler.Interactive(eeprom)
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of prog",
		Long:  `All software has versions. This is prog's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Prog v1.0")
		},
	}

	serialList = &cobra.Command{
		Use:   "serial",
		Short: "List the available serial ports",
		Run: func(cmd *cobra.Command, args []string) {
			ports, err := serial.GetPortsList()
			if err != nil {
				log.Fatal(err)
			}
			if len(ports) == 0 {
				log.Fatal("No serial ports found!")
			}
			fmt.Println("Ports:")
			for _, port := range ports {
				fmt.Printf("\t%v\n", port)
			}
		},
	}

	readCommand = &cobra.Command{
		Use:   "read",
		Short: "read from EEPROM to file",
		Run: func(cmd *cobra.Command, args []string) {
			uart = handler.InitUartHandler(Port)
			eeprom := handler.EEPROMHandler{
				Uart: uart,
				Size: handler.S32k,
			}
			eeprom.Read(FileName)
		},
	}

	writeCommand = &cobra.Command{
		Use:   "write",
		Short: "write to EEPROM from file",
		Run: func(cmd *cobra.Command, args []string) {
			uart = handler.InitUartHandler(Port)
			eeprom := handler.EEPROMHandler{
				Uart: uart,
				Size: handler.S32k,
			}
			eeprom.Write(FileName)
		},
	}

	checkCommand = &cobra.Command{
		Use:   "check",
		Short: "check EEPROM with the given file",
		Run: func(cmd *cobra.Command, args []string) {
			uart = handler.InitUartHandler(Port)
			eeprom := handler.EEPROMHandler{
				Uart: uart,
				Size: handler.S32k,
			}
			eeprom.Check(FileName)
		},
	}

	pingCommand = &cobra.Command{
		Use:   "ping",
		Short: "ping programmer",
		Run: func(cmd *cobra.Command, args []string) {
			uart = handler.InitUartHandler(Port)
			eeprom := handler.EEPROMHandler{
				Uart: uart,
				Size: handler.S32k,
			}
			eeprom.Ping()
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	readCommand.Flags().StringVarP(&Port, "port", "p", "", "Serial port")
	readCommand.MarkFlagRequired("port")
	writeCommand.Flags().StringVarP(&Port, "port", "p", "", "Serial port")
	writeCommand.MarkFlagRequired("port")
	checkCommand.Flags().StringVarP(&Port, "port", "p", "", "Serial port")
	checkCommand.MarkFlagRequired("port")
	pingCommand.Flags().StringVarP(&Port, "port", "p", "", "Serial port")
	pingCommand.MarkFlagRequired("port")
	interactiveCommand.Flags().StringVarP(&Port, "port", "p", "", "Serial port")
	interactiveCommand.MarkFlagRequired("port")

	readCommand.Flags().StringVarP(&FileName, "file", "f", "", "Filename to save")
	checkCommand.Flags().StringVarP(&FileName, "file", "f", "", "Filename to check")
	writeCommand.Flags().StringVarP(&FileName, "file", "f", "", "Filename to read")

	readCommand.MarkFlagRequired("file")
	checkCommand.MarkFlagRequired("file")
	writeCommand.MarkFlagRequired("file")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serialList)

	rootCmd.AddCommand(readCommand)
	rootCmd.AddCommand(writeCommand)
	rootCmd.AddCommand(checkCommand)
	rootCmd.AddCommand(pingCommand)
	rootCmd.AddCommand(interactiveCommand)
}
