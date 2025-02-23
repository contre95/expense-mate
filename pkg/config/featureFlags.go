package config

import "time"

func (c *Config) StorageEngine() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Storage.Engine
}

func (c *Config) LoadSampleData() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Storage.LoadSampleData
}

func (c *Config) SQLitePath() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Storage.SQLitePath
}

func (c *Config) JSONStoragePath() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Storage.JSONStoragePath
}

// HTTP accessors
func (c *Config) CORSAllowList() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.HTTP.CORSAllowList
}

// Telegram accessors
func (c *Config) TelegramEnabled() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Telegram.Enabled
}

func (c *Config) TelegramToken() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Telegram.APIToken
}

// Ollama accessors
func (c *Config) OllamaEnabled() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Ollama.Enabled
}

func (c *Config) OllamaConfig() (endpoint, textModel, visionModel string, timeout time.Duration) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Ollama.Endpoint,
		c.Ollama.TextModel,
		c.Ollama.VisionModel,
		c.Ollama.Timeout
}
