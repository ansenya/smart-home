package main

import (
	"context"
	"fmt"
	"llm-service/internal/agents"
	"llm-service/internal/agents/runtime"
	"llm-service/internal/agents/tools"
	"llm-service/internal/agents/tools/smarthome"
	anthropicClient "llm-service/internal/clients/anthropic"
	openaiClient "llm-service/internal/clients/openaicompat"
	"llm-service/internal/config"
	"llm-service/internal/handlers"
	"llm-service/internal/infra/db"
	"llm-service/internal/infra/mqtt"
	"llm-service/internal/repositories"
	"llm-service/internal/services"
	"log/slog"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.NewConfig()
	cfg.Log.Info("llm service started")

	// DB
	dbClient := db.NewClient(cfg)
	if err := dbClient.Connect(ctx); err != nil {
		cfg.Log.Error("db connect failed", slog.Any("err", err))
		os.Exit(1)
	}
	defer dbClient.Close()

	mqttClient, err := mqtt.NewClient(cfg.Log)
	if err != nil {
		cfg.Log.Error("mqtt connect failed", slog.Any("err", err))
		os.Exit(1)
	}
	defer mqttClient.Close()

	repos := repositories.NewContainer(dbClient)

	// Provider registry
	registry := agents.NewProviderRegistry()

	// OpenAI-compatible (ChatGPT)
	if cfg.OpenaiConfig.ApiKey != "" {
		oai, err := openaiClient.New("openai", "", cfg.OpenaiConfig.ApiKey, cfg.OpenaiConfig.ProxyURL)
		if err != nil {
			cfg.Log.Error("openai client init failed", slog.Any("err", err))
		} else {
			registry.Register(oai, "gpt-4o", "gpt-4o-mini", "gpt-5.2", "o3", "o4-mini")
		}
	}

	// Claude (Anthropic)
	if cfg.AnthropicConfig.APIKey != "" {
		claude, err := anthropicClient.New(cfg.AnthropicConfig.APIKey, cfg.AnthropicConfig.ProxyURL)
		if err != nil {
			cfg.Log.Error("anthropic client init failed", slog.Any("err", err))
		} else {
			registry.RegisterPrefix(claude, "claude-")
		}
	}

	// Tool registry
	toolReg := tools.NewRegistry()
	toolDefs := smarthome.Register(dbClient.DB, mqttClient.Client, toolReg)
	toolExec := tools.NewExecutor(toolReg)
	toolLoop := runtime.NewToolLoop(&tools.RuntimeAdapter{Exec: toolExec}, 8)

	orchestrator := agents.NewOrchestrator(registry, toolLoop, toolDefs)

	svcs, err := services.NewContainer(cfg, repos, orchestrator)
	if err != nil {
		panic(err)
	}

	router := handlers.NewRouter(cfg, svcs, repos)
	if err := router.Run(); err != nil {
		cfg.Log.Error(fmt.Sprintf("failed to start: %s", err))
	}
}
