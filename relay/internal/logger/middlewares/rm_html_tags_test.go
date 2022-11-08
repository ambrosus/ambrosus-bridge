package middlewares

import (
	"regexp"
	"testing"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger"
	mock_logger "github.com/ambrosus/ambrosus-bridge/relay/internal/logger/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRemoveHTMLTags_Log(t *testing.T) {
	tests := []struct {
		name    string
		regexps []*ReplacableRegexp
		log     *logger.ExtLog
		wantLog *logger.ExtLog
	}{
		{
			name: "Ok",
			regexps: []*ReplacableRegexp{
				{regexp.MustCompile(`<script>.*</script>|<.*?>`), ""},
				{regexp.MustCompile(`(?m)^\s+$`), ""},
				{regexp.MustCompile(`(\r\n?|\n){2,}`), ""},
			},

			log:     &logger.ExtLog{Error: `<html><body>test</body></html>`},
			wantLog: &logger.ExtLog{Error: `test`},
		},
		{
			name:    "Ok (real data)",
			regexps: DefaultApplicableRegexps,
			log: &logger.ExtLog{Error: `send event: encodePoSAProof: saveBlocks: fetchBlockCache: 504 Gateway Timeout: <html>
<head><title>504 Gateway Time-out</title></head>
<body>
<center><h1>504 Gateway Time-out</h1></center>
</body>
</html>`},
			wantLog: &logger.ExtLog{Error: `send event: encodePoSAProof: saveBlocks: fetchBlockCache: 504 Gateway Timeout: 
504 Gateway Time-out504 Gateway Time-out`},
		},
		{
			name: "There's no html (should do nothing)",
			regexps: []*ReplacableRegexp{
				{regexp.MustCompile(`tag`), "mag"}, // shouldn't apply this
			},
			log:     &logger.ExtLog{Error: `there's no html tag`},
			wantLog: &logger.ExtLog{Error: `there's no html tag`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			hook := mock_logger.NewMockHook(c)
			hook.EXPECT().Log(tt.wantLog).Times(1)

			md := &RemoveHTMLTags{
				ApplicableRegexps: tt.regexps,
				Hook:              hook,
			}
			md.Log(tt.log)

			assert.Equal(t, tt.wantLog, tt.log)
		})
	}
}
