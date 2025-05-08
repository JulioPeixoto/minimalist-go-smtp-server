package smtp

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type Server struct {
	host     string
	port     int
	listener net.Listener
	wg       sync.WaitGroup
	quit     chan struct{}
}

func NewServer(host string, port int) *Server {
	return &Server{
		host: host,
		port: port,
		quit: make(chan struct{}),
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.listener = listener

	s.wg.Add(1)
	go s.serve()
	return nil
}

func (s *Server) serve() {
	defer s.wg.Done()

	for {
		select {
		case <-s.quit:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				select {
				case <-s.quit:
					return
				default:
					log.Printf("Erro ao aceitar conexão: %v", err)
				}
				continue
			}

			s.wg.Add(1)
			go s.handleConnection(conn)
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	//TODO: create session package
	session := NewSession(conn)
	session.Start()
}

func (s *Server) Stop() error {
	close(s.quit)
	s.listener.Close()

	go func() {
		s.wg.Wait()
	}()

	return nil
}
