from picamera2 import Picamera2
import cv2
import numpy as np
import time

print(1)

# --- Initialize Picamera2 ---
picam2 = Picamera2()
config = picam2.create_preview_configuration({"format": "XBGR8888", "size": (640, 480)})
picam2.configure(config)
picam2.start()
print(2)

# --- Background subtractor ---
bg_subtractor = cv2.createBackgroundSubtractorMOG2()

# --- Current mode ---
mode = "normal"

# --- FPS calculation ---
prev_time = time.time()

# --- Video writer setup ---
save_fps = 30.0
fourcc = cv2.VideoWriter_fourcc(*"XVID")
out = cv2.VideoWriter("output.avi", fourcc, save_fps, (640, 480))
print(3)

while True:
    frame = picam2.capture_array()
    if frame is None:
        print("No frame, skipping")
        continue

    frame = cv2.flip(frame, 1)
    display_frame = frame.copy()

    # --- Mode processing ---
    if mode == "threshold":
        gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)
        _, display_frame = cv2.threshold(gray, 127, 255, cv2.THRESH_BINARY)
        display_frame = cv2.cvtColor(display_frame, cv2.COLOR_GRAY2BGR)

    elif mode == "edge":
        gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)
        edges = cv2.Canny(gray, 100, 200)
        display_frame = cv2.cvtColor(edges, cv2.COLOR_GRAY2BGR)

    elif mode == "bg_sub":
        fg_mask = bg_subtractor.apply(frame)
        display_frame = cv2.bitwise_and(frame, frame, mask=fg_mask)

    elif mode == "contour":
        gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)
        _, thresh = cv2.threshold(gray, 127, 255, cv2.THRESH_BINARY)
        contours, _ = cv2.findContours(thresh, cv2.RETR_TREE, cv2.CHAIN_APPROX_SIMPLE)
        display_frame = cv2.drawContours(frame.copy(), contours, -1, (0, 255, 0), 2)

    # --- Calculate actual FPS ---
    curr_time = time.time()
    processing_fps = 1.0 / (curr_time - prev_time) if curr_time != prev_time else 0.0
    prev_time = curr_time

    # --- Display FPS and mode ---
    cv2.putText(
        display_frame,
        f"FPS: {int(processing_fps)} \nMode: {mode}",
        (10, 30),
        cv2.FONT_HERSHEY_SIMPLEX,
        1,
        (0, 255, 0),
        2,
    )

    # --- Write frame to video ---
    out.write(display_frame)

    # --- Show video ---
    cv2.imshow("Live Video", display_frame)

    # --- Key handling ---
    key = cv2.waitKey(1) & 0xFF
    if key == ord("t"):
        mode = "threshold"
    elif key == ord("e"):
        mode = "edge"
    elif key == ord("b"):
        mode = "bg_sub"
    elif key == ord("c"):
        mode = "contour"
    elif key == ord("n"):
        mode = "normal"
    elif key == ord("q"):
        break

# --- Clean up ---
picam2.stop()
out.release()
cv2.destroyAllWindows()
