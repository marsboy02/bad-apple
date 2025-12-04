# bad-apple

Play the iconic "Bad Apple!!" Touhou animation as ASCII art in your terminal!

## Description

This project plays the famous "Bad Apple!!" animation as ASCII art directly in your terminal, with optional audio playback. The entire video (6500+ frames) and audio are embedded into a single binary, making it easy to distribute and run anywhere.

## Installation

### Via Homebrew (Coming Soon)

```bash
brew tap marsboy02/tap
brew install bad-apple
```

### From Source

```bash
git clone https://github.com/marsboy02/bad-apple.git
cd bad-apple
go build -o bad-apple
./bad-apple
```

## Usage

```bash
# Play video only (30fps default)
./bad-apple

# Play with audio
./bad-apple --audio

# Play audio only
./bad-apple --audio-only

# Loop indefinitely
./bad-apple --repeat

# Play N times
./bad-apple --times=3

# Custom frame rate
./bad-apple --fps=60

# Combine options
./bad-apple --audio --repeat --fps=30
```

## Options

- `--audio`: Play audio along with video
- `--audio-only`: Play audio only (no video)
- `--repeat`: Loop until interrupted (Ctrl+C)
- `--times=N`: Number of times to play (default: 1)
- `--fps=N`: Frames per second for video (default: 30)

## Requirements

- **Runtime**: macOS (uses `afplay` for audio playback)
- **Build**: Go 1.24.5 or later

## Building from Source

1. Clone the repository
2. Run `go build -o bad-apple`
3. The binary will include all assets (frames + audio)

### Regenerating ASCII Frames (Optional)

If you want to modify the video processing:

1. Extract video frames as PNG files to `frames/` directory
2. Install Python dependencies:
   ```bash
   python3 -m venv venv
   source venv/bin/activate
   pip install opencv-python
   ```
3. Run the generator:
   ```bash
   python scripts/generate_ascii_frames.py
   ```
4. Rebuild the Go binary

## Technical Details

- **Binary Size**: ~45MB (includes embedded assets)
- **Frame Format**: 128x40 characters, gzip-compressed
- **Audio Format**: WAV (38MB uncompressed)
- **Total Frames**: 6574 frames
- **Platform**: macOS only (audio playback uses `afplay`)

## License

MIT License - see [LICENSE](LICENSE) file for details

## Credits

- "Bad Apple!!" - Original song by ZUN, arranged by Alstroemeria Records
- Animation from Touhou Project
