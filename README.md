# wsxhub

wsxhub is a websocket server and client for using from other tools.  
**This is in development.**

## Command
- wsxhubd  
    - a websocket server daemon
- wsxhub  
    - a client command for requesting to wsxhubd
    - input request from the command option
    - output responses to stdout
    - filter the received json by --key, --filter options

## Install
```
go get -u github.com/notomo/wsxhub/...
```

## Usage
```
# start server and wait
wsxhubd 

# send {"key":"value"} to server
wsxhub send --json {\"key\":\"value\"} 

# receive only json has {"key":1}
wsxhub --filter {\"key\":1} receive

# receive only json has {"key":any}
wsxhub --key {\"key\":true} receive
```
