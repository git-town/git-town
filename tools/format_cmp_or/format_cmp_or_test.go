package main_test

import (
	"testing"

	formatCmpOr "github.com/git-town/git-town/tools/format_cmp_or"
	"github.com/shoenig/test/must"
)

func TestFormatFileContent(t *testing.T) {
	t.Parallel()

	t.Run("case-insensitive sorting", func(t *testing.T) {
		t.Parallel()
		give := []byte(`
package main

import "cmp"

func foo() {
	if err := cmp.Or(errVerbose, errAutoResolve, errDryRun); err != nil {}
}`[1:])
		want := []byte(`
package main

import "cmp"

func foo() {
	if err := cmp.Or(errAutoResolve, errDryRun, errVerbose); err != nil {}
}`[1:])
		have := formatCmpOr.FormatFileContent("", give)
		must.EqOp(t, string(want), string(have))
	})

	t.Run("multi-line already sorted", func(t *testing.T) {
		t.Parallel()
		give := []byte("package main\n\nimport \"cmp\"\n\nfunc foo() {\n\terr := cmp.Or(\n\t\terrA,\n\t\terrB,\n\t\terrC,\n\t)\n}")
		have := formatCmpOr.FormatFileContent("", give)
		must.EqOp(t, string(give), string(have))
	})

	t.Run("multi-line unsorted", func(t *testing.T) {
		t.Parallel()
		give := []byte("package main\n\nimport \"cmp\"\n\nfunc foo() {\n\terr := cmp.Or(\n\t\terrC,\n\t\terrA,\n\t\terrB,\n\t)\n}")
		want := []byte("package main\n\nimport \"cmp\"\n\nfunc foo() {\n\terr := cmp.Or(\n\t\terrA,\n\t\terrB,\n\t\terrC,\n\t)\n}")
		have := formatCmpOr.FormatFileContent("", give)
		must.EqOp(t, string(want), string(have))
	})

	t.Run("non-identifier args are skipped", func(t *testing.T) {
		t.Parallel()
		give := []byte(`
package main

import "cmp"

func foo() {
	v := cmp.Or(foo.Bar, foo.Baz)
	_ = v
}`[1:])
		have := formatCmpOr.FormatFileContent("", give)
		must.EqOp(t, string(give), string(have))
	})

	t.Run("single argument is skipped", func(t *testing.T) {
		t.Parallel()
		give := []byte(`
package main

import "cmp"

func foo() {
	v := cmp.Or(errA)
	_ = v
}`[1:])
		have := formatCmpOr.FormatFileContent("", give)
		must.EqOp(t, string(give), string(have))
	})

	t.Run("single-line already sorted", func(t *testing.T) {
		t.Parallel()
		give := []byte(`
package main

import "cmp"

func foo() {
	if err := cmp.Or(errA, errB, errC); err != nil {}
}`[1:])
		have := formatCmpOr.FormatFileContent("", give)
		must.EqOp(t, string(give), string(have))
	})

	t.Run("single-line unsorted", func(t *testing.T) {
		t.Parallel()
		give := []byte(`
package main

import "cmp"

func foo() {
	if err := cmp.Or(errC, errA, errB); err != nil {}
}`[1:])
		want := []byte(`
package main

import "cmp"

func foo() {
	if err := cmp.Or(errA, errB, errC); err != nil {}
}`[1:])
		have := formatCmpOr.FormatFileContent("", give)
		must.EqOp(t, string(want), string(have))
	})
}
