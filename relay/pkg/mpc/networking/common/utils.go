package common

var (
	KeygenOperation    = []byte("keygen")
	HeaderTssID        = "X-TSS-ID"
	HeaderTssOperation = "X-TSS-Operation"
	ResultPrefix       = []byte("result")
	EndpointFullMsg    = "/tx"
)

type OpError struct {
	Type string
	Err  error
}
