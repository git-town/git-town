package option_test

import (
	"errors"
	"testing"

	"github.com/git-town/git-town/v14/src/gohacks/option"
	"github.com/shoenig/test/must"
)

func TestUnpack(t *testing.T) {
	t.Parallel()
	t.Run("option of string", func(t *testing.T) {
		t.Parallel()
		t.Run("populated option of string", func(t *testing.T) {
			t.Parallel()
			text := "hello"
			textOpt := &text
			have, err := option.Unpack(textOpt, errors.New("error"))
			must.NoError(t, err)
			must.EqOp(t, text, have)
		})
		t.Run("empty option of string", func(t *testing.T) {
			t.Parallel()
			var textOpt *string = nil
			domainErr := errors.New("error")
			_, haveErr := option.Unpack(textOpt, domainErr)
			must.Eq(t, domainErr, haveErr)
		})
	})
	t.Run("option of struct", func(t *testing.T) {
		t.Parallel()
		type domainType struct {
			id int
		}

		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			domainObj := domainType{1}
			domainOpt := &domainObj
			have, err := option.Unpack(domainOpt, errors.New("error"))
			must.NoError(t, err)
			must.EqOp(t, domainObj, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			var textOpt *domainType = nil
			domainErr := errors.New("error")
			_, haveErr := option.Unpack(textOpt, domainErr)
			must.Eq(t, domainErr, haveErr)
		})
	})
}
