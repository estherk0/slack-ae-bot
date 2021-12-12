FROM golang:1.17 AS builder

# set environment path
ENV PATH /go/bin:$PATH
WORKDIR /slack-ae-bot
COPY . .

RUN go mod tidy && go mod vendor
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/server ./cmd/server

FROM golang:1.17
LABEL AUTHOR Esther Kim (jabbukka@naver.com)

COPY --chown=0:0 --from=builder /go/bin/server /bin/
EXPOSE 3000

ENTRYPOINT ["/bin/server"]
