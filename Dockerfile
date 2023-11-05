FROM golang AS builder
WORKDIR $GOPATH/src/package/app/
COPY . .
RUN go mod download
RUN go get -d -v
RUN go install
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/main
FROM scratch
COPY --from=builder /go/bin/main main
EXPOSE 7000
ENTRYPOINT ["/main"]