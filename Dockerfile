FROM golang:alpine as builder
ENV APPDIR $GOPATH/src/github.com/Finatext/lapper
ENV GO111MODULE on
RUN apk update && apk add --no-cache git && mkdir -p $APPDIR
ADD . $APPDIR/
WORKDIR $APPDIR
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o lapper *.go

FROM public.ecr.aws/lambda/provided:al2
COPY --from=builder /go/src/github.com/Finatext/lapper/lapper ${LAMBDA_RUNTIME_DIR}/bootstrap
ADD entrypoint.sh ${LAMBDA_TASK_ROOT}
ENTRYPOINT ["/lambda-entrypoint.sh"]
CMD ["./entrypoint.sh"]
