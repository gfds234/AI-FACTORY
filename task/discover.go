package task

import (
	"ai-studio/orchestrator/llm"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

// DiscoverQuestions are the default/fallback validation questions
var DiscoverQuestions = []string{
	"What makes it different from existing solutions?",
	"Who is your target audience?",
	"How would you monetize this?",
}

// DiscoverSession represents an active discovery session
type DiscoverSession struct {
	ID        string    `json:"id"`
	RawIdea   string    `json:"raw_idea"`
	Questions []string  `json:"questions"`   // Dynamic AI-generated questions
	IdeaType  string    `json:"idea_type"`   // Detected category (game, app, saas, etc.)
	Answers   []string  `json:"answers"`
	Verdict   string    `json:"verdict"`     // "GO", "REFINE", "PASS", ""
	Reasoning string    `json:"reasoning"`
	Status    string    `json:"status"`      // "answering", "complete"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ResearchModule generates tailored discovery questions
type ResearchModule struct {
	client *llm.Client
}

// NewResearchModule creates a new research module
func NewResearchModule(client *llm.Client) *ResearchModule {
	return &ResearchModule{client: client}
}

// GenerateQuestions analyzes the idea and generates category-specific questions
func (rm *ResearchModule) GenerateQuestions(rawIdea string) ([]string, string, error) {
	prompt := rm.buildQuestionPrompt(rawIdea)

	log.Printf("Generating dynamic questions for idea: %s", rawIdea[:min(len(rawIdea), 50)]+"...")

	// Use Mistral model for question generation (fast + free)
	response, err := rm.client.Generate("mistral:7b-instruct-v0.2-q4_K_M", prompt)
	if err != nil {
		log.Printf("Question generation failed, using defaults: %v", err)
		return getDefaultQuestions(), "unknown", nil
	}

	questions, category, err := rm.parseQuestionResponse(response)
	if err != nil {
		log.Printf("Question parsing failed, using defaults: %v", err)
		return getDefaultQuestions(), "unknown", nil
	}

	log.Printf("Generated %d questions for category: %s", len(questions), category)
	return questions, category, nil
}

// buildQuestionPrompt creates the LLM prompt for question generation
func (rm *ResearchModule) buildQuestionPrompt(rawIdea string) string {
	return fmt.Sprintf(`You are an expert product researcher and startup advisor. Analyze this product idea and generate smart discovery questions.

Product Idea:
%s

Your task:
1. Identify the idea category (game, mobile app, web app, SaaS, hardware, service, marketplace, AI tool, etc.)
2. Generate exactly 4 discovery questions that validate this specific type of product

Questions should:
- Be open-ended (not yes/no)
- Probe different critical aspects (market fit, differentiation, business model, feasibility)
- Be answerable by the person with the idea (not require market research)
- Surface potential risks specific to this category
- Help determine GO/REFINE/PASS verdict

Format your response EXACTLY like this:

CATEGORY: [one word category]

QUESTIONS:
1. [Question about market/need]
2. [Question about differentiation/competition]
3. [Question about business model/monetization]
4. [Question about execution/feasibility]

Be specific to the category. For games ask about retention hooks, for SaaS ask about customer acquisition, for hardware ask about manufacturing, etc.`, rawIdea)
}

// parseQuestionResponse extracts category and questions from LLM response
func (rm *ResearchModule) parseQuestionResponse(response string) ([]string, string, error) {
	lines := strings.Split(response, "\n")
	var questions []string
	var category string

	inQuestions := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Parse category
		if strings.HasPrefix(line, "CATEGORY:") {
			category = strings.TrimSpace(strings.TrimPrefix(line, "CATEGORY:"))
			category = strings.ToLower(category)
			continue
		}

		// Start of questions section
		if strings.HasPrefix(line, "QUESTIONS:") {
			inQuestions = true
			continue
		}

		// Parse numbered questions
		if inQuestions && line != "" {
			// Remove number prefix (1., 2., etc.)
			if len(line) > 2 && line[0] >= '1' && line[0] <= '9' && line[1] == '.' {
				question := strings.TrimSpace(line[2:])
				if question != "" {
					questions = append(questions, question)
				}
			}
		}
	}

	// Validation
	if len(questions) < 3 {
		return nil, "", fmt.Errorf("insufficient questions generated: got %d, need at least 3", len(questions))
	}

	if category == "" {
		category = "unknown"
	}

	// Return 3-4 questions max
	if len(questions) > 4 {
		questions = questions[:4]
	}

	return questions, category, nil
}

// getDefaultQuestions returns hardcoded fallback questions
func getDefaultQuestions() []string {
	return []string{
		"What makes it different from existing solutions?",
		"Who is your target audience?",
		"How would you monetize this?",
	}
}

// DiscoverManager handles discovery sessions
type DiscoverManager struct {
	sessions       map[string]*DiscoverSession
	sessionsMux    sync.RWMutex
	client         *llm.Client
	researchModule *ResearchModule
}

// NewDiscoverManager creates a new discovery manager
func NewDiscoverManager(client *llm.Client) *DiscoverManager {
	return &DiscoverManager{
		sessions:       make(map[string]*DiscoverSession),
		client:         client,
		researchModule: NewResearchModule(client),
	}
}

// StartSession creates a new discovery session with dynamic questions
func (dm *DiscoverManager) StartSession(rawIdea string) *DiscoverSession {
	dm.sessionsMux.Lock()
	defer dm.sessionsMux.Unlock()

	// Generate tailored questions based on idea type
	questions, ideaType, err := dm.researchModule.GenerateQuestions(rawIdea)
	if err != nil {
		// Use default questions on error
		questions = getDefaultQuestions()
		ideaType = "unknown"
	}

	session := &DiscoverSession{
		ID:        fmt.Sprintf("discover_%d", time.Now().Unix()),
		RawIdea:   rawIdea,
		Questions: questions,
		IdeaType:  ideaType,
		Answers:   make([]string, 0, len(questions)),
		Status:    "answering",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	dm.sessions[session.ID] = session
	return session
}

// AddAnswer adds an answer to the session
func (dm *DiscoverManager) AddAnswer(sessionID, answer string) (*DiscoverSession, error) {
	dm.sessionsMux.Lock()
	defer dm.sessionsMux.Unlock()

	session := dm.sessions[sessionID]
	if session == nil {
		return nil, fmt.Errorf("session not found")
	}

	if session.Status == "complete" {
		return nil, fmt.Errorf("session already complete")
	}

	if len(session.Answers) >= len(session.Questions) {
		return nil, fmt.Errorf("all questions already answered")
	}

	session.Answers = append(session.Answers, answer)
	session.UpdatedAt = time.Now()

	// Check if all questions answered
	if len(session.Answers) >= len(session.Questions) {
		dm.scoreSession(session)
	}

	return session, nil
}

// GetSession retrieves a session by ID
func (dm *DiscoverManager) GetSession(sessionID string) *DiscoverSession {
	dm.sessionsMux.RLock()
	defer dm.sessionsMux.RUnlock()
	return dm.sessions[sessionID]
}

// GetAllSessions returns all discovery sessions
func (dm *DiscoverManager) GetAllSessions() []*DiscoverSession {
	dm.sessionsMux.RLock()
	defer dm.sessionsMux.RUnlock()

	sessions := make([]*DiscoverSession, 0, len(dm.sessions))
	for _, session := range dm.sessions {
		sessions = append(sessions, session)
	}
	return sessions
}

// scoreSession determines verdict based on answers
func (dm *DiscoverManager) scoreSession(session *DiscoverSession) {
	session.Status = "complete"

	// Try LLM-based scoring first (primary)
	verdict, reasoning := dm.scoreWithLLM(session)

	session.Verdict = verdict
	session.Reasoning = reasoning
}

// scoreWithLLM uses LLM to intelligently score answers
func (dm *DiscoverManager) scoreWithLLM(session *DiscoverSession) (string, string) {
	prompt := dm.buildScoringPrompt(session)

	// Use mistral model for scoring (same as validate task)
	response, err := dm.client.Generate("mistral:7b-instruct-v0.2-q4_K_M", prompt)
	if err != nil {
		// Fallback to basic scoring
		return dm.scoreBasic(session)
	}

	// Parse LLM response for verdict
	verdict, reasoning := dm.parseVerdictResponse(response)

	// If parsing failed, use fallback
	if verdict == "" {
		return dm.scoreBasic(session)
	}

	return verdict, reasoning
}

// scoreBasic provides fallback rule-based scoring
func (dm *DiscoverManager) scoreBasic(session *DiscoverSession) (string, string) {
	score := 0

	// Check answer quality by length and content
	for _, ans := range session.Answers {
		words := len(strings.Fields(ans))
		if words > 20 {
			score += 2
		} else if words > 10 {
			score += 1
		}
	}

	if score >= 5 {
		return "GO", "Strong answers across all questions. Ready to proceed with development."
	} else if score >= 3 {
		return "REFINE", "Some answers need more detail. Consider expanding your responses with more specifics."
	} else {
		return "PASS", "Answers lack sufficient detail. More thought needed before proceeding."
	}
}

// buildScoringPrompt creates the LLM prompt for scoring (now dynamic)
func (dm *DiscoverManager) buildScoringPrompt(session *DiscoverSession) string {
	// Build Q&A pairs dynamically
	var qaSection strings.Builder
	for i, question := range session.Questions {
		qaSection.WriteString(fmt.Sprintf("\nQuestion %d: %s\n", i+1, question))
		qaSection.WriteString(fmt.Sprintf("Answer: %s\n", session.Answers[i]))
	}

	categoryContext := ""
	if session.IdeaType != "unknown" && session.IdeaType != "" {
		categoryContext = fmt.Sprintf("\nIdea Category: %s\n", session.IdeaType)
	}

	return fmt.Sprintf(`You are evaluating a product idea. The user has provided answers to validation questions.

Product Idea: %s%s

%s

Based on these answers, provide a verdict:
- GO: Strong, well-thought-out answers. Ready for development.
- REFINE: Decent answers but need more detail or clarity.
- PASS: Weak answers, idea needs more thought.

Respond in this exact format:
VERDICT: [GO|REFINE|PASS]
REASONING: [2-3 sentence explanation]`,
		session.RawIdea,
		categoryContext,
		qaSection.String())
}

// parseVerdictResponse extracts verdict and reasoning from LLM response
func (dm *DiscoverManager) parseVerdictResponse(response string) (string, string) {
	lines := strings.Split(response, "\n")
	verdict := ""
	reasoning := ""

	for i, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "VERDICT:") {
			verdict = strings.TrimSpace(strings.TrimPrefix(line, "VERDICT:"))
			// Normalize verdict to uppercase
			verdict = strings.ToUpper(verdict)

			// Validate verdict
			if verdict != "GO" && verdict != "REFINE" && verdict != "PASS" {
				verdict = "REFINE" // Default fallback
			}
		} else if strings.HasPrefix(line, "REASONING:") {
			// Collect reasoning (may span multiple lines)
			reasoning = strings.TrimSpace(strings.TrimPrefix(line, "REASONING:"))
			// Check if reasoning continues on next lines
			for j := i + 1; j < len(lines); j++ {
				nextLine := strings.TrimSpace(lines[j])
				if nextLine != "" && !strings.Contains(nextLine, ":") {
					reasoning += " " + nextLine
				} else {
					break
				}
			}
		}
	}

	// If no reasoning found, provide default
	if reasoning == "" {
		switch verdict {
		case "GO":
			reasoning = "The answers demonstrate a clear understanding of the product vision and market opportunity."
		case "REFINE":
			reasoning = "The answers show potential but could benefit from more specific details and deeper analysis."
		case "PASS":
			reasoning = "The answers need significant improvement before moving forward with implementation."
		default:
			reasoning = "Unable to fully evaluate the responses."
		}
	}

	return verdict, reasoning
}

// GetCurrentQuestion returns the next question to ask based on session state
func (dm *DiscoverManager) GetCurrentQuestion(session *DiscoverSession) (int, string) {
	questionNum := len(session.Answers)

	// Use session-specific questions (dynamic or default)
	if questionNum >= len(session.Questions) {
		return 0, "" // All questions answered
	}

	return questionNum + 1, session.Questions[questionNum]
}

// min returns the minimum of two ints
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
