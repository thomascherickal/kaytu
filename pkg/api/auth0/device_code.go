package auth0

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type DeviceCodeRequest struct {
	ClientId string `json:"client_id"`
	Scope    string `json:"scope"`
	Audience string `json:"audience"`
}

type DeviceCodeResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationUrl         string `json:"verification_uri"`
	VerificationUrlComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

func RequestDeviceCode() (string, error) {
	payload := DeviceCodeRequest{
		ClientId: Auth0ClientID,
		Scope:    "openid profil email api:read",
		Audience: "https://app.kaytu.io",
	}
	payloadEncode, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("[requestDeviceCode] : %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/oauth/device/code", Auth0Hostname), bytes.NewBuffer(payloadEncode))
	if err != nil {
		return "", fmt.Errorf("[requestDeviceCode] : %v", err)
	}
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("[requestDeviceCode] : %v", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("[requestDeviceCode] : %v", err)
	}
	err = res.Body.Close()
	if err != nil {
		return "", fmt.Errorf("[requestDeviceCode] : %v", err)

	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid status code: %d, %s", res.StatusCode, string(body))
	}

	response := DeviceCodeResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("[requestDeviceCode] : %v", err)
	}

	fmt.Println("open this url in your browser:")
	fmt.Println(response.VerificationUrlComplete)

	return response.DeviceCode, nil
}
