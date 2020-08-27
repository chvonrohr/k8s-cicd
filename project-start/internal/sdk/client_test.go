package sdk

import (
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"gitlab.com/letsboot/core/kubernetes-course/solution/code/core/internal/model"
	"gopkg.in/h2non/gock.v1"
)

func TestNewClient(t *testing.T) {
	client := NewClient("http://backend.com/")
	if client == nil {
		panic("this should not happen")
	}
}

func TestClient_GetSites(t *testing.T) {
	defer gock.Off()

	gock.New("http://backend.com").
		Get("/sites").Reply(200).JSON([]model.Site{{
		Model: gorm.Model{
			ID: 10,
		},
		Url:      "http://example.com",
		Interval: 10 * 24 * time.Hour,
	}})

	client := NewClient("http://backend.com")
	gock.InterceptClient(client.HttpClient)
	defer gock.RestoreClient(client.HttpClient)

	sites, err := client.GetSites()
	if err != nil {
		t.Fatal(err)
	}

	if len(sites) != 1 {
		t.Fatal("invalid site length")
	}

	if sites[0].Url != "http://example.com" {
		t.Fatal("invalid url")
	}

}
