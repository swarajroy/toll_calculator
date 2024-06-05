brew install --build-from-source protobuf

go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

set go bin dir in PATH like this PATH="${PATH}:${HOME}/go/bin"

go get google.golang.org/protobuf
go get google.golang.org/grpc

'''
instrumentation
go get github.com/prometheus/client_golang/prometheus
'''

aggregator

- a microservice that is offering a http transport or a grpc transport to
  - accept a post request to aggregate the distance per obu
  - get the invoice by obu id computing the data in real time
