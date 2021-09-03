package host

import (
	"context"
	"testing"

	"github.com/hashicorp/boundary/internal/db"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
)

func TestPlugin(t *testing.T, conn *gorm.DB, name string, prefix string) *Plugin {
	t.Helper()
	p := NewPlugin(name, prefix, "0.0.1")
	id, err := newPluginId()
	require.NoError(t, err)
	p.PublicId = id

	w := db.New(conn)
	require.NoError(t, w.Create(context.Background(), p))
	return p
}
