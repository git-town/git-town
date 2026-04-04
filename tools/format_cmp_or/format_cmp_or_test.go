package main_test

import (
	"testing"

	formatCmpOr "github.com/git-town/git-town/tools/format_cmp_or"
	"github.com/shoenig/test/must"
)

func TestFormatFileContent(t *testing.T) {
	t.Parallel()

	t.Run("single-line unsorted", func(t *testing.T) {
		t.Parallel()
		give := []byte(`package main

import "cmp"

func foo() {
	if err := cmp.Or(errC, errA, errB); err != nil {}
}`)
		want := []byte(`package main

import "cmp"

func foo() {
	if err := cmp.Or(errA, errB, errC); err != nil {}
}`)
		have, err := formatCmpOr.FormatFileContent("", give)
		must.NoError(t, err)
		must.EqOp(t, string(want), string(have))
	})

	t.Run("single-line already sorted", func(t *testing.T) {
		t.Parallel()
		give := []byte(`package main

import "cmp"

func foo() {
	if err := cmp.Or(errA, errB, errC); err != nil {}
}`)
		have, err := formatCmpOr.FormatFileContent("", give)
		must.NoError(t, err)
		must.EqOp(t, string(give), string(have))
	})

	t.Run("multi-line unsorted", func(t *testing.T) {
		t.Parallel()
		give := []byte("package main\n\nimport \"cmp\"\n\nfunc foo() {\n\terr := cmp.Or(\n\t\terrC,\n\t\terrA,\n\t\terrB,\n\t)\n}")
		want := []byte("package main\n\nimport \"cmp\"\n\nfunc foo() {\n\terr := cmp.Or(\n\t\terrA,\n\t\terrB,\n\t\terrC,\n\t)\n}")
		have, err := formatCmpOr.FormatFileContent("", give)
		must.NoError(t, err)
		must.EqOp(t, string(want), string(have))
	})

	t.Run("multi-line already sorted", func(t *testing.T) {
		t.Parallel()
		give := []byte("package main\n\nimport \"cmp\"\n\nfunc foo() {\n\terr := cmp.Or(\n\t\terrA,\n\t\terrB,\n\t\terrC,\n\t)\n}")
		have, err := formatCmpOr.FormatFileContent("", give)
		must.NoError(t, err)
		must.EqOp(t, string(give), string(have))
	})

	t.Run("non-identifier args are skipped", func(t *testing.T) {
		t.Parallel()
		give := []byte(`package main

import "cmp"

func foo() {
	v := cmp.Or(foo.Bar, foo.Baz)
	_ = v
}`)
		have, err := formatCmpOr.FormatFileContent("", give)
		must.NoError(t, err)
		must.EqOp(t, string(give), string(have))
	})

	t.Run("keep-sorted calls are skipped", func(t *testing.T) {
		t.Parallel()
		give := []byte("package main\n\nimport \"cmp\"\n\nfunc foo() {\n\terr := cmp.Or(\n\t\t// keep-sorted start\n\t\terrC,\n\t\terrA,\n\t\t// keep-sorted end\n\t)\n\t_ = err\n}")
		have, err := formatCmpOr.FormatFileContent("myfile.go", give)
		must.NoError(t, err)
		must.EqOp(t, string(give), string(have))
	})

	t.Run("single argument is skipped", func(t *testing.T) {
		t.Parallel()
		give := []byte(`package main

import "cmp"

func foo() {
	v := cmp.Or(errA)
	_ = v
}`)
		have, err := formatCmpOr.FormatFileContent("", give)
		must.NoError(t, err)
		must.EqOp(t, string(give), string(have))
	})

	t.Run("case-insensitive sorting", func(t *testing.T) {
		t.Parallel()
		give := []byte(`package main

import "cmp"

func foo() {
	if err := cmp.Or(errVerbose, errAutoResolve, errDryRun); err != nil {}
}`)
		want := []byte(`package main

import "cmp"

func foo() {
	if err := cmp.Or(errAutoResolve, errDryRun, errVerbose); err != nil {}
}`)
		have, err := formatCmpOr.FormatFileContent("", give)
		must.NoError(t, err)
		must.EqOp(t, string(want), string(have))
	})
}
