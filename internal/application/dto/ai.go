package dto

type AIPayloadDto struct {
	Model             string      `json:"model"`
	Temperature       float64     `json:"temperature,omitempty"`
	TopP              float64     `json:"top_p,omitempty"`
	TopK              uint64      `json:"top_k,omitempty"`
	FrequencyPenalty  float64     `json:"frequency_penalty,omitempty"`
	PresencePenalty   float64     `json:"presence_penalty,omitempty"`
	RepetitionPenalty float64     `json:"repetition_penalty,omitempty"`
	MinP              float64     `json:"min_p,omitempty"`
	TopA              float64     `json:"top_a,omitempty"`
	Messages          []AIMessage `json:"messages"`
}

type AIResponseDto struct {
	Id                string      `json:"id"`
	Provider          string      `json:"provider"`
	Model             string      `json:"model"`
	Object            string      `json:"object"`
	Created           int64       `json:"created"`
	Choices           []AIChoice  `json:"choices"`
	SystemFingerprint interface{} `json:"system_fingerprint"`
	Usage             APIUsage    `json:"usage"`
}

type AIMessage struct {
	Role    string      `json:"role"`
	Content string      `json:"content,omitempty"`
	Refusal interface{} `json:"refusal,omitempty"`
}

type APIUsage struct {
	PromptTokens     int64 `json:"prompt_tokens"`
	CompletionTokens int64 `json:"completion_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}

type AIChoice struct {
	Logprobs           interface{} `json:"logprobs"`
	FinishReason       string      `json:"finish_reason"`
	NativeFinishReason string      `json:"native_finish_reason"`
	Index              int64       `json:"index"`
	Message            AIMessage   `json:"message"`
}
