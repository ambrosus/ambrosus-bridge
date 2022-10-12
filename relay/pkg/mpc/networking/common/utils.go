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
	if e.Err != nil {
		if e.Type != needType {
			return fmt.Errorf("not %v error: %w", needType, e.Err)
		} else {
			return fmt.Errorf("%v error: %w", needType, e.Err)
		}
	}
	return nil
}
