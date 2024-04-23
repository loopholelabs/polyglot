module benchmark

go 1.21.4

replace github.com/loopholelabs/polyglot/v2 => ../

require (
	github.com/loopholelabs/polyglot/v2 v2.0.0
	google.golang.org/grpc v1.59.0
	google.golang.org/protobuf v1.33.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.23.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect; indirect..
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231016165738-49dd2c1f3d0b // indirect
)
