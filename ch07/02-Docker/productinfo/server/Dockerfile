# 멀티 단계(Multistage) 빌드

# I 단계 빌드:
FROM golang AS build
ENV location /Users/vincent/Workspace/Writing/grpc-up-and-running/ch07/02-Docker/productinfo
WORKDIR ${location}/server

ADD ./server ${location}/server
ADD ./proto-gen ${location}/proto-gen

RUN go get -d ./...
RUN go install ./...

RUN CGO_ENABLED=0 go build -o /bin/grpc-productinfo-server

# II 단계 빌드:
FROM scratch
COPY --from=build /bin/grpc-productinfo-server /bin/grpc-productinfo-server

ENTRYPOINT ["/bin/grpc-productinfo-server"]
EXPOSE 50051