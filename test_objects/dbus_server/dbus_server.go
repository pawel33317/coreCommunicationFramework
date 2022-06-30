package main

import (
	"fmt"
	"os"

	"github.com/godbus/dbus/v5"
)

//!!!!!!!!!!!!!!!! for dbus in wsl setup .bash_sysinit

type exportedStringType string
type exportedEmptyType struct{}

func (memberParam exportedStringType) ExportedMethod1(u uint64) (string, *dbus.Error) {
	fmt.Println("ExportedMethod1 called, received param:", u, ", memberParam:", memberParam)
	return string("Return val from ExportedMethod1"), nil
}

func (memberParam exportedStringType) ExportedMethod2(u uint64) (string, *dbus.Error) {
	fmt.Println("ExportedMethod2 called, received param:", u, ", memberParam:", memberParam)
	return string("Return val from ExportedMethod2"), nil
}

func (exportedEmptyType) ExportedMethodForEmptyObj() (uint, *dbus.Error) {
	fmt.Println("ExportedMethodForEmptyObj called")
	return 1, nil
}

func main() {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	f := exportedStringType("exportedStringType")
	conn.Export(f, "/github/com/pawel33317/CoreCommunicationFramework", "github.com.pawel33317.coreCommunicationFramework")

	f2 := exportedEmptyType{}
	conn.Export(f2, "/github/com/pawel33317/CoreCommunicationFramework2", "github.com.pawel33317.coreCommunicationFramework")

	reply, err := conn.RequestName("github.com.pawel33317.coreCommunicationFramework", dbus.NameFlagDoNotQueue)
	if err != nil {
		panic(err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		fmt.Fprintln(os.Stderr, "name already taken")
		os.Exit(1)
	}

	fmt.Println("Listening on github.com.pawel33317.coreCommunicationFramework / /github/com/pawel33317/CoreCommunicationFramework ...")
	select {}
}
