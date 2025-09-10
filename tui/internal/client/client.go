package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"

	"github.com/linuxunsw/vote/tui/internal/tui/messages"
	"github.com/spf13/viper"
	"golang.org/x/net/publicsuffix"
)

var (
	ErrServerUnavailable      = errors.New("server unavailable")
	ErrServerInternalError    = errors.New("internal error")
	ErrServerNotAuthenticated = errors.New("not authenticated")
)

// Makes a new client with cookie jar
func NewClient() *http.Client {
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})

	return &http.Client{
		Jar: jar,
	}
}

func generateOTP(client *http.Client, ctx context.Context, zID string) error {
	baseURL := viper.GetString("tui.server")
	url := baseURL + "/api/v1/otp/generate"

	requestBody := generateOTPRequest{
		ZID: zID,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// TODO: make more descriptive
	if resp.StatusCode != http.StatusNoContent {
		unmarshaledResponse := &responseError{}
		err := json.Unmarshal(respBody, unmarshaledResponse)
		if err != nil {
			return err
		}

		return ErrServerInternalError
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

func submitOTP(client *http.Client, ctx context.Context, zID string, otp string) error {
	baseURL := viper.GetString("tui.server")
	reqUrl := baseURL + "/api/v1/otp/submit"

	requestBody := submitOTPRequest{
		ZID: zID,
		OTP: otp,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqUrl, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// TODO: make more descriptive
	if resp.StatusCode != http.StatusNoContent {
		unmarshaledResponse := &responseError{}
		err := json.Unmarshal(respBody, unmarshaledResponse)

		if err != nil {
			return err
		}
		return ErrServerInternalError
	}

	// FIX: REMOVE
	cookie := &http.Cookie{
		Name:  "SESSION",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiaXNzIjoidm90ZS1hcGkiLCJpc0FkbWluIjp0cnVlfQ.dbZl2hECaSFzt697PkpFD1y7KJBZD4dhTWYtEinFgs4",
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return err
	}

	client.Jar.SetCookies(parsedURL, []*http.Cookie{cookie})

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

func submitNomination(client *http.Client, ctx context.Context, nomination messages.Submission) (refCode string, err error) {
	baseURL := viper.GetString("tui.server")
	url := baseURL + "/api/v1/nomination"

	body, err := json.Marshal(nomination)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	respBody, err := io.ReadAll(resp.Body)
	respID := resp.Header.Get("X-Request-ID")
	if err != nil {
		return respID, err
	}

	// TODO: deal with whatever this is
	if resp.StatusCode != http.StatusOK {
		unmarshaledResponse := &responseError{}
		err := json.Unmarshal(respBody, unmarshaledResponse)
		if err != nil {
			return respID, err
		}
		// TODO: change this
		return respID, errors.New(strconv.Itoa(resp.StatusCode))
	}

	unmarshaledResponse := &submitNominationResponse{}
	err = json.Unmarshal(respBody, unmarshaledResponse)
	if err != nil {
		return unmarshaledResponse.ID, err
	}

	err = resp.Body.Close()
	if err != nil {
		return unmarshaledResponse.ID, err
	}

	return unmarshaledResponse.ID, nil
}
