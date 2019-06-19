# wsxhub

wsxhub is a websocket server and client for using from other tools.  
**This is in development.**

## Install
```
go get -u github.com/notomo/wsxhub/...
```

## Usage
```
# start server and wait
wsxhub server

# send {"key":"value"} to server
echo '{"key":"value"}' | wsxhub send

# receive only json has {"key":1}
wsxhub receive --filter '{"operator": "and", "filters": [{"type": "exact", "map": {"key":1}}]}'
```
