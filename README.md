# afk-x

A terminal-based idle/incremental game built in Go. Ascend a computational tower while managing resources — CPU cycles, memory shards, and process threads — all rendered in a TUI with an ancient mainframe aesthetic.

## Features

- **Resource Tick Engine** — Three base resources accumulate at configurable rates with a 62.5 Hz tick loop
- **TUI Boot Sequence** — Themed system initialization sequence before entering the game
- **Narrative Milestones** — Flavor text at key floors as you ascend the tower
- **Pure Game Logic** — Engine and game state are fully deterministic with no side effects
- **JSON-First State** — Game state is designed for serialization from day one

## Requirements

- Go 1.26+
- Terminal: minimum 80x24

## Quick Start

```bash
# Build
make build

# Run
make run

# Or run directly
./bin/afk-x
```

## Development

```bash
# Run tests
make test

# Lint
make lint

# Clean build artifacts
make clean

# Build release archives
make release
```

## Project Structure

```
cmd/afk-x/          Entry point
internal/
  game/              Game state schema (pure data structures)
  engine/            Resource tick engine (pure functions)
  content/           Static game content — resources, narrative, constants
  ui/                Bubbletea TUI model, screens, and styles
  save/              Save/load persistence (planned)
  combat/            Combat system (planned)
```

## Controls

| Key          | Action |
|--------------|--------|
| `q`          | Quit   |
| `ctrl+c`    | Quit   |
| Any key      | Skip boot sequence |

## License

MIT
