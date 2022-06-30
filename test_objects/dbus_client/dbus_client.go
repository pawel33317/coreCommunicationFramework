package main

import (
	"fmt"
	"os"

	"github.com/godbus/dbus/v5"
)

func main() {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
		os.Exit(1)
	}
	defer conn.Close()

	var s string
	obj := conn.Object("github.com.pawel33317.coreCommunicationFramework", "/github/com/pawel33317/CoreCommunicationFramework")
	err = obj.Call("github.com.pawel33317.coreCommunicationFramework.ExportedMethod1", 0, uint64(33)).Store(&s)
	if err != nil {
		fmt.Println(os.Stderr, "Call ExportedMethod1 failed:", err)
		os.Exit(1)
	}
	fmt.Println("ExportedMethod1 result", s)

	err = obj.Call("github.com.pawel33317.coreCommunicationFramework.ExportedMethod2", 0, uint64(1111)).Store(&s)
	if err != nil {
		fmt.Println("Call ExportedMethod2 failed:", err)
		os.Exit(1)
	}
	fmt.Println("ExportedMethod2 result", s)

	var s2 uint
	obj2 := conn.Object("github.com.pawel33317.coreCommunicationFramework", "/github/com/pawel33317/CoreCommunicationFramework2")
	err = obj2.Call("github.com.pawel33317.coreCommunicationFramework.ExportedMethodForEmptyObj", 0).Store(&s2)
	if err != nil {
		fmt.Println("Call ExportedMethodForEmptyObj failed:", err)
		os.Exit(1)
	}
	fmt.Println("Result from ExportedMethodForEmptyObj", s2)
}
