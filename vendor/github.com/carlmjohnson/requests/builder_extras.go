package requests

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// URL creates a new Builder suitable for method chaining.
// It is equivalent to calling BaseURL on an empty Builder.
func URL(baseurl string) *Builder {
	return (&Builder{}).BaseURL(baseurl)
}

// New creates a new Builder suitable for method chaining by applying the specified Configs.
// It is equivalent to calling Config on an empty Builder.
// The zero value of Builder is usable,
// so it is not necessary to call New
// when you do not have any Configs to apply.
func New(cfgs ...Config) *Builder {
	return (&Builder{}).Config(cfgs...)
}

// Head sets HTTP method to HEAD.
func (rb *Builder) Head() *Builder {
	return rb.Method(http.MethodHead)
}

// Post sets HTTP method to POST.
//
// Note that setting a Body causes a request to be POST by default.
func (rb *Builder) Post() *Builder {
	return rb.Method(http.MethodPost)
}

// Put sets HTTP method to PUT.
func (rb *Builder) Put() *Builder {
	return rb.Method(http.MethodPut)
}

// Patch sets HTTP method to PATCH.
func (rb *Builder) Patch() *Builder {
	return rb.Method(http.MethodPatch)
}

// Delete sets HTTP method to DELETE.
func (rb *Builder) Delete() *Builder {
	return rb.Method(http.MethodDelete)
}

// Hostf calls Host with fmt.Sprintf.
func (rb *Builder) Hostf(format string, a ...any) *Builder {
	return rb.Host(fmt.Sprintf(format, a...))
}

// Pathf calls Path with fmt.Sprintf.
//
// Note that for security reasons, you must not use %s
// with a user provided string!
func (rb *Builder) Pathf(format string, a ...any) *Builder {
	return rb.Path(fmt.Sprintf(format, a...))
}

// ParamInt converts value to a string and calls Param.
func (rb *Builder) ParamInt(key string, value int) *Builder {
	return rb.Param(key, strconv.Itoa(value))
}

// Params calls Param with all the members of m.
func (rb *Builder) Params(m map[string][]string) *Builder {
	for k, vv := range m {
		rb.Param(k, vv...)
	}
	return rb
}

// Headers calls Header with all the members of m.
func (rb *Builder) Headers(m map[string][]string) *Builder {
	for k, vv := range m {
		rb.Header(k, vv...)
	}
	return rb
}

// Accept sets the Accept header for a request.
func (rb *Builder) Accept(contentTypes string) *Builder {
	return rb.Header("Accept", contentTypes)
}

// CacheControl sets the client-side Cache-Control directive for a request.
func (rb *Builder) CacheControl(directive string) *Builder {
	return rb.Header("Cache-Control", directive)
}

// ContentType sets the Content-Type header on a request.
func (rb *Builder) ContentType(ct string) *Builder {
	return rb.Header("Content-Type", ct)
}

// UserAgent sets the User-Agent header.
func (rb *Builder) UserAgent(s string) *Builder {
	return rb.Header("User-Agent", s)
}

// BasicAuth sets the Authorization header to a basic auth credential.
func (rb *Builder) BasicAuth(username, password string) *Builder {
	auth := username + ":" + password
	v := base64.StdEncoding.EncodeToString([]byte(auth))
	return rb.Header("Authorization", "Basic "+v)
}

// Bearer sets the Authorization header to a bearer token.
func (rb *Builder) Bearer(token string) *Builder {
	return rb.Header("Authorization", "Bearer "+token)
}

// BodyReader sets the Builder's request body to r.
func (rb *Builder) BodyReader(r io.Reader) *Builder {
	return rb.Body(BodyReader(r))
}

// BodyWriter pipes writes from w to the Builder's request body.
func (rb *Builder) BodyWriter(f func(w io.Writer) error) *Builder {
	return rb.Body(BodyWriter(f))
}

// BodyBytes sets the Builder's request body to b.
func (rb *Builder) BodyBytes(b []byte) *Builder {
	return rb.Body(BodyBytes(b))
}

// BodySerializer sets the Builder's request body
// to the serialized object.
func (rb *Builder) BodySerializer(s Serializer, v any) *Builder {
	return rb.
		Body(BodySerializer(s, v))
}

// BodyJSON sets the Builder's request body to the marshaled JSON.
// It uses [JSONSerializer] to marshal the object.
// It also sets ContentType to "application/json"
// if it is not otherwise set.
func (rb *Builder) BodyJSON(v any) *Builder {
	return rb.
		Body(BodyJSON(v)).
		HeaderOptional("Content-Type", "application/json")
}

// BodyForm sets the Builder's request body to the encoded form.
// It also sets the ContentType to "application/x-www-form-urlencoded".
func (rb *Builder) BodyForm(data url.Values) *Builder {
	return rb.
		Body(BodyForm(data)).
		ContentType("application/x-www-form-urlencoded")
}

// BodyFile sets the Builder's request body to read from the given file path.
func (rb *Builder) BodyFile(name string) *Builder {
	return rb.Body(BodyFile(name))
}

// CheckStatus adds a validator for status code of a response.
func (rb *Builder) CheckStatus(acceptStatuses ...int) *Builder {
	return rb.AddValidator(CheckStatus(acceptStatuses...))
}

// CheckContentType adds a validator for the content type header of a response.
func (rb *Builder) CheckContentType(cts ...string) *Builder {
	return rb.AddValidator(CheckContentType(cts...))
}

// CheckPeek adds a validator that peeks at the first n bytes of a response body.
func (rb *Builder) CheckPeek(n int, f func([]byte) error) *Builder {
	return rb.AddValidator(CheckPeek(n, f))
}

// ToDeserializer sets the Builder to decode a response into v
// using a [Deserializer].
func (rb *Builder) ToDeserializer(d Deserializer, v any) *Builder {
	return rb.
		Handle(ToDeserializer(d, v))
}

// ToJSON sets the Builder to decode a response as a JSON object.
//
// It uses [JSONDeserializer] to unmarshal the object.
func (rb *Builder) ToJSON(v any) *Builder {
	return rb.Handle(ToJSON(v))
}

// ToString sets the Builder to write the response body to the provided string pointer.
func (rb *Builder) ToString(sp *string) *Builder {
	return rb.Handle(ToString(sp))
}

// ToBytesBuffer sets the Builder to write the response body to the provided bytes.Buffer.
func (rb *Builder) ToBytesBuffer(buf *bytes.Buffer) *Builder {
	return rb.Handle(ToBytesBuffer(buf))
}

// ToWriter sets the Builder to copy the response body into w.
func (rb *Builder) ToWriter(w io.Writer) *Builder {
	return rb.Handle(ToWriter(w))
}

// ToFile sets the Builder to write the response body to the given file name.
// The file and its parent directories are created automatically.
// For more advanced use cases, use ToWriter.
func (rb *Builder) ToFile(name string) *Builder {
	return rb.Handle(ToFile(name))
}

// CopyHeaders adds a validator which copies the response headers to h.
// Note that because CopyHeaders adds a validator,
// the DefaultValidator is disabled and must be added back manually
// if status code validation is desired.
func (rb *Builder) CopyHeaders(h map[string][]string) *Builder {
	return rb.
		AddValidator(CopyHeaders(h))
}

// ToHeaders sets the method to HEAD and adds a handler which copies the response headers to h.
// To just copy headers, see Builder.CopyHeaders.
func (rb *Builder) ToHeaders(h map[string][]string) *Builder {
	return rb.
		Head().
		Handle(ChainHandlers(CopyHeaders(h), consumeBody))
}

// ErrorJSON adds a validator that applies DefaultValidator
// and decodes the response as a JSON object
// if the DefaultValidator check fails.
func (rb *Builder) ErrorJSON(v any) *Builder {
	return rb.AddValidator(ErrorJSON(v))
}
