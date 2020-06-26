package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/letsboot/core/kubernetes-course/solution/code/core/internal/model"
)

const (
	ApiUrl = "http://backend"
)

func PageCallback(page model.Page, urls []string) error {
	bs, err := json.Marshal(&urls)
	if err != nil {
		return err
	}
	r := bytes.NewReader(bs)
	_, err = http.Post(fmt.Sprintf("%s/callback/%d", ApiUrl, page.Id), "application/json", r)
	return err
}
