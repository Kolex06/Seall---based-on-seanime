package continuity

import (
	"path/filepath"
	"seall/internal/database/db"
	"seall/internal/util"
	"seall/internal/util/filecache"
	"testing"

	"github.com/stretchr/testify/require"
)

func NewTestManager(t *testing.T, db *db.Database) *Manager {
	logger := util.NewLogger()
	cacher, err := filecache.NewCacher(filepath.Join(t.TempDir(), "cache"))
	require.NoError(t, err)

	manager := NewManager(&NewManagerOptions{
		FileCacher: cacher,
		Logger:     logger,
		Database:   db,
	})

	return manager
}
