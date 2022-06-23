package tcp_server

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/pawel33317/coreCommunicationFramework/app_state_manager"
	"github.com/pawel33317/coreCommunicationFramework/app_state_manager/app_state"
	"github.com/pawel33317/coreCommunicationFramework/logger"
)

type TcpServer struct {
	closeServerChan  chan bool
	closedServerChan chan bool
	dataStreamChan   chan<- string
	log              logger.LogWrapper
	asm              app_state_manager.AppStateClientHandler
	conn             net.Conn     //handler to close connection
	listener         net.Listener //handler to close connection
}

func NewTcpServer(logg logger.Logger, asmClient app_state_manager.AppStateClientHandler, dataStreamChannel chan<- string) *TcpServer {
	tcpServ := &TcpServer{
		closeServerChan:  make(chan bool, 1),
		closedServerChan: make(chan bool, 1),
		log:              *logger.NewLogWrapper(logg, "TCPS"),
		asm:              asmClient,
		dataStreamChan:   dataStreamChannel,
	}
	asmClient.RegisterObserver(tcpServ)
	return tcpServ
}

func (tcpServ *TcpServer) OnAppStateChanged(state app_state.State) {
	switch state {
	case app_state.INITIALIZED:
		tcpServ.asm.RegisterLockState(tcpServ, app_state.SHUTDOWN)
	case app_state.CONFIGURED:
		tcpServ.log.Log(logger.DEBUG, "CONFIGURED received")
		tcpServ.RunTcpServer()
	case app_state.SHUTDOWN:
		tcpServ.log.Log(logger.DEBUG, "SHUTDOWN received")
		tcpServ.closeServerChan <- true
		<-tcpServ.closedServerChan
		tcpServ.asm.UnlockState(tcpServ)
	case app_state.DISABLED:
	}
}

func (tcpServ *TcpServer) RunTcpServer() {
	PORT := "127.0.0.1:8081"
	tcpServ.log.Log(logger.INFO, "Starting tcps server", PORT)

	startConnectionReader := func() {
		for {
			tcpServ.log.Log(logger.DEBUG, "Receiving mode")
			netData, err := bufio.NewReader(tcpServ.conn).ReadString('\n') //TODO close also reader
			if err != nil {
				tcpServ.log.Log(logger.WARN, "Cannot read tcp data", err)
				return
			}
			if strings.TrimSpace(string(netData)) == "STOP" {
				tcpServ.log.Log(logger.INFO, "Client closed connection Exiting TCP server!")
				return
			}

			fmt.Print("TCPClientData-> ", string(netData))
			tcpServ.log.Log(logger.DEBUG, "Responding mode, responding")
			//fmt.Fprintf(tcpServ.conn, "response\n")
			tcpServ.conn.Write([]byte("response\n"))
			tcpServ.dataStreamChan <- string(netData)
		}
	}

	startClientListener := func() {
		listener, err := net.Listen("tcp", PORT)
		tcpServ.listener = listener
		if err != nil {
			tcpServ.log.Log(logger.FATAL, "Cannot start tcp server on port", PORT)
			panic("Cannot start tcp server on port " + PORT + " " + err.Error())
		}
		tcpServ.log.Log(logger.INFO, "TCP server started", PORT)

		tcpServ.log.Log(logger.INFO, "Starting listener")
		tcpServ.conn, err = tcpServ.listener.Accept() //not in while so accept only one listener
		if err != nil {
			tcpServ.log.Log(logger.WARN, "Cannot create tcp connection", err)
			return
		}
		tcpServ.log.Log(logger.INFO, "Stopping listener, client detected")
		startConnectionReader() //start with go if more clients
	}

	observeCloseChannel := func() {
		<-tcpServ.closeServerChan
		tcpServ.log.Log(logger.INFO, "Stopping tcps server")
		if tcpServ.conn != nil {
			tcpServ.log.Log(logger.INFO, "Closing connection")
			tcpServ.conn.Close()
		}
		tcpServ.listener.Close()
		tcpServ.closedServerChan <- true
	}

	go observeCloseChannel()
	go startClientListener()
}
