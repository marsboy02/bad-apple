# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a "Bad Apple!!" ASCII animation player written in Go that plays the famous Touhou animation as ASCII art in the terminal. The project embeds both compressed ASCII frames and audio into a single binary using Go's embed feature.

## Architecture

### Core Components

1. **main.go**: Single-file Go application that handles:
   - Frame loading from embedded gzip-compressed data
   - Audio playback using macOS `afplay` command
   - Terminal control (cursor hiding, screen clearing)
   - Playback loop with configurable FPS
   - Signal handling for graceful shutdown

2. **Embedded Assets** (via `//go:embed`):
   - `assets/frames.txt.gz`: Gzip-compressed ASCII frames separated by `\n---FRAME---\n`
   - `assets/bad_apple.wav`: Audio file (38MB WAV)

3. **Frame Generation Pipeline** (`scripts/generate_ascii_frames.py`):
   - Converts PNG frames to ASCII art using OpenCV
   - Uses character gradient from dark to bright: ` .`^\",:;Il!i~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$`
   - Output dimensions: 128 columns × 40 rows
   - Compresses output to `assets/frames.txt.gz`

### Key Design Patterns

- **Single Binary Distribution**: All assets embedded at build time, no external files needed
- **Async Audio Playback**: When playing video+audio together, audio plays in background goroutine
- **Blocking Audio Playback**: Audio-only mode uses synchronous playback for proper looping
- **Platform Dependency**: Uses macOS `afplay` command for audio (not cross-platform)

## Building and Running

### Build the Binary
```bash
go build -o bad-apple
```

This creates a ~45MB binary with all assets embedded.

### Run Commands
```bash
# Basic playback (video only, 30fps)
./bad-apple

# With audio
./bad-apple --audio

# Audio only
./bad-apple --audio-only

# Loop indefinitely
./bad-apple --repeat

# Play N times
./bad-apple --times=3

# Custom FPS
./bad-apple --fps=60
```

## Development Workflow

### Regenerating ASCII Frames

If you need to modify the video or ASCII conversion:

1. Extract video frames to `frames/` directory as `frame_0001.png`, `frame_0002.png`, etc.
2. Ensure Python environment is activated:
   ```bash
   source venv/bin/activate
   ```
3. Run the frame generator (requires opencv-python):
   ```bash
   python scripts/generate_ascii_frames.py
   ```
4. This creates `assets/frames.txt.gz` which gets embedded on next build

### Replacing Audio

Replace `assets/bad_apple.wav` with your audio file, then rebuild. The file must be WAV format compatible with `afplay`.

## Important Implementation Details

- **Frame Format**: Frames are stored as plain text separated by `\n---FRAME---\n`, then gzip-compressed
- **Terminal Control**: Uses ANSI escape sequences (e.g., `\033[2J` for clear, `\033[?25l` for hiding cursor)
- **Signal Handling**: Catches Ctrl+C to restore cursor visibility before exit
- **Temp File Management**: Audio playback writes WAV to temp file, deleted after playback completes
- **Ticker Precision**: Uses `time.Ticker` for frame timing, calculated as `time.Second / fps`

## Module Information

- Module path: `github.com/marsboy02/bad-apple`
- Go version: 1.24.5
- No external Go dependencies (uses only standard library)
