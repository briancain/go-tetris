# Go Tetris

[![Go Tetris CI](https://github.com/briancain/go-tetris/actions/workflows/main.yml/badge.svg)](https://github.com/briancain/go-tetris/actions/workflows/main.yml)

A Tetris clone written in Go using the Ebiten game library. Compliant with the official Tetris Guidelines.

## Features

- Classic Tetris gameplay following official Tetris Guidelines
- Standard 10×22 playfield (with top 2 rows hidden)
- Official Tetrimino colors (Cyan I, Yellow O, Purple T, Green S, Red Z, Blue J, Orange L)
- Super Rotation System (SRS) with proper wall kicks
- 7-bag Random Generator for fair piece distribution
- T-Spin detection and bonus scoring
- Back-to-Back bonus scoring
- Increasing difficulty levels
- Next piece preview
- Hold piece functionality
- Ghost piece showing where the current piece will land
- Pause functionality

## Controls

- **Arrow Left/Right**: Move piece horizontally
- **Arrow Down**: Soft drop (accelerate downward)
- **Arrow Up**: Rotate piece clockwise
- **Space**: Hard drop (instantly drop piece)
- **Escape**: Pause/Resume game
- **Shift**: Hold current piece for later use
- **Enter**: Start new game (from menu or game over screen)

## Requirements

- Go 1.18 or higher
- Ebiten v2 game library

## Building and Running

### Using Make

```bash
# Build the game
make build

# Run the game
make run

# Clean build artifacts
make clean
```

### Using Go directly

```bash
# Build the game
go build -o bin/tetris ./cmd

# Run the game
go run ./cmd
```

### Server Configuration

The multiplayer server supports configuration via CLI flags, environment variables, or defaults:

```bash
# Build and run the server
go build -o bin/server ./cmd/server
./bin/server

# With CLI flags (highest priority)
./bin/server -port 9000 -redis-url redis://prod:6379 -server-url https://api.example.com

# With environment variables (fallback)
PORT=9000 REDIS_URL=redis://prod:6379 ./bin/server

# Mixed (CLI flags override env vars)
PORT=7000 ./bin/server -port 9000  # Uses port 9000

# Show help
./bin/server -h
```

Configuration priority: **CLI flags** → **Environment variables** → **Defaults**

Configuration options:
- `PORT` / `-port`: Server port (default: 8080)
- `REDIS_URL` / `-redis-url`: Redis connection URL (default: redis://localhost:6379)
- `SERVER_URL` / `-server-url`: Public server URL (default: http://localhost:8080)

## Tetris Logo

To fully comply with the Tetris Guidelines, you need to obtain the official Tetris logo from The Tetris Company and place it in the `internal/ui/assets` directory as `tetris_logo.png`.

The Tetris brand and Tetris logos are trademarks of The Tetris Company, LLC.

## Project Structure

- `cmd/`: Entry point for the application
- `internal/tetris/`: Core game logic and mechanics
- `internal/ui/`: Rendering and user interface components
- `internal/ui/assets/`: Game assets including the Tetris logo

## License

MIT
