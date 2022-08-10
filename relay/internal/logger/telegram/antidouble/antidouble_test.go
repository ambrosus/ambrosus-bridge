package antidouble

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger/telegram"
	mocks_telegram "github.com/ambrosus/ambrosus-bridge/relay/internal/logger/telegram/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setStaticNowFunc(t *testing.T) time.Time {
	t.Helper()

	now := time.Now()
	nowFunc = func() time.Time { return now }

	return now
}

// checks:
// 	`firstSentAt` (in `checkEveryLogFunc`)
// 	`sentCount` (in `checkAfterLogsFunc` and in some cases in `checkEveryLogFunc`)
func Test_AntiDoubleTgLogger_Log_send_or_edit_expired_log(t *testing.T) {
	const numberOfLoggings = 5
	extLog := &logger.ExtLog{
		Level:   "error",
		Bridge:  "amb",
		Service: "lalala",
		Message: "test test test",
		Error:   "test error msg",
	}
	msg := telegram.BuildMessage(extLog)

	tests := []struct {
		name               string
		mock               func(tgLogger *mocks_telegram.MockTgLoggerI)
		mockBeforeNextLog  func(firstSentAtBeforeNextLog time.Time)
		checkEveryLogFunc  func(t *testing.T, firstSentAtBeforeNextLog time.Time, antiDoubleTgLogger *AntiDoubleTgLogger)
		checkAfterLogsFunc func(t *testing.T, antiDoubleTgLogger *AntiDoubleTgLogger)
	}{
		{
			name: "Log is expired (should send that again)",
			mock: func(tgLogger *mocks_telegram.MockTgLoggerI) {
				tgLogger.EXPECT().Send(gomock.Eq(msg)).Return([]uint64{1}, []string{msg}, nil).Times(1 * numberOfLoggings)
			},
			mockBeforeNextLog: func(firstSentAtBeforeNextLog time.Time) {
				nowFunc = func() time.Time { return firstSentAtBeforeNextLog.Add(logExpireTime + 1) }
			},
			checkEveryLogFunc: func(t *testing.T, firstSentAtBeforeNextLog time.Time, antiDoubleTgLogger *AntiDoubleTgLogger) {
				log := antiDoubleTgLogger.logs[msg]
				assert.Equal(t, nowFunc(), log.firstSentAt)
				assert.Equal(t, uint64(1), log.sentCount)
			},
			checkAfterLogsFunc: func(t *testing.T, antiDoubleTgLogger *AntiDoubleTgLogger) {
				assert.Equal(t, uint64(1), antiDoubleTgLogger.logs[msg].sentCount)
			},
		},
		{
			name: "Log is not expired (should edit that msg)",
			mock: func(tgLogger *mocks_telegram.MockTgLoggerI) {
				tgLogger.EXPECT().Send(gomock.Eq(msg)).Return([]uint64{1}, []string{msg}, nil).Times(1)
				tgLogger.EXPECT().EditText(gomock.Any(), gomock.Any()).Return(nil).Times(numberOfLoggings - 1)
			},
			mockBeforeNextLog: func(firstSentAtBeforeNextLog time.Time) {
				nowFunc = func() time.Time { return firstSentAtBeforeNextLog.Add(logExpireTime - 1) }
			},
			checkEveryLogFunc: func(t *testing.T, firstSentAtBeforeNextLog time.Time, antiDoubleTgLogger *AntiDoubleTgLogger) {
				assert.Equal(t, firstSentAtBeforeNextLog, antiDoubleTgLogger.logs[msg].firstSentAt)
			},
			checkAfterLogsFunc: func(t *testing.T, antiDoubleTgLogger *AntiDoubleTgLogger) {
				assert.Equal(t, uint64(numberOfLoggings), antiDoubleTgLogger.logs[msg].sentCount)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			setStaticNowFunc(t)

			tgLogger := mocks_telegram.NewMockTgLoggerI(c)
			antiDoubleTgLogger := NewAntiDoubleTgLogger(tgLogger)
			tt.mock(tgLogger)

			// first log
			antiDoubleTgLogger.Log(extLog)

			// other logs
			for i := 0; i < numberOfLoggings-1; i++ {
				firstSentAtBeforeNextLog := antiDoubleTgLogger.logs[msg].firstSentAt
				tt.mockBeforeNextLog(firstSentAtBeforeNextLog)

				antiDoubleTgLogger.Log(extLog)
				tt.checkEveryLogFunc(t, firstSentAtBeforeNextLog, antiDoubleTgLogger)
			}
			tt.checkAfterLogsFunc(t, antiDoubleTgLogger)
		})
	}
}

func Test_AntiDoubleTgLogger_clearExpiredLogs(t *testing.T) {
	const shouldBeCleared = "should be cleared"
	const shouldNotBeCleared = "should not be cleared"

	var uniqueKeyCounter uint
	genUniqueKey := func(s string) string {
		res := s + fmt.Sprint(uniqueKeyCounter)
		uniqueKeyCounter++
		return res
	}

	c := gomock.NewController(t)
	defer c.Finish()
	now := setStaticNowFunc(t)

	antiDoubleTgLogger := NewAntiDoubleTgLogger(mocks_telegram.NewMockTgLoggerI(c))
	antiDoubleTgLogger.logs = map[string]LogMsg{
		genUniqueKey(shouldBeCleared):    {firstSentAt: now.Add(-(logExpireTime + 1))},
		genUniqueKey(shouldBeCleared):    {firstSentAt: now.Add(-(logExpireTime + 2))},
		genUniqueKey(shouldNotBeCleared): {firstSentAt: now.Add(-(logExpireTime))},
		genUniqueKey(shouldNotBeCleared): {firstSentAt: now.Add(-(logExpireTime))},
		genUniqueKey(shouldNotBeCleared): {firstSentAt: now.Add(-(logExpireTime - 1))},
		genUniqueKey(shouldNotBeCleared): {firstSentAt: now.Add(-(logExpireTime - 2))},
	}

	expectedLogs := map[string]LogMsg{}
	for k, v := range antiDoubleTgLogger.logs {
		if strings.HasPrefix(k, "should not be cleared") {
			expectedLogs[k] = v
		}
	}

	antiDoubleTgLogger.clearExpiredLogs()
	assert.Equal(t, expectedLogs, antiDoubleTgLogger.logs)
}

func Test_LogMsg_isExpired(t *testing.T) {
	now := setStaticNowFunc(t)

	tests := []struct {
		name        string
		firstSentAt time.Time
		want        bool
	}{
		{
			name:        "True",
			firstSentAt: now.Add(-(logExpireTime + 1)),
			want:        true,
		},
		{
			name:        "False (equal)",
			firstSentAt: now.Add(-logExpireTime),
			want:        false,
		},
		{
			name:        "False",
			firstSentAt: now.Add(-(logExpireTime - 1)),
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logMsg := LogMsg{firstSentAt: tt.firstSentAt}

			assert.Equal(t, tt.want, logMsg.isExpired())
		})
	}
}
