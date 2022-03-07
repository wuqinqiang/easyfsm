package log

import (
	"testing"
)

func TestInfo(t *testing.T) {
	logger := DefaultLogger
	logger.Log(LevelInfo, "key1", "value1")
}
