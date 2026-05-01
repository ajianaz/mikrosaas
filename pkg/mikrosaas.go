// Package mikrosaas provides the public API for the SaaS Core Engine.
// Import this package to bootstrap a new SaaS application:
//
//	app := mikrosaas.NewApp(mikrosaas.Config{
//	    DatabaseURL: "postgres://...",
//	    JWTSecret:   "...",
//	})
//	app.Start(":8080")
package mikrosaas

// version is set at build time via -ldflags "-X github.com/ajianaz/mikrosaas.version=X.Y.Z"
var version = "dev"
