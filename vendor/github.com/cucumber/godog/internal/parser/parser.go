package parser

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/cucumber/gherkin-go/v19"
	"github.com/cucumber/messages-go/v16"

	"github.com/cucumber/godog/internal/models"
	"github.com/cucumber/godog/internal/tags"
)

var pathLineRe = regexp.MustCompile(`:([\d]+)$`)

// ExtractFeaturePathLine ...
func ExtractFeaturePathLine(p string) (string, int) {
	line := -1
	retPath := p
	if m := pathLineRe.FindStringSubmatch(p); len(m) > 0 {
		if i, err := strconv.Atoi(m[1]); err == nil {
			line = i
			retPath = p[:strings.LastIndexByte(p, ':')]
		}
	}
	return retPath, line
}

func parseFeatureFile(path string, newIDFunc func() string) (*models.Feature, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer reader.Close()

	var buf bytes.Buffer
	gherkinDocument, err := gherkin.ParseGherkinDocument(io.TeeReader(reader, &buf), newIDFunc)
	if err != nil {
		return nil, fmt.Errorf("%s - %v", path, err)
	}

	gherkinDocument.Uri = path
	pickles := gherkin.Pickles(*gherkinDocument, path, newIDFunc)

	f := models.Feature{GherkinDocument: gherkinDocument, Pickles: pickles, Content: buf.Bytes()}
	return &f, nil
}

func parseFeatureDir(dir string, newIDFunc func() string) ([]*models.Feature, error) {
	var features []*models.Feature
	return features, filepath.Walk(dir, func(p string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		if !strings.HasSuffix(p, ".feature") {
			return nil
		}

		feat, err := parseFeatureFile(p, newIDFunc)
		if err != nil {
			return err
		}

		features = append(features, feat)
		return nil
	})
}

func parsePath(path string, newIDFunc func() string) ([]*models.Feature, error) {
	var features []*models.Feature

	path, line := ExtractFeaturePathLine(path)

	fi, err := os.Stat(path)
	if err != nil {
		return features, err
	}

	if fi.IsDir() {
		return parseFeatureDir(path, newIDFunc)
	}

	ft, err := parseFeatureFile(path, newIDFunc)
	if err != nil {
		return features, err
	}

	// filter scenario by line number
	var pickles []*messages.Pickle

	if line != -1 {
		ft.Uri += ":" + strconv.Itoa(line)
	}

	for _, pickle := range ft.Pickles {
		sc := ft.FindScenario(pickle.AstNodeIds[0])

		if line == -1 || int64(line) == sc.Location.Line {
			if line != -1 {
				pickle.Uri += ":" + strconv.Itoa(line)
			}

			pickles = append(pickles, pickle)
		}
	}
	ft.Pickles = pickles

	return append(features, ft), nil
}

// ParseFeatures ...
func ParseFeatures(filter string, paths []string) ([]*models.Feature, error) {
	var order int

	featureIdxs := make(map[string]int)
	uniqueFeatureURI := make(map[string]*models.Feature)
	newIDFunc := (&messages.Incrementing{}).NewId
	for _, path := range paths {
		feats, err := parsePath(path, newIDFunc)

		switch {
		case os.IsNotExist(err):
			return nil, fmt.Errorf(`feature path "%s" is not available`, path)
		case os.IsPermission(err):
			return nil, fmt.Errorf(`feature path "%s" is not accessible`, path)
		case err != nil:
			return nil, err
		}

		for _, ft := range feats {
			if _, duplicate := uniqueFeatureURI[ft.Uri]; duplicate {
				continue
			}

			uniqueFeatureURI[ft.Uri] = ft
			featureIdxs[ft.Uri] = order

			order++
		}
	}

	var features = make([]*models.Feature, len(uniqueFeatureURI))
	for uri, feature := range uniqueFeatureURI {
		idx := featureIdxs[uri]
		features[idx] = feature
	}

	features = filterFeatures(filter, features)

	return features, nil
}

func filterFeatures(filter string, features []*models.Feature) (result []*models.Feature) {
	for _, ft := range features {
		ft.Pickles = tags.ApplyTagFilter(filter, ft.Pickles)

		if ft.Feature != nil && len(ft.Pickles) > 0 {
			result = append(result, ft)
		}
	}

	return
}
