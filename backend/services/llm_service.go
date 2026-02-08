package services

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GenerateAnswer calls Google Gemini to generate a natural answer based on context.
func GenerateAnswer(matches []ChunkMatch, question string) (string, error) {
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

	// 3. Construct Prompt with source information
	contextBlock := ""
	for i, match := range matches {
		sourceInfo := ""
		if match.Source != "" {
			sourceInfo = fmt.Sprintf(" (from %s", match.Source)
			if match.Page > 0 {
				sourceInfo += fmt.Sprintf(", Page %d", match.Page)
			}
			sourceInfo += ")"
		}
		contextBlock += fmt.Sprintf("Source %d%s:\n%s\n\n", i+1, sourceInfo, match.Text)
	}

	prompt := fmt.Sprintf(`You are a high-level technical assistant for the BPT Knowledge Center.
Your goal is to provide a comprehensive, clear, and professional answer based **ONLY** on the context provided below.

### Formatting Guidelines:
1. **Direct Summary**: Start with a 1-2 sentence high-level summary of the answer.
2. **Structured Details**: Use bullet points or numbered lists for technical features, steps, or list items.
3. **Emphasis**: Use **bold** text for key terms, categories, or important entities.
4. **Tone**: Maintain a professional, objective, and helpful tone.
5. **Constraints**: 
   - If the information is not in the context, explicitly state: "I couldn't find that specific information in the available documents."
   - Do not mention the context sources (e.g., "Source 1 says...") directly in the narrative unless necessary for clarity.
   - Use standard Markdown formatting for best readability in a web interface.

---
**Context Documents:**
%s

---
**User Question:** %s

**Structured Answer:**`, contextBlock, question)

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
