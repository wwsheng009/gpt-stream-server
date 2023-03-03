package main

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CompletionRequest is a request for the completions API
type ChatRequest struct {
	Model string `json:"model"`
	// A list of string prompts to use.
	// TODO there are other prompt types here for using token integers that we could add support for.
	Messages []Message `json:"messages"`
	// How many tokens to complete up to. Max of 512
	MaxTokens *int `json:"max_tokens,omitempty"`
	// Sampling temperature to use
	Temperature *float32 `json:"temperature,omitempty"`
	// Alternative to temperature for nucleus sampling
	TopP *float32 `json:"top_p,omitempty"`
	// How many choice to create for each prompt
	N    *int     `json:"n"`
	Stop []string `json:"stop,omitempty"`
	// PresencePenalty number between 0 and 1 that penalizes tokens that have already appeared in the text so far.
	PresencePenalty float32 `json:"presence_penalty"`
	// FrequencyPenalty number between 0 and 1 that penalizes tokens on existing frequency in the text so far.
	FrequencyPenalty float32 `json:"frequency_penalty"`

	// Whether to stream back results or not. Don't set this value in the request yourself
	// as it will be overriden depending on if you use CompletionStream or Completion methods.
	Stream bool `json:"stream,omitempty"`
}

// CompletionResponseChoice is one of the choices returned in the response to the Completions API
type ChatResponseChoice struct {
	Delta struct {
		Role    string `json:"role,omitempty"`
		Content string `json:"content" `
	}
	Index        int    `json:"index"`
	FinishReason string `json:"finish_reason"`
}

// CompletionResponse is the full response from a request to the completions API
type ChatResponse struct {
	ID      string               `json:"id"`
	Object  string               `json:"object"`
	Created int                  `json:"created"`
	Model   string               `json:"model"`
	Choices []ChatResponseChoice `json:"choices"`
}

// EnginesResponse is returned from the Engines API
type EnginesResponse struct {
	Data ChatResponse `json:"data"`
	// Object string       `json:"object"`
}
