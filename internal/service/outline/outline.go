package outline

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"

	"github.com/0x16f/vpn-resolver/internal/entity"
	"github.com/goccy/go-json"
	"github.com/sirupsen/logrus"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) CreateUser(ctx context.Context, req entity.OutlineCreateUserReq) (entity.OutlineCreateUserResp, error) {
	outlineReq := createUser{
		Method:   "chacha20-ietf-poly1305",
		Name:     req.Name,
		Password: req.Password,
	}

	encoded, err := json.Marshal(outlineReq)
	if err != nil {
		return entity.OutlineCreateUserResp{}, err
	}

	var (
		r    *http.Request
		furl = fmt.Sprintf("%s/access-keys", newOutlineURL(req.OutlineInfo))
	)

	r, err = http.NewRequestWithContext(ctx, http.MethodPost, furl, bytes.NewBuffer(encoded))
	if err != nil {
		return entity.OutlineCreateUserResp{}, err
	}

	r.Header.Add("Content-Type", "application/json")

	var (
		resp *http.Response
	)

	client := http.Client{
		Transport: getTransport(),
	}

	resp, err = client.Do(r)
	if err != nil {
		return entity.OutlineCreateUserResp{}, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusCreated {
		bodyRes, err := io.ReadAll(resp.Body)
		if err != nil {
			return entity.OutlineCreateUserResp{}, err
		}

		logrus.Errorf("failed to create user: %s", bodyRes)

		return entity.OutlineCreateUserResp{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var outlineResp entity.OutlineCreateUserResp

	if err = json.NewDecoder(resp.Body).Decode(&outlineResp); err != nil {
		return entity.OutlineCreateUserResp{}, err
	}

	return outlineResp, nil
}

func (s *Service) DeleteUser(ctx context.Context, req entity.OutlineDeleteUserReq) error {
	furl := fmt.Sprintf("%s/access-keys/%s", newOutlineURL(req.OutlineInfo), req.UserID)

	r, err := http.NewRequestWithContext(ctx, http.MethodDelete, furl, nil)
	if err != nil {
		return err
	}

	var (
		resp *http.Response
	)

	client := http.Client{
		Transport: getTransport(),
	}

	resp, err = client.Do(r)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusNoContent {
		bodyRes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		logrus.Errorf("failed to create user: %s", bodyRes)

		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func getTransport() http.RoundTripper {
	var defaultTransport http.RoundTripper = &http.Transport{
		Proxy: nil,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	return defaultTransport
}

func newOutlineURL(info entity.OutlineInfo) string {
	return fmt.Sprintf("https://%s:%d/%s", info.OutlineURL, info.OutlinePort, info.OutlineSecret)
}
