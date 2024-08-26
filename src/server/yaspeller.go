package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type wordCorrection struct {
	Code int      `json:"code"`
	Pos  int      `json:"pos"`
	Row  int      `json:"row"`
	Col  int      `json:"col"`
	Len  int      `json:"len"`
	Word string   `json:"word"`
	S    []string `json:"s"`
}

func validateAndFixErrors(note string) (string, error) {
	corrections, err := getCorrections(note)
	if err != nil {
		return "", err
	}
	var correctedText strings.Builder
	currentPos := 0
	runes := []rune(note)
	for _, correction := range corrections {
		correctedText.WriteString(string(runes[currentPos:correction.Pos]))
		correctedText.WriteString(correction.S[0])
		currentPos = correction.Pos + correction.Len
	}

	if currentPos < len(runes) {
		correctedText.WriteString(string(runes[currentPos:]))
	}
	return correctedText.String(), nil
}

func getCorrections(note string) ([]wordCorrection, error) {
	resp, err := http.PostForm("https://speller.yandex.net/services/spellservice.json/checkText",
		url.Values{"text": {note}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var corrections []wordCorrection
	if err := json.Unmarshal(body, &corrections); err != nil {
		return nil, err
	}
	return corrections, nil
}
