

#### RTSP

```shell script
cd E:/code/livego
ffmpeg -re -i demo.flv -rtsp_transport tcp -vcodec h264 -f rtsp rtsp://192.168.6.80/test
ffmpeg -re -i demo.flv -rtsp_transport udp -vcodec h264 -f rtsp rtsp://192.168.6.80/test

ffplay -i rtsp://192.168.6.80/test
```
