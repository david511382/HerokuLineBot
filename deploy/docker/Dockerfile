# Build Stage
FROM golang:1.18 AS builder

# 複製原始碼
COPY ./bootstrap /app/bootstrap
COPY ./cmd /app/cmd
COPY ./src /app/src
COPY ./docs /app/docs
COPY ./main.go /app/main.go
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
WORKDIR /app

# 進行編譯
RUN go build -o heroku-line-bot

# Final Stage
FROM golang:1.18
COPY --from=builder /app/heroku-line-bot /app/heroku-line-bot
COPY --from=builder /app/docs /app/docs

COPY ./config /app/config
COPY ./resource /app/resource
WORKDIR /app

CMD [ "sh", "-c", "./heroku-line-bot server" ]