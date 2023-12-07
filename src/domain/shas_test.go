package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/shoenig/test/must"
)

func TestSHAs(t *testing.T) {
	t.Parallel()

	t.Run("Join", func(t *testing.T) {
		t.Parallel()
		t.Run("contains elements", func(t *testing.T) {
			t.Parallel()
			give := domain.SHAs{
				domain.NewSHA("111111"),
				domain.NewSHA("222222"),
			}
			have := give.Join(", ")
			want := "111111, 222222"
			must.EqOp(t, want, have)
		})
		t.Run("empty list", func(t *testing.T) {
			t.Parallel()
			give := domain.SHAs{}
			have := give.Join(", ")
			want := ""
			must.EqOp(t, want, have)
		})
	})

	t.Run("Strings", func(t *testing.T) {
		t.Parallel()
		t.Run("contains elements", func(t *testing.T) {
			t.Parallel()
			give := domain.SHAs{
				domain.NewSHA("111111"),
				domain.NewSHA("222222"),
			}
			have := give.Strings()
			want := []string{"111111", "222222"}
			must.Eq(t, want, have)
		})
		t.Run("empty list", func(t *testing.T) {
			t.Parallel()
			give := domain.SHAs{}
			have := give.Strings()
			want := []string{}
			must.Eq(t, want, have)
		})
	})
}
