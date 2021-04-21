package render

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	headerContentType = "Content-Type"

	mimeJSON = "application/json"
)

// WriteJSONwithCode sets http Header Content-Type to JSON mime
// and writes encoded JSON data with custom http status code.
func WriteJSONwithCode(w http.ResponseWriter, v interface{}, code int) error {

	buf, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshal: %s", err)
	}

	w.Header().Set(headerContentType, mimeJSON)
	w.WriteHeader(code)

	if _, err := w.Write(buf); err != nil {
		return fmt.Errorf("write response: %s", err)
	}

	return nil
}

// ReadJSON reads and decodes data from request body, assuming the content is JSON (no check)
func ReadJSON(r *http.Request, v interface{}) error {

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("read response: %s", err)
	}

	if len(buf) == 0 {
		return errors.New("no data to parse")
	}

	if err := json.Unmarshal(buf, v); err != nil {
		return fmt.Errorf("unmarshal: %s", err)
	}

	return nil
}
