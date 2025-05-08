package smtp

import (
	"fmt"
	"log"
	"strings"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) HandleCommand(session *Session, cmd, arg string) bool {
	switch cmd {
	case "HELO", "EHLO":
		h.handleHELO(session, arg)
	case "MAIL":
		h.handleMAIL(session, arg)
	case "RCPT":
		h.handleRCPT(session, arg)
	case "DATA":
		h.handleDATA(session)
	case "QUIT":
		h.handleQUIT(session)
		return true
	default:
		session.sendResponse(502, "Comando não implementado")
	}

	return false
}

func (h *Handler) handleHELO(s *Session, arg string) {
	if arg == "" {
		s.sendResponse(501, "Parâmetro obrigatório")
		return
	}

	s.state = StateEHLO
	s.sendResponse(250, fmt.Sprintf("Olá %s", arg))
}

func (h *Handler) handleMAIL(s *Session, arg string) {
	if s.state < StateEHLO {
		s.sendResponse(503, "HELO/EHLO primeiro")
		return
	}

	if !strings.HasPrefix(strings.ToUpper(arg), "FROM:") {
		s.sendResponse(501, "Sintaxe: MAIL FROM:<endereço>")
		return
	}

	address := strings.Trim(arg[5:], "<>")
	s.mailFrom = address
	s.state = StateMail
	s.sendResponse(250, "OK")
}

func (h *Handler) handleRCPT(s *Session, arg string) {
	if s.state < StateMail {
		s.sendResponse(503, "MAIL FROM primeiro")
		return
	}

	if !strings.HasPrefix(strings.ToUpper(arg), "TO:") {
		s.sendResponse(501, "Sintaxe: RCPT TO:<endereço>")
		return
	}

	address := strings.Trim(arg[3:], "<>")
	s.recipients = append(s.recipients, address)
	s.state = StateRcpt
	s.sendResponse(250, "OK")
}

func (h *Handler) handleDATA(s *Session) {
	if s.state < StateRcpt {
		s.sendResponse(503, "RCPT TO primeiro")
		return
	}

	s.sendResponse(354, "Fim dos dados com <CR><LF>.<CR><LF>")

	var dataLines []string
	for {
		line, err := s.readLine()
		if err != nil {
			return
		}

		if line == "." {
			break
		}

		if strings.HasPrefix(line, "..") {
			line = line[1:]
		}

		dataLines = append(dataLines, line)
	}

	s.data = strings.Join(dataLines, "\n")
	s.state = StateData

	log.Printf("Email recebido de %s para %v", s.mailFrom, s.recipients)
	log.Printf("Conteúdo:\n%s", s.data)

	s.sendResponse(250, "OK - Mensagem aceita")
}

func (h *Handler) handleQUIT(s *Session) {
	s.sendResponse(221, "Fechando conexão")
}
