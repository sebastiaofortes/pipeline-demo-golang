FROM golang:1.20 as builder
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .
#ENTRYPOINT [ "/go/src/app/main" ]

FROM scratch
COPY --from=builder /go/src/app/main /main
ENTRYPOINT [ "/main" ]