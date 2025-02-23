package config

import (
	"os"
	"strconv"
	"sync"
	"time"
)

const OLLAMA_ENDPOINT = "OLLAMA_ENDPOINT"
const OLLAMA_TEXT_MODEL = "OLLAMA_TEXT_MODEL"
const OLLAMA_VISION_MODEL = "OLLAMA_VISION_MODEL"
const TELEGRAM_APITOKEN = "TELEGRAM_APITOKEN"

type Config struct {
	mu sync.RWMutex

	Storage struct {
		Engine          string
		LoadSampleData  bool
		SQLitePath      string
		JSONStoragePath string
	}

	HTTP struct {
		CORSAllowList string
	}

	Telegram struct {
		Enabled  bool
		APIToken string
	}

	Ollama struct {
		Enabled     bool
		Endpoint    string
		TextModel   string
		VisionModel string
		Timeout     time.Duration
	}
}

// Loads initial/default configurations based on env variables
func Load() *Config {
	cfg := &Config{
		Ollama: struct {
			Enabled     bool
			Endpoint    string
			TextModel   string
			VisionModel string
			Timeout     time.Duration
		}{
			Timeout: 300 * time.Second,
		},
	}

	// Storage configuration
	cfg.Storage.Engine = os.Getenv("STORAGE_ENGINE")
	if cfg.Storage.Engine == "" {
		cfg.Storage.Engine = "sqlite"
	}

	cfg.Storage.SQLitePath = os.Getenv("SQLITE_PATH")
	if cfg.Storage.SQLitePath == "" {
		cfg.Storage.SQLitePath = "./exp.db"
	}

	cfg.Storage.JSONStoragePath = os.Getenv("JSON_STORAGE_PATH")
	if cfg.Storage.JSONStoragePath == "" {
		cfg.Storage.JSONStoragePath = "./users.json"
	}

	if loadSample, err := strconv.ParseBool(os.Getenv("LOAD_SAMPLE_DATA")); err == nil {
		cfg.Storage.LoadSampleData = loadSample
	}

	// HTTP configuration
	cfg.HTTP.CORSAllowList = os.Getenv("CORS_ALLOWLIST")
	if cfg.HTTP.CORSAllowList == "" {
		cfg.HTTP.CORSAllowList = "*"
	}
	// Telegram configuration
	if token := os.Getenv(TELEGRAM_APITOKEN); token != "" {
		cfg.Telegram.Enabled = true
		cfg.Telegram.APIToken = token
	}

	// NOTE: Currently Ollma can be load separately from telegram
	// cause I'm planning to build ollama into the Web UI as well.

	// Ollama configuration
	cfg.Ollama.Endpoint = os.Getenv(OLLAMA_ENDPOINT)
	cfg.Ollama.TextModel = os.Getenv(OLLAMA_TEXT_MODEL)
	cfg.Ollama.VisionModel = os.Getenv(OLLAMA_VISION_MODEL)
	// Enforcing that at least on model should be specified
	cfg.Ollama.Enabled = cfg.Ollama.Endpoint != "" && (cfg.Ollama.TextModel != "" || cfg.Ollama.VisionModel != "")

	return cfg
}

// Update. Use the Update closure for thread-safe configuration changes.
// All modifications within the closure execute atomically under mutex protection.
// This is the only method that should be used to udpate configurations on runtime.
func (c *Config) Update(fn func(*Config)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	fn(c)
}
