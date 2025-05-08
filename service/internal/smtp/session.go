package smtp

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

const (
	StateInit = iota
	StateEHLO
	StateMail
	StateRcpt
	StateData
	StateQuit
)

type Session struct {
	conn       net.Conn
	reader     *bufio.Reader
	mailFrom   string
	recipients []string
	data       string
	mailData   []byte
	state      int
	handler    *Handler
}

func NewSession(conn net.Conn) *Session {
	return &Session{
		conn:       conn,
		reader:     bufio.NewReader(conn),
		recipients: make([]string, 0),
		state:      StateInit,
		handler:    NewHandler(),
	}
}

func (s *Session) Start() {
	s.sendResponse(220, "SMTP Server Ready")

	for {
		line, err := s.readLine()
		if err != nil {
			log.Printf("Erro at line: %v", err)
			return
		}

		if s.processCommand(line) {
			return
		}
	}
}

func (s *Session) readLine() (string, error) {
	s.conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	line, err := s.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(line, "\r\n"), nil

}

func (s *Session) sendResponse(code int, message string) {
	response := fmt.Sprintf("%d %s\r\n", code, message)
	s.conn.Write([]byte(response))
}

func (s *Session) processCommand(line string) bool {
	parts := strings.SplitN(line, " ", 2)
	cmd := strings.ToUpper(parts[0])

	var arg string
	if len(parts) > 1 {
		arg = parts[1]
	}
	return s.handler.HandleCommand(s, cmd, arg)
}
