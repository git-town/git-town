package gherkin

import (
	"fmt"
	"github.com/cucumber/messages-go/v10"
	gio "github.com/gogo/protobuf/io"
	"io"
	"io/ioutil"
	"math"
	"strings"
)

func Messages(
	paths []string,
	sourceStream io.Reader,
	language string,
	includeSource bool,
	includeGherkinDocument bool,
	includePickles bool,
	writer gio.WriteCloser,
	newId func() string,
) ([]messages.Envelope, error) {
	var result []messages.Envelope
	var err error

	handleMessage := func(result []messages.Envelope, message *messages.Envelope) ([]messages.Envelope, error) {
		if writer != nil {
			err = writer.WriteMsg(message)
			return result, err
		} else {
			result = append(result, *message)
		}

		return result, err
	}

	processSource := func(source *messages.Source) error {
		if includeSource {
			result, err = handleMessage(result, &messages.Envelope{
				Message: &messages.Envelope_Source{
					Source: source,
				},
			})
		}
		doc, err := ParseGherkinDocumentForLanguage(strings.NewReader(source.Data), language, newId)
		if errs, ok := err.(parseErrors); ok {
			// expected parse errors
			for _, err := range errs {
				if pe, ok := err.(*parseError); ok {
					result, err = handleMessage(result, pe.asAttachment(source.Uri))
				} else {
					return fmt.Errorf("parse feature file: %s, unexpected error: %+v\n", source.Uri, err)
				}
			}
			return nil
		}

		if includeGherkinDocument {
			doc.Uri = source.Uri
			result, err = handleMessage(result, &messages.Envelope{
				Message: &messages.Envelope_GherkinDocument{
					GherkinDocument: doc,
				},
			})
		}

		if includePickles {
			for _, pickle := range Pickles(*doc, source.Uri, newId) {
				result, err = handleMessage(result, &messages.Envelope{
					Message: &messages.Envelope_Pickle{
						Pickle: pickle,
					},
				})
			}
		}
		return nil
	}

	if len(paths) == 0 {
		reader := gio.NewDelimitedReader(sourceStream, math.MaxInt32)
		for {
			envelope := &messages.Envelope{}
			err := reader.ReadMsg(envelope)
			if err == io.EOF {
				break
			}

			switch t := envelope.Message.(type) {
			case *messages.Envelope_Source:
				processSource(t.Source)
			}
		}
	} else {
		for _, path := range paths {
			in, err := ioutil.ReadFile(path)
			if err != nil {
				return result, fmt.Errorf("read feature file: %s - %+v", path, err)
			}
			source := &messages.Source{
				Uri:       path,
				Data:      string(in),
				MediaType: "text/x.cucumber.gherkin+plain",
			}
			processSource(source)
		}
	}

	return result, err
}

func (a *parseError) asAttachment(uri string) *messages.Envelope {
	return &messages.Envelope{
		Message: &messages.Envelope_Attachment{
			Attachment: &messages.Attachment{
				Body: &messages.Attachment_Text{
					Text: a.Error(),
				},
				Source: &messages.SourceReference{
					Uri: uri,
					Location: &messages.Location{
						Line:   uint32(a.loc.Line),
						Column: uint32(a.loc.Column),
					},
				},
			}},
	}
}
