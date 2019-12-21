package tcp

import (
	"net"
	"sync"
	"time"
)

const (
	defalutDeadline      = time.Second
	defaultReadDeadline  = time.Second
	defaultWriteDeadline = time.Second
	defaultChanLimit     = 4
)

type Packet interface {
	Serialize() []byte
}

type Protocol interface {
	ReadPacket(conn *net.TCPConn) (Packet, error)
}

type Options struct {
	Host          string
	Port          string
	Deadline      time.Duration
	ReadDeadline  time.Duration
	WriteDeadline time.Duration
	// the limit of packet send channel
	PacketSendChanLimit uint32
	// the limit of packet receive channel
	PacketReceiveChanLimit uint32
}

type Server struct {
	options   *Options        // server configuration
	callback  ConnCallback    // message callbacks in connection
	protocol  Protocol        // customize packet protocol
	exitChan  chan struct{}   // notify all goroutines to shutdown
	waitGroup *sync.WaitGroup // wait for all goroutines
}

func SampleOptions(host, port string) *Options {
	return &Options{
		Host:                   host,
		Port:                   port,
		Deadline:               defalutDeadline,
		ReadDeadline:           defaultReadDeadline,
		WriteDeadline:          defaultWriteDeadline,
		PacketSendChanLimit:    defaultChanLimit,
		PacketReceiveChanLimit: defaultChanLimit,
	}
}

// NewServer creates a server
func NewServer(options *Options, callback ConnCallback, protocol Protocol) *Server {
	return &Server{
		options:   options,
		callback:  callback,
		protocol:  protocol,
		exitChan:  make(chan struct{}),
		waitGroup: &sync.WaitGroup{},
	}
}

func (s *Server) Listener() (*net.TCPListener, error) {
	addr := net.JoinHostPort(s.options.Host, s.options.Port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	listen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, err
	}
	return listen, nil
}

// Start starts service
func (s *Server) Start() error {
	listen, err := s.Listener()
	if err != nil {
		return nil
	}
	s.waitGroup.Add(1)
	defer func() {
		listen.Close()
		s.waitGroup.Done()
	}()

	for {
		select {
		case <-s.exitChan:
			return nil
		default:
		}
		//it will stick if deadline had not set
		listen.SetDeadline(time.Now().Add(s.options.Deadline))
		conn, err := listen.AcceptTCP()
		if err != nil {
			continue
		}
		s.waitGroup.Add(1)
		go func() {
			newConn(conn, s).Do()
			s.waitGroup.Done()
		}()
	}
}

// Stop stops service
func (s *Server) Stop() {
	close(s.exitChan)
	s.waitGroup.Wait()
}
