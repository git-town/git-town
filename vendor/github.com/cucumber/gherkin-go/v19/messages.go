package gherkin

import (
	"encoding/json"
	"fmt"
	"github.com/cucumber/messages-go/v16"
	"io"
	"io/ioutil"
	"strings"
)

func Messages(
	paths []string,
	decoder *json.Decoder,
	language string,
	includeSource bool,
	includeGherkinDocument bool,
	includePickles bool,
	encoder *json.Encoder,
	newId func() string,
) ([]messages.Envelope, error) {
	var result []messages.Envelope
	var err error

	handleMessage := func(result []messages.Envelope, message *messages.Envelope) ([]messages.Envelope, error) {
		if encoder != nil {
			err = encoder.Encode(message)
			return result, err
		} else {
			result = append(result, *message)
		}

		return result, err
	}

	processSource := func(source *messages.Source) error {
		if includeSource {
			result, err = handleMessage(result, &messages.Envelope{
				Source: source,
			})
		}
		doc, err := ParseGherkinDocumentForLanguage(strings.NewReader(source.Data), language, newId)
		if errs, ok := err.(parseErrors); ok {
			// expected parse errors
			for _, err := range errs {
				if pe, ok := err.(*parseError); ok {
					result, err = handleMessage(result, pe.asMessage(source.Uri))
				} else {
					return fmt.Errorf("parse feature file: %s, unexpected error: %+v\n", source.Uri, err)
				}
			}
			return nil
		}

		if includeGherkinDocument {
			doc.Uri = source.Uri
			result, err = handleMessage(result, &messages.Envelope{
				GherkinDocument: doc,
			})
		}

		if includePickles {
			for _, pickle := range Pickles(*doc, source.Uri, newId) {
				result, err = handleMessage(result, &messages.Envelope{
					Pickle: pickle,
				})
			}
		}
		return nil
	}

	if len(paths) == 0 {
		for {
			envelope := &messages.Envelope{}
			err := decoder.Decode(envelope)
			//marshal, err := json.Marshal(envelope)
			//fmt.Println(string(marshal))
			if err == io.EOF {
				break
			}

			if envelope.Source != nil {
				err = processSource(envelope.Source)
				if err != nil {
					return result, err
				}
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

func (a *parseError) asMessage(uri string) *messages.Envelope {
	return &messages.Envelope{
		ParseError: &messages.ParseError{
			Message: a.Error(),
			Source: &messages.SourceReference{
				Uri: uri,
				Location: &messages.Location{
					Line:   int64(a.loc.Line),
					Column: int64(a.loc.Column),
				},
			},
		},
	}
}
