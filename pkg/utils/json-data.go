package utils

import "encoding/json"

func Marshal(v any) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func Unmarshal(data []byte, v any) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	return nil
}
