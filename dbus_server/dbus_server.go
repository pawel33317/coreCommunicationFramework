package dbus_server

import (
	"github.com/godbus/dbus/v5"
	"github.com/pawel33317/coreCommunicationFramework/app_state_manager"
	"github.com/pawel33317/coreCommunicationFramework/app_state_manager/app_state"
	"github.com/pawel33317/coreCommunicationFramework/logger"
)

//!!!!!!!!!!!!!!!! for dbus in wsl setup .bash_sysinit
//Use dbus_client.go for testing

var globalLogger *logger.LogWrapper = nil

type DBUSServer struct {
	closeSignalChan chan struct{}
	log             logger.LogWrapper
	asm             app_state_manager.AppStateClientHandler
}

func NewDBUSServer(logg logger.Logger, asmClient app_state_manager.AppStateClientHandler) *DBUSServer {
	dbusServer := &DBUSServer{
		closeSignalChan: make(chan struct{}),
		log:             *logger.NewLogWrapper(logg, "DBSS"),
		asm:             asmClient,
	}
	globalLogger = &dbusServer.log
	asmClient.RegisterObserver(dbusServer)
	return dbusServer
}

func (dbusServer *DBUSServer) OnAppStateChanged(state app_state.State) {
	switch state {
	case app_state.INITIALIZED:
		dbusServer.asm.RegisterLockState(dbusServer, app_state.SHUTDOWN)
	case app_state.CONFIGURED:
		dbusServer.log.Log(logger.DEBUG, "CONFIGURED received")
		dbusServer.RunDBUSServer()
	case app_state.SHUTDOWN:
		dbusServer.log.Log(logger.DEBUG, "SHUTDOWN received")
		dbusServer.closeSignalChan <- struct{}{}
		<-dbusServer.closeSignalChan
		dbusServer.asm.UnlockState(dbusServer)
	case app_state.DISABLED:
	}
}

type exportedStringType string
type exportedEmptyType struct{}

func (memberParam exportedStringType) ExportedMethod1(u uint64) (string, *dbus.Error) {
	globalLogger.Log(logger.INFO, "ExportedMethod1 called, received param:", u, ", memberParam:", memberParam)
	return string("Return val from ExportedMethod1"), nil
}

func (memberParam exportedStringType) ExportedMethod2(u uint64) (string, *dbus.Error) {
	globalLogger.Log(logger.INFO, "ExportedMethod2 called, received param:", u, ", memberParam:", memberParam)
	return string("Return val from ExportedMethod2"), nil
}

func (exportedEmptyType) ExportedMethodForEmptyObj() (uint, *dbus.Error) {
	globalLogger.Log(logger.INFO, "ExportedMethodForEmptyObj called")
	return 1, nil
}

func (dbusServer *DBUSServer) RunDBUSServer() {
	DBUS_PATH := "/github/com/pawel33317/CoreCommunicationFramework"
	DBUS_IFACE := "github.com.pawel33317.coreCommunicationFramework"
	dbusServer.log.Log(logger.INFO, "Starting DBUS server", DBUS_PATH, DBUS_IFACE)

	runDbusThread := func() {
		conn, err := dbus.ConnectSessionBus()
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		f := exportedStringType("exportedStringType")
		conn.Export(f, dbus.ObjectPath(DBUS_PATH), DBUS_IFACE)

		f2 := exportedEmptyType{}
		conn.Export(f2, dbus.ObjectPath(DBUS_PATH+"2"), DBUS_IFACE)

		reply, err := conn.RequestName(DBUS_IFACE, dbus.NameFlagDoNotQueue)
		if err != nil {
			panic(err)
		}
		if reply != dbus.RequestNameReplyPrimaryOwner {
			globalLogger.Log(logger.FATAL, "name already taken")
			panic("name already taken")
		}

		<-dbusServer.closeSignalChan
		dbusServer.log.Log(logger.INFO, "Stopping DBUS server")
		dbusServer.closeSignalChan <- struct{}{}
	}

	go runDbusThread()
}
