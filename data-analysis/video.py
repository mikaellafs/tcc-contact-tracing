import os
import datasource
import processing
import graph

from datetime import timedelta
import pandas as pd
import numpy as np
import cv2

# Get data
constant_df, metrics_df = datasource.get_from_generated()
constant_df = processing.parse_timestamp_types(constant_df).reset_index()

# Get total of days
min_day = constant_df['firstcontacttimestamp'].min()
max_day = constant_df['lastcontacttimestamp'].max()

total_days = min_day - max_day
total_days = abs(total_days.days)

print(min_day, max_day)

# Generate images
for shift in range(0, total_days):
    current_day = min_day.date() + timedelta(days=shift)
    next_day = current_day + timedelta(days=1)

    day_contacts = constant_df[constant_df['firstcontacttimestamp'].dt.day == current_day.day]
    day_contacts["total"] = pd.DataFrame(np.zeros(day_contacts.shape[0]), columns=["total"], index=day_contacts.index)

    print(day_contacts["total"])

    G = graph.build_graph_from_df(day_contacts, "distance", ["distance", "duration", "total"], isMulti=False)
    graph.draw(G, filename=f"video/graph-{shift}.png")

# Generate video
frameSize = (500, 500)
out = cv2.VideoWriter('generated/graph/contacts.mp4',cv2.VideoWriter_fourcc(*'MP4V'), total_days+1, frameSize)

for filename in os.listdir('generated/graph/video'):
    print(filename)
    if filename.endswith(".png"):
        img = cv2.imread(filename)
        out.write(img)

out.release()