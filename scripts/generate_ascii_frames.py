#!/usr/bin/env python3
import cv2
import gzip
import glob
import os

COLS = 128
ROWS = 40

# 어두움 → 밝음으로 갈수록 디테일한 문자
CHARS = " .`^\",:;Il!i~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"

def frame_to_ascii(path):
    img = cv2.imread(path, cv2.IMREAD_GRAYSCALE)
    if img is None:
        raise RuntimeError(f"failed to read {path}")
    # ffmpeg에서 이미 scale 했어도, 안전하게 한 번 더 맞춰주기
    img = cv2.resize(img, (COLS, ROWS))

    # 0.0 ~ 1.0으로 정규화
    norm = img.astype("float32") / 255.0

    # 감마 보정 (조금 더 밝게 보고 싶으면 0.8 같은 값으로 줄여도 됨)
    gamma = 1.0
    norm = norm ** gamma

    # 밝기 → 문자 인덱스
    idxs = (norm * (len(CHARS) - 1)).astype("int32")

    lines = []
    for y in range(ROWS):
        row = idxs[y]
        line = "".join(CHARS[i] for i in row)
        lines.append(line)
    return "\n".join(lines)

def main():
    frame_paths = sorted(glob.glob("frames/frame_*.png"))
    print(f"{len(frame_paths)} frames found")
    if not frame_paths:
        raise SystemExit("no frames found in frames/frame_*.png")

    os.makedirs("assets", exist_ok=True)

    frames_ascii = []
    for i, path in enumerate(frame_paths):
        ascii_frame = frame_to_ascii(path)
        frames_ascii.append(ascii_frame)

        # 디버깅용: 첫 프레임 미리보기 파일 생성
        if i == 0:
            with open("assets/preview_first_frame.txt", "w", encoding="utf-8") as f:
                f.write(ascii_frame)

    out_path = os.path.join("assets", "frames.txt.gz")
    print(f"writing {out_path}")

    with gzip.open(out_path, "wt", encoding="utf-8") as f:
        f.write("\n---FRAME---\n".join(frames_ascii))

if __name__ == "__main__":
    main()

