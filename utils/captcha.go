package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
)

type VerifyResponse struct {
	Success bool `json:"success"`
}

func VerifyCaptcha(response string) (bool, error) {
	secret := "0x0000000000000000000000000000000000000000"
	if os.Getenv("CAPTCHA_SECRET") != "" {
		secret = os.Getenv("CAPTCHA_SECRET")
	}

	data := url.Values{}
	data.Set("response", response)
	data.Set("secret", secret)

	res, err := http.PostForm("https://hcaptcha.com/siteverify", data)
	if err != nil {
		return false, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	var verifyResponse VerifyResponse
	err = json.Unmarshal(body, &verifyResponse)
	if err != nil {
		return false, err
	}

	return verifyResponse.Success, nil
}