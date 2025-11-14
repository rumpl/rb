package config

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/goccy/go-yaml"

	latest "github.com/rumpl/rb/pkg/config/v2"
	"github.com/rumpl/rb/pkg/environment"
	"github.com/rumpl/rb/pkg/filesystem"
)

func LoadConfig(path string, fs filesystem.FS) (*latest.Config, error) {
	data, err := fs.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %s: %w", path, err)
	}

	return LoadConfigBytes(data)
}

func LoadConfigBytes(data []byte) (*latest.Config, error) {
	config, err := parseCurrentVersion(data)
	if err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// CheckRequiredEnvVars checks which environment variables are required by the models and tools.
//
// This allows exiting early with a proper error message instead of failing later when trying to use a model or tool.
func CheckRequiredEnvVars(ctx context.Context, cfg *latest.Config, env environment.Provider, runtimeConfig RuntimeConfig) error {
	missing, err := gatherMissingEnvVars(ctx, cfg, env, runtimeConfig)
	if err != nil {
		// If there's a tool preflight error, log it but continue
		slog.Warn("Failed to preflight toolset environment variables; continuing", "error", err)
	}

	// Return error if there are missing environment variables
	if len(missing) > 0 {
		return &environment.RequiredEnvError{
			Missing: missing,
		}
	}

	return nil
}

func parseCurrentVersion(data []byte) (latest.Config, error) {
	options := []yaml.DecodeOption{yaml.Strict()}

	var cfg latest.Config
	err := yaml.UnmarshalWithOptions(data, &cfg, options...)
	return cfg, err
}

func validateConfig(cfg *latest.Config) error {
	if cfg.Models == nil {
		cfg.Models = map[string]latest.ModelConfig{}
	}

	for name := range cfg.Models {
		if cfg.Models[name].ParallelToolCalls == nil {
			m := cfg.Models[name]
			m.ParallelToolCalls = boolPtr(true)
			cfg.Models[name] = m
		}
	}

	for agentName := range cfg.Agents {
		agent := cfg.Agents[agentName]

		modelNames := strings.SplitSeq(agent.Model, ",")
		for modelName := range modelNames {
			if _, exists := cfg.Models[modelName]; exists {
				continue
			}

			provider, model, ok := strings.Cut(modelName, "/")
			if !ok {
				return fmt.Errorf("agent '%s' references non-existent model '%s'", agentName, modelName)
			}

			cfg.Models[modelName] = latest.ModelConfig{
				Provider: provider,
				Model:    model,
			}
		}

		for _, subAgentName := range agent.SubAgents {
			if _, exists := cfg.Agents[subAgentName]; !exists {
				return fmt.Errorf("agent '%s' references non-existent sub-agent '%s'", agentName, subAgentName)
			}
		}
	}

	return nil
}

func boolPtr(b bool) *bool {
	return &b
}
