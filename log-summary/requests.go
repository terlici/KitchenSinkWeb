package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// curl -X POST \
//   -d '{"input": "What is the meaning of life?"}' \
//   -H 'Content-Type: application/json' \
//   https://api.deepinfra.com/v1/inference/meta-llama/Llama-2-7b-chat-hf

func Request(url string, data any) []byte {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(dataJSON))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")

	// add authorization header to the req
	// req.Header.Add("Authorization", "Bearer <token>")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Response status not OK. ", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return body
}

func ChatInference(input string) string {
	data := struct {
		Input string `json:"input"`
	}{input}

	body := Request("https://api.deepinfra.com/v1/inference/meta-llama/Llama-2-7b-chat-hf", data)

	var response struct {
		RequestID        string `json:"request_id"`
		TokensCount      int    `json:"num_tokens"`
		InputTokensCount int    `json:"num_input_tokens"`
		Results          []struct {
			Text string `json:"generated_text"`
		} `json:"results"`
		Status struct {
			Status          string  `json:"status"`
			Duration        int64   `json:"runtime_ms"`
			Cost            float64 `json:"cost"`
			TokensGenerated int     `json:"tokens_generated"`
			TokensInput     int     `json:"tokens_input"`
		} `json:"inference_status"`
	}

	err := json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err)
	}

	if response.Status.Status != "succeeded" {
		log.Println(response)
		log.Fatal("Inference status not succeeded. ", response.Status.Status)
	} else if len(response.Results) == 0 {
		log.Println(response)
		log.Fatal("No results generated.")
	}

	return response.Results[0].Text
}

func ImageInference(prompt string) string {
	data := struct {
		Prompt string `json:"prompt"`
	}{prompt}

	body := Request("https://api.deepinfra.com/v1/inference/runwayml/stable-diffusion-v1-5", data)

	var response struct {
		RequestID           string   `json:"request_id"`
		Images              []string `json:"images"`
		NSFWContentDetected []bool   `json:"nsfw_content_detected"`
		Seed                int64    `json:"seed"`
		Status              struct {
			Status   string  `json:"status"`
			Duration int64   `json:"runtime_ms"`
			Cost     float64 `json:"cost"`
		} `json:"inference_status"`
	}

	err := json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err)
	}

	if response.Status.Status != "succeeded" {
		log.Println(response)
		log.Fatal("Inference status not succeeded. ", response.Status.Status)
	} else if len(response.Images) == 0 {
		log.Println(response)
		log.Fatal("No images generated.")
	}

	return response.Images[0]
}

func Summarize(entries []string) string {
	var input string
	input += "Write a warm summary of my day based only on the following events that happened to me:\n"
	for _, entry := range entries {
		input += "- " + entry + "\n"
	}

	return ChatInference(input)
}

func SummaryImage(summary string) string {
	return ImageInference(summary + ", photorealistic, fantasy")
}
