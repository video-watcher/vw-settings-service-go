#!/bin/bash
set -e

echo "Building Docker image..."
docker build -t vw-settings-service-go:latest .

echo "Saving image..."
docker save vw-settings-service-go:latest | gzip > /tmp/vw-settings-service-go.tar.gz

echo "Copying to staging server..."
scp -i ~/.ssh/video_watcher_staging /tmp/vw-settings-service-go.tar.gz root@43.228.212.83:/tmp/

echo "Loading image on staging..."
ssh -i ~/.ssh/video_watcher_staging root@43.228.212.83 "docker load < /tmp/vw-settings-service-go.tar.gz"

echo "Deploy complete!"
