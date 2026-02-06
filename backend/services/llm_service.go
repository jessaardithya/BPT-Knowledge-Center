package services

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GenerateAnswer calls Google Gemini to generate a natural answer based on context.
func GenerateAnswer(matches []string, question string) (string, error) {
	// 1. Check for API Key (Try GOOGLE_API_KEY first as per user, then GEMINI_API_KEY)
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
	}

	fmt.Printf("DEBUG: LLM Service called. Key Length: %d\n", len(apiKey))

	if apiKey == "" {
		fmt.Println("DEBUG: No API Key found. Skipping LLM.")
		return "", nil // No key = No generation (Fallback to raw search)
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		fmt.Printf("DEBUG: Client Creation Error: %v\n", err)
		return "", fmt.Errorf("failed to create gemini client: %w", err)
	}
	defer client.Close()

	// 2. Select Model (User requested gemini-3-flash-preview)
	model := client.GenerativeModel("gemini-3-flash-preview")
	model.SetTemperature(0.2) // Low temperature for factual answers

	// 3. Construct Prompt
	contextBlock := ""
	for i, match := range matches {
		contextBlock += fmt.Sprintf("Source %d:\n%s\n\n", i+1, match)
	}

	prompt := fmt.Sprintf(`You are a helpful assistant for the BPT Knowledge Center.
Answer the user's question using ONLY the provided context information below.
If the answer is not in the context, say "I couldn't find that information in the documents."
Do not make up information. Keep the answer concise.

Context:
%s

Question: %s
Answer:`, contextBlock, question)

	// 4. Generate
	fmt.Println("DEBUG: Sending request to Gemini...")
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		fmt.Printf("DEBUG: Generation Error: %v\n", err)
		return "", fmt.Errorf("gemini generation failed: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		fmt.Println("DEBUG: Empty response candidates")
		return "", fmt.Errorf("empty response from gemini")
	}

	// Extract text
	var answer string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			answer += string(txt)
		}
	}

	return answer, nil
}
