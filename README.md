# Go Tetris

[![Go Tetris CI](https://github.com/briancain/go-tetris/actions/workflows/main.yml/badge.svg)](https://github.com/briancain/go-tetris/actions/workflows/main.yml)

A Tetris clone written in Go using the Ebiten game library.

Made with Amazon Q

## Features

- Classic Tetris gameplay
- Increasing difficulty levels
- Score tracking
- Next piece preview
- Ghost piece showing where the current piece will land
- Pause functionality

## Controls

- **Arrow Left/Right**: Move piece horizontally
- **Arrow Down**: Soft drop (accelerate downward)
- **Arrow Up**: Rotate piece
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

## Project Structure

- `cmd/`: Entry point for the application
- `internal/tetris/`: Core game logic and mechanics
- `internal/ui/`: Rendering and user interface components

## License

MIT
