# toll_calculator

toll_calculator

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

the toll_calculator is a mono-repo consisting of the following microservices

- obu - [to generate random sample data]
- data_receiver
- distance_calculator
- aggregator
- gateway

pre-requisities

- run kafka locally using docker
- run the following services in order
  - agg receiver calc obu
  - test using the gateway by passing a valid obuid
