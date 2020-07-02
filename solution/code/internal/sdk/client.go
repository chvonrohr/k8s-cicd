package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"net/http"

	"gitlab.com/letsboot/core/kubernetes-course/solution/code/core/internal/model"
)

func PageCallback(page model.Page, urls []string) error {
	bs, err := json.Marshal(&urls)
	if err != nil {
		return err
	}
	r := bytes.NewReader(bs)
	_, err = http.Post(fmt.Sprintf("%s/callback/%d", viper.GetString("backend.url"), page.ID), "application/json", r)
	return err
}
