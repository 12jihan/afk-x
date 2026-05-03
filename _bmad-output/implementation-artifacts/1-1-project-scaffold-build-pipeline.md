# Story 1.1: Project Scaffold & Build Pipeline

Status: ready-for-dev

## Story

As a developer,
I want the Go project initialized with all dependencies, build tooling, and CI pipeline,
so that all subsequent development has a clean, reproducible foundation.

## Acceptance Criteria

1. **Given** the repository is cloned, **When** `make build` is run, **Then** a binary `bin/afk-x` is produced that exits cleanly with code 0
2. **Given** a commit is pushed to main, **When** GitHub Actions CI runs, **Then** `go test ./...` and `go vet ./...` both pass
3. **Given** the binary exists, **When** run with `--version`, **Then** it prints a version string and exits with code 0
4. **Given** `.goreleaser.yml` exists, **When** `goreleaser check` is run, **Then** no configuration errors are reported

## Tasks / Subtasks

- [ ] Task 1: Initialize Go module and install dependencies (AC: 1, 2)
  - [ ] Run `go mod init github.com/<username>/afk-x` (use a reasonable module path)
  - [ ] Run `go get github.com/charmbracelet/bubbletea@v1.3.10`
  - [ ] Run `go get github.com/charmbracelet/bubbles@v2.0.0`
  - [ ] Run `go get github.com/charmbracelet/lipgloss@v1.1.0`
  - [ ] Verify `go.mod` and `go.sum` are generated correctly

- [ ] Task 2: Create full project directory structure (AC: 1)
  - [ ] Create `cmd/afk-x/` directory
  - [ ] Create `internal/game/` directory
  - [ ] Create `internal/content/` directory
  - [ ] Create `internal/engine/` directory
  - [ ] Create `internal/combat/` directory
  - [ ] Create `internal/save/` directory
  - [ ] Create `internal/ui/` directory
  - [ ] Create `internal/ui/screens/` directory
  - [ ] Create `internal/ui/styles/` directory
  - [ ] Add `.gitignore` (standard Go .gitignore: bin/, vendor/, *.exe, etc.)

- [ ] Task 3: Create minimal `cmd/afk-x/main.go` entry point (AC: 1, 3)
  - [ ] Parse `--version` flag using standard `flag` package
  - [ ] If `--version`: print version string (hardcode `"0.1.0-dev"` initially) and `os.Exit(0)`
  - [ ] Default (no flags): print "afk-x: TUI not yet initialized" and exit cleanly (placeholder until Story 1.5)
  - [ ] Ensure binary exits with code 0 for all valid paths

- [ ] Task 4: Create stub package files so `go build ./...` succeeds (AC: 1, 2)
  - [ ] `internal/game/state.go` — package declaration + TODO comment (Story 1.2 will implement)
  - [ ] `internal/content/content.go` — package declaration + TODO comment (Story 1.3 will implement)
  - [ ] `internal/engine/engine.go` — package declaration + TODO comment (Story 1.4 will implement)
  - [ ] `internal/combat/combat.go` — package declaration + TODO comment
  - [ ] `internal/save/save.go` — package declaration + TODO comment
  - [ ] `internal/ui/ui.go` — package declaration + TODO comment
  - [ ] All stub files must compile without errors

- [ ] Task 5: Create Makefile (AC: 1, 2)
  - [ ] `build` target: `go build -o bin/afk-x ./cmd/afk-x`
  - [ ] `test` target: `go test ./...`
  - [ ] `lint` target: `go vet ./...`
  - [ ] `release` target: `goreleaser release --clean`
  - [ ] `clean` target: `rm -rf bin/`
  - [ ] Default target should be `build`

- [ ] Task 6: Create GitHub Actions CI workflow (AC: 2)
  - [ ] Create `.github/workflows/ci.yml`
  - [ ] Trigger: `push` to `main` and `pull_request` to `main`
  - [ ] Steps: checkout, setup-go (Go 1.21+), `go test ./...`, `go vet ./...`
  - [ ] Use `actions/checkout@v4` and `actions/setup-go@v5`

- [ ] Task 7: Create GoReleaser configuration (AC: 4)
  - [ ] Create `.goreleaser.yml` at project root
  - [ ] Configure builds for: `linux/amd64`, `linux/arm64`, `macos/amd64`, `macos/arm64`
  - [ ] Set `ldflags: ["-s", "-w"]` for binary size optimization (NFR15: ≤20MB)
  - [ ] Set binary name to `afk-x`
  - [ ] Run `goreleaser check` to validate (no actual release)

- [ ] Task 8: Write minimal test to satisfy CI (AC: 2)
  - [ ] Create `cmd/afk-x/main_test.go` with at least one passing test
  - [ ] Test can be a simple smoke test: verify `go test ./...` passes with non-zero test count
  - [ ] Run `go test ./...` locally to confirm

## Dev Notes

### Context: This is Story 1.1 — the Foundation

This is the very first story. No project files exist yet. Every file created here is NEW. The project does not exist on disk — the developer must create it from scratch.

**The output of this story is not a working game — it is a compilable, testable project scaffold that all subsequent stories build upon.** Story 1.5 will add the actual TUI. This story only needs `bin/afk-x --version` to work.

### Architecture Compliance (MUST FOLLOW)

From `architecture.md`:

**Package dependency rules (enforced by Go's import system):**
```
cmd/afk-x  →  ui, save, game
ui          →  engine, combat, game, content
engine      →  game, content
combat      →  game, content
save        →  game
content     →  (no internal imports)
game        →  (no internal imports)
```

**CRITICAL:** Do NOT import any `internal/` packages from `main.go` in this story. The main.go only needs `flag` and `fmt` from the standard library. No game logic, no Bubbletea, no save system in this story.

**Naming conventions:**
- Package names: lowercase single word (`engine`, `game`, `combat`, `save`, `content`, `ui`)
- File names: snake_case (`state.go`, `tick.go`, `atomic_write.go`)

### Technical Stack (EXACT VERSIONS — DO NOT DEVIATE)

| Dependency | Version | Import Path |
|---|---|---|
| bubbletea | v1.3.10 | `github.com/charmbracelet/bubbletea` |
| bubbles | v2.0.0 | `github.com/charmbracelet/bubbles` |
| lipgloss | v1.1.0 | `github.com/charmbracelet/lipgloss` |
| Go | 1.21+ | — |

**Why v1.3.10 for bubbletea:** v2 is RC, not production-ready. Use v1 stable.

**No CGo.** Standard `go build` only. Single binary output.

### Expected File Tree After This Story

```
afk-x/
├── .github/
│   └── workflows/
│       └── ci.yml
├── .gitignore
├── .goreleaser.yml
├── Makefile
├── README.md  (optional but nice)
├── go.mod
├── go.sum
├── bin/           (created by make build, git-ignored)
│   └── afk-x
└── cmd/
    └── afk-x/
        ├── main.go
        └── main_test.go
└── internal/
    ├── game/
    │   └── state.go        (stub: package game)
    ├── content/
    │   └── content.go      (stub: package content)
    ├── engine/
    │   └── engine.go       (stub: package engine)
    ├── combat/
    │   └── combat.go       (stub: package combat)
    ├── save/
    │   └── save.go         (stub: package save)
    └── ui/
        ├── ui.go           (stub: package ui)
        ├── screens/
        │   └── .gitkeep    (or stub file)
        └── styles/
            └── .gitkeep    (or stub file)
```

### main.go Implementation Pattern

```go
package main

import (
    "flag"
    "fmt"
    "os"
)

var version = "0.1.0-dev"

func main() {
    showVersion := flag.Bool("version", false, "print version and exit")
    flag.Parse()

    if *showVersion {
        fmt.Printf("afk-x %s\n", version)
        os.Exit(0)
    }

    // Placeholder — Story 1.5 will replace this with tea.NewProgram(...)
    fmt.Println("afk-x: TUI not yet initialized")
}
```

### Makefile Pattern

```makefile
.PHONY: build test lint clean release

build:
	go build -o bin/afk-x ./cmd/afk-x

test:
	go test ./...

lint:
	go vet ./...

clean:
	rm -rf bin/

release:
	goreleaser release --clean
```

### GoReleaser Pattern

```yaml
# .goreleaser.yml
version: 2
builds:
  - id: afk-x
    binary: afk-x
    main: ./cmd/afk-x
    goos:
      - linux
      - darwin
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w
archives:
  - format: tar.gz
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
```

### GitHub Actions CI Pattern

```yaml
# .github/workflows/ci.yml
name: CI
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - run: go test ./...
      - run: go vet ./...
```

### Testing Standards

- `go test ./...` must pass with at least one test (AC: 2 requires non-zero test execution)
- Create `cmd/afk-x/main_test.go` with a minimal smoke test
- Tests in `internal/*/` stubs can be empty — they just need to compile

### Stub File Pattern

All `internal/*/` stub files should follow this pattern:
```go
// Package game defines the core game state types.
// TODO: Implemented in Story 1.2
package game
```

No blank files — Go requires at least a package declaration.

### Project Structure Notes

- Module path should be `github.com/<username>/afk-x` — the developer can choose their GitHub username or use a placeholder like `github.com/player/afk-x`
- `bin/` must be in `.gitignore`
- GoReleaser v2 syntax uses `version: 2` at the top of `.goreleaser.yml`

### References

- Architecture: `_bmad-output/planning-artifacts/architecture.md` → Starter Template Evaluation, Package Structure, Infrastructure & Deployment sections
- PRD: `_bmad-output/planning-artifacts/prd.md` → NFR11 (single binary), NFR12 (Linux/macOS amd64/arm64), NFR15 (≤20MB)
- Epics: `_bmad-output/planning-artifacts/epics.md` → Epic 1, Story 1.1

## Dev Agent Record

### Agent Model Used

claude-sonnet-4-6

### Debug Log References

### Completion Notes List

- Ultimate context engine analysis completed - comprehensive developer guide created

### File List

(populate when story is implemented)
