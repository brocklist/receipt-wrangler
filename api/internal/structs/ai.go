package structs

import "context"

type AiClientMessage struct {
	Role    string   `json:"role"`
	Content string   `json:"content"`
	Images  []string `json:"images"`
}

type AiChatCompletionOptions struct {
	// Optional worker context. Existing synchronous callers retain their behavior.
	Context context.Context `json:"-"`
	// Messages to send to the AI model
	Messages []AiClientMessage `json:"messages"`

	// Determines whether to decrypt the key
	DecryptKey bool `json:"decryptKey"`
}

func (options AiChatCompletionOptions) RequestContext() context.Context {
	if options.Context != nil {
		return options.Context
	}
	return context.Background()
}
