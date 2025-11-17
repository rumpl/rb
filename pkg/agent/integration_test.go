package agent_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rumpl/rb/pkg/config"
	"github.com/rumpl/rb/pkg/teamloader"
)

func TestDynamicInstructionIntegration(t *testing.T) {
	// Create a temporary directory with a test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	require.NoError(t, os.WriteFile(testFile, []byte("test content"), 0o644))

	// Create a test configuration with dynamic instruction
	configContent := `version: "2"

models:
  dmr/test:1.0:
    provider: dmr
    model: test:1.0

agents:
  root:
    model: dmr/test:1.0
    instruction: |
      Files in current directory:
      ${list_directory({path: "."})}
    toolsets:
      - type: filesystem
        tools:
          - list_directory
`

	configPath := filepath.Join(tmpDir, "config.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte(configContent), 0o644))

	// Change to temp directory so list_directory works
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(oldWd) }()
	require.NoError(t, os.Chdir(tmpDir))

	// Load the team
	ctx := t.Context()
	team, err := teamloader.Load(ctx, configPath, config.RuntimeConfig{})
	require.NoError(t, err)
	require.NotNil(t, team)

	// Get the root agent
	agent, err := team.Agent("root")
	require.NoError(t, err)
	require.NotNil(t, agent)

	// Get the instruction - this should trigger expansion
	instruction := agent.Instruction(ctx)

	// Verify the instruction contains the expanded directory listing
	assert.Contains(t, instruction, "Files in current directory")
	assert.Contains(t, instruction, "test.txt", "instruction should contain the test file from directory listing")
	assert.Contains(t, instruction, "config.yaml", "instruction should contain config file from directory listing")
}

func TestDynamicInstructionWithNoTools(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test configuration without tools in instruction
	configContent := `version: "2"

models:
  dmr/test:1.0:
    provider: dmr
    model: test:1.0

agents:
  root:
    model: dmr/test:1.0
    instruction: |
      You are a simple assistant.
      No dynamic content here.
    toolsets:
      - type: filesystem
`

	configPath := filepath.Join(tmpDir, "config.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte(configContent), 0o644))

	// Load the team
	ctx := t.Context()
	team, err := teamloader.Load(ctx, configPath, config.RuntimeConfig{})
	require.NoError(t, err)

	// Get the root agent
	agent, err := team.Agent("root")
	require.NoError(t, err)

	// Get the instruction
	instruction := agent.Instruction(ctx)

	// Verify the instruction is unchanged
	assert.Contains(t, instruction, "You are a simple assistant")
	assert.Contains(t, instruction, "No dynamic content here")
	assert.NotContains(t, instruction, "${")
}
