#!/bin/bash
#echo "hello world"

xhost local:root

docker run -it --name gstreamer_container --privileged --net=host \
\
-v ~/.Xauthority:/root/.Xauthority \
-v /tmp/.X11-unix:/tmp/.X11-unix \
-e DISPLAY=$DISPLAY \
-e HTTP_PROXY=$HTTP_PROXY \
-e HTTPS_PROXY=$HTTPS_PROXY \
-e http_proxy=$http_proxy \
-e https_proxy=$https_proxy \
\
-v ~/gva/data/models/intel:/root/intel_models:ro \
-v ~/gva/data/models/common:/root/common_models:ro \
-e MODELS_PATH=/root/intel_models:/root/common_models \
\
-v ~/gva/data/video:/root/video-examples:ro \
-e VIDEO_EXAMPLES_DIR=/root/video-examples \
\
adi6496/gst-video-analytics:latest 