// Build pipeline for ttyd
//
// Builds the ttyd web terminal, including the ghostty-web frontend
// (Node.js/Yarn 3) and the C backend (cmake/libwebsockets).

package main

import (
	"context"
	"dagger/ttyd/internal/dagger"
)

type Ttyd struct{}

// Build the frontend (html.h + wasm.h) from the html/ directory.
// Returns a directory containing the generated src/html.h and src/wasm.h.
func (m *Ttyd) Frontend(ctx context.Context, source *dagger.Directory) *dagger.Directory {
	return m.frontendContainer(source).
		Directory("/out")
}

func (m *Ttyd) frontendContainer(source *dagger.Directory) *dagger.Container {
	htmlDir := source.Directory("html")

	return dag.Container().
		From("node:18-bookworm").
		// Enable corepack for Yarn 3
		WithExec([]string{"corepack", "enable"}).
		// Mount html/ source
		WithDirectory("/app", htmlDir).
		WithWorkdir("/app").
		// Install dependencies
		WithExec([]string{"yarn", "install", "--immutable"}).
		// Production webpack build
		WithEnvVariable("NODE_ENV", "production").
		WithExec([]string{"yarn", "webpack"}).
		// Gulp: inline HTML, generate html.h and wasm.h
		WithExec([]string{"yarn", "gulp"}).
		// Copy generated headers to output
		WithExec([]string{"mkdir", "-p", "/out"}).
		WithExec([]string{"cp", "/app/../src/html.h", "/out/html.h"}).
		WithExec([]string{"cp", "/app/../src/wasm.h", "/out/wasm.h"})
}

// Build the ttyd binary. Returns a directory containing the compiled binary.
func (m *Ttyd) Build(ctx context.Context, source *dagger.Directory) *dagger.Directory {
	// First build the frontend to get html.h and wasm.h
	headers := m.Frontend(ctx, source)

	return dag.Container().
		From("debian:bookworm").
		// Install C build dependencies
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{
			"apt-get", "install", "-y", "--no-install-recommends",
			"build-essential", "cmake", "git",
			"libjson-c-dev", "libwebsockets-dev", "libuv1-dev",
			"zlib1g-dev", "libssl-dev",
		}).
		// Mount source tree, excluding host build artifacts
		WithDirectory("/src", source, dagger.ContainerWithDirectoryOpts{
			Exclude: []string{"build/", "html/node_modules/", "html/dist/", ".dagger/"},
		}).
		// Overlay generated headers into src/
		WithFile("/src/src/html.h", headers.File("html.h")).
		WithFile("/src/src/wasm.h", headers.File("wasm.h")).
		// cmake configure + build
		WithExec([]string{"mkdir", "-p", "/src/_build"}).
		WithWorkdir("/src/_build").
		WithExec([]string{"cmake", ".."}).
		WithExec([]string{"make", "-j"}).
		// Copy binary to clean output
		WithExec([]string{"mkdir", "-p", "/out"}).
		WithExec([]string{"cp", "/src/_build/ttyd", "/out/ttyd"}).
		Directory("/out")
}

// Build and export the ttyd binary to the host at the given path.
func (m *Ttyd) BuildLocal(ctx context.Context, source *dagger.Directory, output string) (string, error) {
	return m.Build(ctx, source).Export(ctx, output)
}

// Regenerate src/html.h and src/wasm.h from the frontend source.
// Run this after bumping ghostty-web or changing html/src/.
//
// Usage: dagger call generate --source=. export --path=./src
func (m *Ttyd) Generate(ctx context.Context, source *dagger.Directory) *dagger.Directory {
	return m.Frontend(ctx, source)
}

// Run the frontend yarn install (useful for regenerating yarn.lock).
// Returns the updated html/ directory including the new yarn.lock.
func (m *Ttyd) YarnInstall(ctx context.Context, source *dagger.Directory) *dagger.Directory {
	htmlDir := source.Directory("html")

	return dag.Container().
		From("node:18-bookworm").
		WithExec([]string{"corepack", "enable"}).
		WithDirectory("/app", htmlDir).
		WithWorkdir("/app").
		WithExec([]string{"yarn", "install"}).
		Directory("/app")
}
