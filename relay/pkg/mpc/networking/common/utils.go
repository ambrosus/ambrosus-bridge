package common

import "fmt"

var (
	KeygenOperation    = []byte("keygen")
	HeaderTssID        = "X-TSS-ID"
	HeaderTssOperation = "X-TSS-Operation"
	ResultPrefix       = []byte("result")
)

type OpError struct {
	Type string
	Err  error
}

func (e OpError) Check(needType string) error {
	if e.Type != needType {
		return fmt.Errorf("%v (need %v) error: %w", e.Type, needType, e.Err)
	}
	if e.Err != nil {
		return fmt.Errorf("%v error: %w", needType, e.Err)
	}
	return nil
}
