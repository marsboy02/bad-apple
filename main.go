package main

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const frameDivider = "\n---FRAME---\n"

// 프레임 및 오디오를 바이너리에 임베드
//go:embed assets/frames.txt.gz
var framesGz []byte

//go:embed assets/bad_apple.wav
var audioWav []byte

func main() {
	// ------- 플래그 정의 -------
	audioFlag := flag.Bool("audio", false, "play audio along with video")
	audioOnlyFlag := flag.Bool("audio-only", false, "play audio only (no video)")
	repeatFlag := flag.Bool("repeat", false, "loop until interrupted")
	timesFlag := flag.Int("times", 1, "number of times to play (ignored if --repeat)")
	fpsFlag := flag.Float64("fps", 30, "frames per second for video")
	flag.Parse()

	// Ctrl+C 시 터미널 상태 복구
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		showCursor()
		fmt.Print("\033[0m") // 색 초기화
		os.Exit(0)
	}()

	// ------- 오디오 전용 모드 -------
	if *audioOnlyFlag {
		playCount := 1
		if !*repeatFlag {
			if *timesFlag > 0 {
				playCount = *timesFlag
			}
		}

		if *repeatFlag {
			for {
				if err := playAudioBlocking(); err != nil {
					fmt.Fprintf(os.Stderr, "failed to play audio: %v\n", err)
					break
				}
			}
		} else {
			for i := 0; i < playCount; i++ {
				if err := playAudioBlocking(); err != nil {
					fmt.Fprintf(os.Stderr, "failed to play audio: %v\n", err)
					break
				}
			}
		}
		return
	}

	// ------- 비디오(ASCII) + 선택적 오디오 -------

	frames, err := loadFrames()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load frames: %v\n", err)
		os.Exit(1)
	}

	hideCursor()
	defer showCursor()
	clearScreen()

	// 한 번 재생하는 함수 (오디오는 선택적으로)
	playOnce := func(withAudio bool) {
		if withAudio {
			if err := playAudioAsync(); err != nil {
				fmt.Fprintf(os.Stderr, "failed to play audio: %v\n", err)
			}
		}

		interval := time.Duration(float64(time.Second) / *fpsFlag)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for _, frame := range frames {
			<-ticker.C
			clearScreen()
			fmt.Print(frame)
		}
	}

	if *repeatFlag {
		for {
			playOnce(*audioFlag)
		}
	} else {
		playCount := 1
		if *timesFlag > 0 {
			playCount = *timesFlag
		}

		for i := 0; i < playCount; i++ {
			playOnce(*audioFlag)
		}
	}
}

// ----------------- 프레임 로딩 -----------------

func loadFrames() ([]string, error) {
	r, err := gzip.NewReader(bytes.NewReader(framesGz))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	raw := string(buf)
	parts := strings.Split(raw, frameDivider)
	return parts, nil
}

// ----------------- 오디오 재생 -----------------

// 오디오를 동기(blocking) 방식으로 1회 재생 (audio-only 모드에서 사용)
func playAudioBlocking() error {
	tmpFile, err := os.CreateTemp("", "bad-apple-*.wav")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(audioWav); err != nil {
		tmpFile.Close()
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}

	cmd := exec.Command("afplay", tmpFile.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// 비디오와 함께 재생할 때는 비동기(async)로 재생
func playAudioAsync() error {
	tmpFile, err := os.CreateTemp("", "bad-apple-*.wav")
	if err != nil {
		return err
	}
	if _, err := tmpFile.Write(audioWav); err != nil {
		tmpFile.Close()
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}

	cmd := exec.Command("afplay", tmpFile.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		_ = cmd.Wait()
		_ = os.Remove(tmpFile.Name())
	}()

	return nil
}

// ----------------- 터미널 유틸 -----------------

func clearScreen() {
	fmt.Print("\033[2J")
	moveCursorHome()
}

func moveCursorHome() {
	fmt.Print("\033[H")
}

func hideCursor() {
	fmt.Print("\033[?25l")
}

func showCursor() {
	fmt.Print("\033[?25h")
}
