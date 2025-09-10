package pg

import (
	"testing"

	"github.com/linuxunsw/vote/backend/internal/store/pg/harness"
)

func TestMain(m *testing.M) {
	harness.HarnessMain(m)
}
