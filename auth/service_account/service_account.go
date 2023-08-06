package service_account

import (
	"encoding/json"
	"io"
	"os"
)

type ServiceAccount struct {
	ClientEmail  string `json:"client_email"`
	ClientID     string `json:"client_id"`
	PrivateKey   string `json:"private_key"`
	PrivateKeyID string `json:"private_key_id"`
}

func FromFile(path string) (ServiceAccount, error) {

	file, err := os.Open(path)
	if err != nil {
		return ServiceAccount{}, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return ServiceAccount{}, err
	}

	var sa ServiceAccount
	err = json.Unmarshal(bytes, &sa)
	if err != nil {
		return ServiceAccount{}, err
	}

	return sa, nil
}
