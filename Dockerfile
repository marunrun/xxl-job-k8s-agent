FROM golang:1.18.4-alpine as app-builder

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add git

WORKDIR /workspace

COPY . .
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go env -w GO111MODULE=on
RUN go build -o agent  ./cmd/agent.go


FROM alpine:3.16

WORKDIR /

COPY --from=app-builder /workspace/agent .

EXPOSE 9999

CMD ./agent