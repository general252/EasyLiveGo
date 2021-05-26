

#### RTSP

```shell script
ffmpeg -re -i demo.flv -rtsp_transport tcp -vcodec h264 -f rtsp rtsp://192.168.6.80/test
ffmpeg -re -i demo.flv -rtsp_transport udp -vcodec h264 -f rtsp rtsp://192.168.6.80/test

ffplay -i rtsp://192.168.6.80/test
```

push
```
OPTIONS - ANNOUNCE - SETUP - SETUP - RECORD
```

pull
```
OPTIONS - DESCRIBE - SETUP - SETUP - PLAY
```

---

push
```

2021/05/26 19:54:08 收到数据:
OPTIONS rtsp://192.168.6.80:554/test RTSP/1.0
CSeq: 1
User-Agent: Lavf58.45.100

2021/05/26 19:54:08 发送数据:
RTSP/1.0 200 OK
CSeq: 1
Session: d761648ac176ecc3
Public: DESCRIBE, SETUP, TEARDOWN, PLAY, PAUSE, OPTIONS, ANNOUNCE, RECORD

2021/05/26 19:54:08 收到数据:
ANNOUNCE rtsp://192.168.6.80:554/test RTSP/1.0
Content-Type: application/sdp
CSeq: 2
User-Agent: Lavf58.45.100
Session: d761648ac176ecc3
Content-Length: 493

v=0
o=- 0 0 IN IP4 127.0.0.1
s=No Name
c=IN IP4 192.168.6.80
t=0 0
a=tool:libavformat 58.45.100
m=video 0 RTP/AVP 96
a=rtpmap:96 H264/90000
a=fmtp:96 packetization-mode=1; sprop-parameter-sets=Z2QAH6zZQFAFuwEQAAADABAAAAMDoPGDGWA=,aOvjyyLA; profile-level-id=64001F
a=control:streamid=0
m=audio 0 RTP/AVP 97
b=AS:128
a=rtpmap:97 MPEG4-GENERIC/48000/2
a=fmtp:97 profile-level-id=1;mode=AAC-hbr;sizelength=13;indexlength=3;indexdeltalength=3; config=119056E500
a=control:streamid=1
2021/05/26 19:54:08 audio codec[aac]
2021/05/26 19:54:08 video codec[h264]
2021/05/26 19:54:08 发送数据:
RTSP/1.0 200 OK
CSeq: 2
Session: d761648ac176ecc3

2021/05/26 19:54:08 收到数据:
SETUP rtsp://192.168.6.80:554/test/streamid=0 RTSP/1.0
Transport: RTP/AVP/UDP;unicast;client_port=32566-32567;mode=record
CSeq: 3
User-Agent: Lavf58.45.100
Session: d761648ac176ecc3

2021/05/26 19:54:08 发送数据:
RTSP/1.0 200 OK
Session: d761648ac176ecc3
Transport: RTP/AVP/UDP;unicast;client_port=32566-32567;server_port=5020-5021;mode=record
CSeq: 3

2021/05/26 19:54:08 收到数据:
SETUP rtsp://192.168.6.80:554/test/streamid=1 RTSP/1.0
CSeq: 4
User-Agent: Lavf58.45.100
Session: d761648ac176ecc3
Transport: RTP/AVP/UDP;unicast;client_port=32568-32569;mode=record

2021/05/26 19:54:08 发送数据:
RTSP/1.0 200 OK
CSeq: 4
Session: d761648ac176ecc3
Transport: RTP/AVP/UDP;unicast;client_port=32568-32569;server_port=5020-5021;mode=record

2021/05/26 19:54:08 收到数据:
RECORD rtsp://192.168.6.80:554/test RTSP/1.0
Range: npt=0.000-
CSeq: 5
User-Agent: Lavf58.45.100
Session: d761648ac176ecc3

2021/05/26 19:54:08 发送数据:
RTSP/1.0 200 OK
CSeq: 5
Session: d761648ac176ecc3
```

pull
```
2021/05/26 19:56:54 收到数据:
OPTIONS rtsp://192.168.6.80:554/test RTSP/1.0
CSeq: 1
User-Agent: Lavf58.45.100

2021/05/26 19:56:54 发送数据:
RTSP/1.0 200 OK
CSeq: 1
Session: 742ce62026bcf929
Public: DESCRIBE, SETUP, TEARDOWN, PLAY, PAUSE, OPTIONS, ANNOUNCE, RECORD

2021/05/26 19:56:54 收到数据:
DESCRIBE rtsp://192.168.6.80:554/test RTSP/1.0
User-Agent: Lavf58.45.100
Session: 742ce62026bcf929
Accept: application/sdp
CSeq: 2

2021/05/26 19:56:54 发送数据:
RTSP/1.0 200 OK
Session: 742ce62026bcf929
Content-Length: 493
CSeq: 2

v=0
o=- 0 0 IN IP4 127.0.0.1
s=No Name
c=IN IP4 192.168.6.80
t=0 0
a=tool:libavformat 58.45.100
m=video 0 RTP/AVP 96
a=rtpmap:96 H264/90000
a=fmtp:96 packetization-mode=1; sprop-parameter-sets=Z2QAH6zZQFAFuwEQAAADABAAAAMDoPGDGWA=,aOvjyyLA; profile-level-id=64001F
a=control:streamid=0
m=audio 0 RTP/AVP 97
b=AS:128
a=rtpmap:97 MPEG4-GENERIC/48000/2
a=fmtp:97 profile-level-id=1;mode=AAC-hbr;sizelength=13;indexlength=3;indexdeltalength=3; config=119056E500
a=control:streamid=1
2021/05/26 19:56:54 收到数据:
SETUP rtsp://192.168.6.80:554/test/streamid=0 RTSP/1.0
User-Agent: Lavf58.45.100
Session: 742ce62026bcf929
Transport: RTP/AVP/UDP;unicast;client_port=14832-14833
CSeq: 3

2021/05/26 19:56:54 发送数据:
RTSP/1.0 200 OK
CSeq: 3
Session: 742ce62026bcf929
Transport: RTP/AVP/UDP;unicast;client_port=14832-14833

2021/05/26 19:56:54 收到数据:
SETUP rtsp://192.168.6.80:554/test/streamid=1 RTSP/1.0
Transport: RTP/AVP/UDP;unicast;client_port=14834-14835
CSeq: 4
User-Agent: Lavf58.45.100
Session: 742ce62026bcf929

2021/05/26 19:56:54 发送数据:
RTSP/1.0 200 OK
CSeq: 4
Session: 742ce62026bcf929
Transport: RTP/AVP/UDP;unicast;client_port=14834-14835

2021/05/26 19:56:54 收到数据:
PLAY rtsp://192.168.6.80:554/test RTSP/1.0
User-Agent: Lavf58.45.100
Session: 742ce62026bcf929
Range: npt=0.000-
CSeq: 5

2021/05/26 19:56:54 发送数据:
RTSP/1.0 200 OK
Session: 742ce62026bcf929
Range: npt=0.000-
CSeq: 5

```
