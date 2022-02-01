# Build Stage
FROM golang:1.16 AS builder

# 複製原始碼
COPY . /app
WORKDIR /app

# 進行編譯
RUN go build -o heroku-line-bot

# Final Stage
FROM golang:1.16
COPY --from=builder /app/heroku-line-bot /app/heroku-line-bot
COPY --from=builder /app/resource /app/resource
WORKDIR /app

ENV WORK_DIR=app

CMD [ "sh", "-c", "./heroku-line-bot server" ]