FROM golang AS builder
RUN apt update && apt upgrade
ENV $GOPATH=/usr/local/go/src
WORKDIR $GOPATH/app
COPY . .

RUN go get -d -v
RUN go install
RUN go build -o /go/bin/main
FROM scratch

COPY --from=builder /go/bin/main /go/bin/main
ENTRYPOINT ["/go/bin/main"]