package httputil

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/ghetzel/go-stockutil/maputil"
	"github.com/ghetzel/go-stockutil/typeutil"
)

type EncoderFunc func(any) (io.Reader, error)
type DecoderFunc func(io.Reader, any) error
type ErrorDecoderFunc func(*http.Response) error
type InitRequestFunc func() error
type InterceptRequestFunc func(*http.Request) (any, error)
type InterceptResponseFunc func(*http.Response, any) error

var DefaultMultipartFormFileField = `filename`

func JSONEncoder(in any) (io.Reader, error) {
	if req, ok := in.(*http.Request); ok {
		req.Header.Set(`Content-Type`, `application/json`)

		return nil, nil
	} else if data, err := json.Marshal(in); err == nil {
		return bytes.NewBuffer(data), nil
	} else {
		return nil, err
	}
}

func JSONDecoder(in io.Reader, out any) error {
	return json.NewDecoder(in).Decode(out)
}

func XMLDecoder(in io.Reader, out any) error {
	return xml.NewDecoder(in).Decode(out)
}

type MultipartFormFile struct {
	Filename string    `json:"filename"`
	Data     io.Reader `json:"data"`
}

type multipartFormRequest struct {
	data        io.Reader
	boundary    string
	contentType string
}

func (self *multipartFormRequest) Read(b []byte) (int, error) {
	return self.data.Read(b)
}

// Specifies that the given data should be encoded as a multipart/form-data request.
func MultipartFormEncoder(in any) (io.Reader, error) {
	if req, ok := in.(*http.Request); ok {
		req.Header.Set(`Content-Type`, `multipart/form-data`)

		return nil, nil
	} else {
		var output = bytes.NewBuffer(nil)
		var mp = multipart.NewWriter(output)
		var fields = make(map[string]any)

		if typeutil.IsMap(in) {
			fields = maputil.M(in).MapNative()
		} else if inR, ok := in.(io.Reader); ok {
			fields[DefaultMultipartFormFileField] = inR
		}

		// encode all fields as multipart form parts
		for field, value := range fields {
			var filename string

			// make []byte into a reader
			if vBytes, ok := value.([]byte); ok {
				value = bytes.NewBuffer(vBytes)
			}

			// get the filename from os.File or MultipartFormFile (if provided)
			if vFile, ok := value.(*os.File); ok {
				filename = vFile.Name()
			} else {
				var vMPFF *MultipartFormFile

				if v, ok := value.(MultipartFormFile); ok {
					vMPFF = &v
				} else if v, ok := value.(*MultipartFormFile); ok {
					vMPFF = v
				}

				if vMPFF != nil {
					filename = vMPFF.Filename

					if vMPFF.Data != nil {
						value = vMPFF.Data
					} else if f, err := os.Open(filename); err == nil {
						defer f.Close()
						value = f
					} else {
						return nil, fmt.Errorf("Cannot add multipart file %q: %v", filename, err)
					}
				}
			}

			if valueR, ok := value.(io.Reader); ok {
				// if a filename is present, then create a file part
				if filename != `` {
					if outW, err := mp.CreateFormFile(field, filename); err == nil {
						if _, err := io.Copy(outW, valueR); err != nil {
							return nil, fmt.Errorf("Cannot write multipart form field %q: %v", field, err)
						}
					} else {
						return nil, fmt.Errorf("Cannot create multipart form field %q: %v", field, err)
					}
				} else if outW, err := mp.CreateFormField(field); err == nil {
					// readers are handled as straight binary copies
					if _, err := io.Copy(outW, valueR); err != nil {
						return nil, fmt.Errorf("Cannot write multipart form field %q: %v", field, err)
					}
				} else {
					return nil, fmt.Errorf("Cannot create multipart form field %q: %v", field, err)
				}
			} else if err := mp.WriteField(field, typeutil.String(value)); err != nil {
				// everything else is stringified
				return nil, fmt.Errorf("Cannot encode multipart form field %q: %v", field, err)
			}
		}

		return &multipartFormRequest{
			data:        output,
			boundary:    mp.Boundary(),
			contentType: mp.FormDataContentType(),
		}, mp.Close()
	}
}
