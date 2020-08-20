package utils

import (
	"../../api/responses"
	"io/ioutil"
	. "net/http"
)

func GetBodyFromRequest(w ResponseWriter, r *Request) []byte {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, StatusUnprocessableEntity, err)
		return nil
	}

	return body
}
