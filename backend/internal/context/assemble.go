// Package context assembles retrieved results into prompts for the LLM.
package context

import (
	"fmt"
	"strings"

	"tutorai/backend/internal/retrieval"
)

const systemPrompt = `You are TutorAI, a Magic: The Gathering assistant.
Answer the user's question using only the context provided below.
If the context does not contain enough information to answer fully, say so clearly — do not invent card names, rules, or prices.
Keep your answer concise and focused on what the user asked.`

// maxQueryLen is the maximum number of characters accepted from a user query.
// Queries longer than this are truncated before being included in any prompt,
// limiting the impact of prompt injection attempts.
const maxQueryLen = 500

// SanitizeQuery trims whitespace, enforces the length cap, and strips
// characters that could break prompt structure.
func SanitizeQuery(q string) string {
	q = strings.TrimSpace(q)
	if len(q) > maxQueryLen {
		q = q[:maxQueryLen]
	}
	// Remove null bytes and other C0 control characters except newline and tab,
	// which are legitimate in multi-line card lists.
	var b strings.Builder
	for _, r := range q {
		if r == '\n' || r == '\t' || r >= 0x20 {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// Assemble builds the full prompt to send to the LLM.
// The system instructions always precede retrieved context, and user input is
// always placed last and clearly delimited, limiting prompt injection surface.
func Assemble(query string, results []retrieval.Result) string {
	var sb strings.Builder

	sb.WriteString(systemPrompt)
	sb.WriteString("\n\n")

	if len(results) == 0 {
		sb.WriteString("No relevant context was found for this query.\n\n")
	} else {
		sb.WriteString("--- CONTEXT ---\n")
		for i, r := range results {
			if r.Name != "" {
				fmt.Fprintf(&sb, "[%d] %s\n%s\n\n", i+1, r.Name, r.Text)
			} else {
				fmt.Fprintf(&sb, "[%d] %s\n\n", i+1, r.Text)
			}
		}
		sb.WriteString("--- END CONTEXT ---\n\n")
	}

	sb.WriteString("--- USER QUESTION ---\n")
	sb.WriteString(SanitizeQuery(query))
	sb.WriteString("\n--- END USER QUESTION ---")

	return sb.String()
}
