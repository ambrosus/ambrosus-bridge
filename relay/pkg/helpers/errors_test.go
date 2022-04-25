package helpers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_decodeRevertMessage(t *testing.T) {
	type args struct {
		errStr string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Case 1",
			args{"0x08c379a00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000002750726f7669646564206164647265737320697320616c726561647920612076616c696461746f7200000000000000000000000000000000000000000000000000"},
			"Provided address is already a validator",
			false,
		},
		{"Case 2",
			args{"0x08c379a000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000094416363657373436f6e74726f6c3a206163636f756e7420307832393563323730373331396164346265636136623562623430383636313766643666323430633165206973206d697373696e6720726f6c6520307830373761316435323661346365386137373336333261623133623466626266316663633935346333646162323663643237656130653261363735306461356437000000000000000000000000"},
			"AccessControl: account 0x295c2707319ad4beca6b5bb4086617fd6f240c1e is missing role 0x077a1d526a4ce8a773632ab13b4fbbf1fcc954c3dab26cd27ea0e2a6750da5d7",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decodedMsg, err := decodeRevertMessage(tt.args.errStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeRevertMessage() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			assert.Equalf(t, tt.want, decodedMsg, "decodeRevertMessage(%v)", tt.args.errStr)
		})
	}
}

func Test_parseError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			"Reverted error",
			args{NewTestError("The execution failed due to an exception.", "Reverted 0x08c379a00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000002750726f7669646564206164647265737320697320616c726561647920612076616c696461746f7200000000000000000000000000000000000000000000000000")},
			fmt.Errorf("The execution failed due to an exception. Provided address is already a validator"),
		},
		{
			"VM error",
			args{NewTestError("VM execution error.", "Out of gas")},
			fmt.Errorf("VM execution error. Out of gas"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want.Error(), ParseError(tt.args.err).Error(), "parseError(%v)", tt.args.err)
		})
	}
}

type testError struct {
	message string
	data    string
}

func NewTestError(message string, data string) *testError {
	return &testError{
		message: message,
		data:    data,
	}
}

func (e *testError) Error() string {
	return e.message
}

func (e *testError) ErrorData() interface{} {
	return e.data
}
