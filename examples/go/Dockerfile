FROM golang:alpine AS builder
ENV APPDIR=$GOPATH/src/github.com/Finatext/lapper/examples/go
ENV GO111MODULE=on
RUN apk update && rm -rf /var/cache/apk/* && mkdir -p $APPDIR
ADD . $APPDIR/
WORKDIR $APPDIR
RUN CGO_ENABLED=0 go build -mod=vendor -ldflags "-s -w" -o go main.go

FROM lapper
WORKDIR ${LAMBDA_TASK_ROOT}
COPY --from=builder /go/src/github.com/Finatext/lapper/examples/go ${LAMBDA_TASK_ROOT}
CMD ["./go"]
