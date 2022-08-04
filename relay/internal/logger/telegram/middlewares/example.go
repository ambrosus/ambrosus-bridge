package middlewares

import (
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger"
)

type ExampleMiddleware struct {
	next logger.Hook
}

func NewExampleMiddleware(next logger.Hook) *ExampleMiddleware {
	return &ExampleMiddleware{next: next}
}

func (m *ExampleMiddleware) Log(l *logger.ExtLog) {
	fmt.Println("before")
	m.next.Log(l)
	fmt.Println("after")
}
