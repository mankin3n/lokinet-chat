package storage

import (
	"encoding/json"
	"os"
)

type ChatHistory struct {
	Messages []string
}

func SaveHistory(history ChatHistory) error {
	file, err := os.Create("history.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(history)
}

func LoadHistory() (ChatHistory, error) {
	file, err := os.Open("history.json")
	if err != nil {
		return ChatHistory{}, err
	}
	defer file.Close()

	var history ChatHistory
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&history)
	return history, err
}
