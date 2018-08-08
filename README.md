# wsbridge

A tiny websocket bridge server mainly for websocket client which does not support HTTP proxy.
This server refers `http_proxy`, `https_proxy`, and `no_proxy` environment variable and connect to the server.

## Usage

First of all, install *wsbridge* and start it for `wss://echo.websocket.org/echo` like:

```
$ go install github.com/lambdalisue/wsbridge
$ wsbridge wss://echo.websocket.org/echo
*******************************************************************************************
wsbridge - Tiny websocket connection bridge server

  listen: ws://127.0.0.1:8080/echo
  server: wss://echo.websocket.org/echo
  proxy:  system

*******************************************************************************************
```

Then connect to `ws://127.0.0.1:8080/echo`.
