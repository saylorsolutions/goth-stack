package audit

import "log"

type LogDelegate interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}

func StdDelegate(logger *log.Logger, enableDebug bool) LogDelegate {
	return &stdLogDelegate{
		logger:       logger,
		debugEnabled: enableDebug,
	}
}

type stdLogDelegate struct {
	logger       *log.Logger
	debugEnabled bool
}

func (s *stdLogDelegate) Info(msg string, args ...any) {
	s.logger.Printf("[INF] "+msg+"\n", args...)
}

func (s *stdLogDelegate) Error(msg string, args ...any) {
	s.logger.Printf("[ERR] "+msg+"\n", args...)
}

func (s *stdLogDelegate) Debug(msg string, args ...any) {
	if !s.debugEnabled {
		return
	}
	s.logger.Printf("[DBG] "+msg+"\n", args...)
}
