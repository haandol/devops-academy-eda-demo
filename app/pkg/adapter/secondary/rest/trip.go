package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/haandol/devops-academy-eda-demo/pkg/config"
	"github.com/haandol/devops-academy-eda-demo/pkg/constant"
	"github.com/haandol/devops-academy-eda-demo/pkg/util"
)

type TripRestAdapter struct {
	AuthHeader string
	Host       string
}

func NewTripRestAdapter(cfg *config.Config) *TripRestAdapter {
	return &TripRestAdapter{
		AuthHeader: cfg.App.AuthHeader,
		Host:       cfg.Rest.HotelHost,
	}
}

func (a *TripRestAdapter) InjectError(ctx context.Context, flag bool) error {
	logger := util.GetLogger().With(
		"module", "TripRestAdapter",
		"func", "InjectError",
	)

	buf, err := json.Marshal(struct {
		Flag bool `json:"flag"`
	}{
		Flag: flag,
	})
	if err != nil {
		return err
	}

	errorInjectionURL := fmt.Sprintf("%s/v1/hotels/error/", a.Host)
	req, err := http.NewRequestWithContext(ctx,
		"PUT", errorInjectionURL, bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	req.Header.Add(constant.AuthHeaderKey, a.AuthHeader)
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{Timeout: time.Duration(30) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to inject error")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var data struct {
		Success bool `json:"success"`
		Data    bool `json:"data"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (a *TripRestAdapter) GetInjectionStatus(ctx context.Context) (bool, error) {
	logger := util.GetLogger().With(
		"module", "TripRestAdapter",
		"func", "GetInjectionStatus",
	)
	errorInjectionURL := fmt.Sprintf("%s/v1/hotels/error/", a.Host)
	req, err := http.NewRequestWithContext(ctx, "GET", errorInjectionURL, http.NoBody)
	if err != nil {
		return false, err
	}
	req.Header.Add(constant.AuthHeaderKey, a.AuthHeader)
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{Timeout: time.Duration(30) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	logger.Infow("response body", "body", string(body))

	var data struct {
		Success bool `json:"success"`
		Data    bool `json:"data"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		logger.Error(err.Error())
		return false, err
	}

	return data.Data, nil
}
