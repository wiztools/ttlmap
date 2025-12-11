package ttlmap_test

import (
	"context"
	"testing"
	"time"

	"github.com/wiztools/ttlmap"
)

func TestFullRefreshMap(t *testing.T) {
	t.Run("full-refresh-map-create/get", func(t *testing.T) {
		var counter int64 = 0
		m := ttlmap.NewFullRefreshMap(func(ctx context.Context) (map[string]int64, error) {
			counter++
			return map[string]int64{"key": counter}, nil
		}, time.Second)
		ctx := context.Background()
		{
			v, ok, _ := m.Get(ctx, "key")
			if !ok {
				t.Error("expected key to be present")
				return
			}
			if v != 1 {
				t.Error("expected value to be 1")
				return
			}
		}
		time.Sleep(3 * time.Second)
		{
			v, ok, _ := m.Get(ctx, "key")
			if !ok {
				t.Error("expected cache to be refreshed")
				return
			}
			if v != 2 {
				t.Error("expected value to be 2")
				return
			}
		}
	})

}
