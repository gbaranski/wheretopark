import os
import sys
import cv2
import copy
from pydantic import BaseModel
import yaml

url = str(sys.argv[1])
capture = cv2.VideoCapture(url)
ret, frame = capture.read()
print("Got frame from VideoCpature")

spots = []

points: list[(int, int)] = []
frames: list[any] = []

def draw(event, x, y, flags, param):
    if event == 1:
        points.append((x, y))
        frames.append(copy.deepcopy(frame))
        if len(points) > 1:
            cv2.line(frame, points[-1], points[-2], (0, 0, 255), 2)
        if len(points) == 4:
            cv2.line(frame, points[-1], points[0], (0, 0, 255), 2)
            spots.append(points.copy())
            points.clear()
            frames.clear()


cv2.namedWindow('frame')
cv2.setMouseCallback('frame', draw)

while True:
    cv2.imshow('frame', frame)
    k = cv2.waitKey(20) & 0xFF
    if k == 27:
        break
    elif k == 127:
        points.remove(points[-1])
        frame = frames.pop()
    elif k == ord('r'):
        os.execl(sys.executable, os.path.abspath(__file__), *sys.argv)


for points in spots:
    output = ("- points: %s" % points).replace("(", "[").replace(")", "]")
    print(output)

capture.release()
cv2.destroyAllWindows()