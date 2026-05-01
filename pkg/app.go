package mikrosaas

import (
	"fmt"
	"log/slog"
	"net/http"
)

// App is the core SaaS application. It wires together all dependencies
// and provides a ready-to-run http.Handler.
type App struct {
	config *Config
	router *http.ServeMux
	logger *slog.Logger
	server *http.Server
}

// NewApp creates a new SaaS application with the given configuration.
// Call Start() to begin serving requests.
func NewApp(cfg Config) (*App, error) {
	cfg.SetDefaults()
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("new app: %w", err)
	}

	logger := slog.Default()
	if cfg.IsDevelopment() {
		logger = slog.New(slog.NewTextHandler(http.ResponseWriter(nil), &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	app := &App{
		config: &cfg,
		router: http.NewServeMux(),
		logger: logger,
	}

	// TODO: Register routes when domain packages are implemented
	// app.registerRoutes()

	app.server = &http.Server{
		Addr:         cfg.Addr(),
		Handler:      app.router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return app, nil
}

// Handler returns the HTTP handler for testing or custom server setup.
func (a *App) Handler() http.Handler {
	return a.router
}

// Start begins serving HTTP requests. This blocks until the server is shut down.
func (a *App) Start() error {
	a.logger.Info("starting server",
		"addr", a.config.Addr(),
		"env", a.config.Env,
		"version", version,
	)
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}

// Shutdown gracefully stops the server.
func (a *App) Shutdown() error {
	a.logger.Info("shutting down server")
	return a.server.Close()
}

// Config returns a read-only copy of the application config.
func (a *App) Config() Config {
	return *a.config
}

// Version returns the application version.
func Version() string {
	return version
}
