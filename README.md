# go-wave
IN-PROGRESS Golang library for interacting with the Genki Instruments Wave

cmd/wave-cli-stream contains example code for connecting to a Wave ring and subscribing to the incoming Notifications. 

With go installed, you can build and run the example with:
```
cd cmd/wave-cli/stream
export GO111MODULE=on
go build
sudo ./wave-cli-stream
```
It will connect to the first detected device with the friendly name "Wave"