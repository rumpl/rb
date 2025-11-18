package agent

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"strings"

	"github.com/rumpl/rb/pkg/model/provider"
	"github.com/rumpl/rb/pkg/tools"
)

// Agent represents an AI agent
type Agent struct {
	name                string
	description         string
	welcomeMessage      string
	instruction         string
	expandedInstruction string // cached expanded instruction
	instructionExpanded bool   // flag to track if expansion was done
	expansionError      error  // any error during expansion
	toolsets            []*StartableToolSet
	models              []provider.Provider
	subAgents           []*Agent
	handoffs            []*Agent
	parents             []*Agent
	addDate             bool
	addEnvironmentInfo  bool
	maxIterations       int
	numHistoryItems     int
	addPromptFiles      []string
	tools               []tools.Tool
	commands            map[string]string
	pendingWarnings     []string
}

// New creates a new agent
func New(name, prompt string, opts ...Opt) *Agent {
	agent := &Agent{
		name:        name,
		instruction: prompt,
	}

	for _, opt := range opts {
		opt(agent)
	}

	return agent
}

func (a *Agent) Name() string {
	return a.name
}

// Instruction returns the agent's instructions, expanding any JavaScript template expressions.
// Template expressions use ${tool_name({args})} syntax.
// The expansion is performed once and cached for subsequent calls.
func (a *Agent) Instruction(ctx context.Context) string {
	// Return cached result if already expanded
	if a.instructionExpanded {
		return a.expandedInstruction
	}

	// Mark as expanded to avoid re-expansion
	a.instructionExpanded = true
	a.expandedInstruction = a.instruction

	// Check if expansion is needed (fast path for static instructions)
	if !strings.Contains(a.instruction, "${") {
		return a.expandedInstruction
	}

	// Ensure toolsets are started before expansion
	a.ensureToolSetsAreStarted(ctx)

	// Collect all started toolsets
	var toolSets []tools.ToolSet
	for _, ts := range a.toolsets {
		if ts.started.Load() {
			toolSets = append(toolSets, ts.ToolSet)
		}
	}

	// Expand instruction template
	expanded, err := ExpandInstruction(ctx, a.instruction, toolSets)
	if err != nil {
		// Log warning but continue with original instruction
		slog.Warn("Failed to expand instruction", "agent", a.Name(), "error", err)
		a.expansionError = err
		a.addToolWarning(fmt.Sprintf("instruction expansion failed: %v", err))
		return a.expandedInstruction
	}

	a.expandedInstruction = expanded
	return a.expandedInstruction
}

func (a *Agent) AddDate() bool {
	return a.addDate
}

func (a *Agent) AddEnvironmentInfo() bool {
	return a.addEnvironmentInfo
}

func (a *Agent) MaxIterations() int {
	return a.maxIterations
}

func (a *Agent) NumHistoryItems() int {
	return a.numHistoryItems
}

func (a *Agent) AddPromptFiles() []string {
	return a.addPromptFiles
}

// Description returns the agent's description
func (a *Agent) Description() string {
	return a.description
}

// WelcomeMessage returns the agent's welcome message
func (a *Agent) WelcomeMessage() string {
	return a.welcomeMessage
}

// SubAgents returns the list of sub-agent names
func (a *Agent) SubAgents() []*Agent {
	return a.subAgents
}

func (a *Agent) Handoffs() []*Agent {
	return a.handoffs
}

// Parents returns the list of parent agent names
func (a *Agent) Parents() []*Agent {
	return a.parents
}

// HasSubAgents checks if the agent has sub-agents
func (a *Agent) HasSubAgents() bool {
	return len(a.subAgents) > 0
}

// Model returns a random model from the available models
func (a *Agent) Model() provider.Provider {
	return a.models[rand.Intn(len(a.models))]
}

// Commands returns the named commands configured for this agent.
func (a *Agent) Commands() map[string]string {
	return a.commands
}

// Tools returns the tools available to this agent
func (a *Agent) Tools(ctx context.Context) ([]tools.Tool, error) {
	a.ensureToolSetsAreStarted(ctx)

	var agentTools []tools.Tool
	for _, toolSet := range a.toolsets {
		if !toolSet.started.Load() {
			// Toolset failed to start; skip it
			continue
		}
		ta, err := toolSet.Tools(ctx)
		if err != nil {
			slog.Warn("Toolset listing failed; skipping", "agent", a.Name(), "toolset", fmt.Sprintf("%T", toolSet.ToolSet), "error", err)
			a.addToolWarning(fmt.Sprintf("%T list failed: %v", toolSet.ToolSet, err))
			continue
		}
		agentTools = append(agentTools, ta...)
	}

	agentTools = append(agentTools, a.tools...)

	return agentTools, nil
}

func (a *Agent) ToolSets() []tools.ToolSet {
	var toolSets []tools.ToolSet

	for _, ts := range a.toolsets {
		toolSets = append(toolSets, ts)
	}

	return toolSets
}

func (a *Agent) ensureToolSetsAreStarted(ctx context.Context) {
	for _, toolSet := range a.toolsets {
		// Skip if toolset is already started
		if toolSet.started.Load() {
			continue
		}

		if err := toolSet.Start(ctx); err != nil {
			slog.Warn("Toolset start failed; skipping", "agent", a.Name(), "toolset", fmt.Sprintf("%T", toolSet.ToolSet), "error", err)
			a.addToolWarning(fmt.Sprintf("%T start failed: %v", toolSet.ToolSet, err))
			continue
		}

		// Mark toolset as started
		toolSet.started.Store(true)
	}
}

// addToolWarning records a warning generated while loading or starting toolsets.
func (a *Agent) addToolWarning(msg string) {
	if msg == "" {
		return
	}
	a.pendingWarnings = append(a.pendingWarnings, msg)
}

// DrainWarnings returns pending warnings and clears them.
func (a *Agent) DrainWarnings() []string {
	if len(a.pendingWarnings) == 0 {
		return nil
	}
	warnings := a.pendingWarnings
	a.pendingWarnings = nil
	return warnings
}

func (a *Agent) StopToolSets(ctx context.Context) error {
	for _, toolSet := range a.toolsets {
		// Only stop toolsets that are marked as started
		if !toolSet.started.Load() {
			continue
		}

		if err := toolSet.Stop(ctx); err != nil {
			return fmt.Errorf("failed to stop toolset: %w", err)
		}

		// Mark toolset as stopped
		toolSet.started.Store(false)
	}

	return nil
}
