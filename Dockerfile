FROM golang:1.15-alpine
LABEL creator = "Email: zumikq78@gmail.com Telegram: @YUART"
ENV DB_TOKEN = "postgres://ttsemhpm:Resem3sFwwvOjTRjo5tU9aSSQwvKaJVE@balarama.db.elephantsql.com:5432/ttsemhpm"
ENV BOT_TOKEN = "1022122500:AAEtOpzz2EKXi0Os2kiJXCvYg6zbh8M3kF8"
RUN apk add --no-cache git
WORKDIR /app/anonchat-bot
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o ./out/anonchat-bot .
EXPOSE 5432
VOLUME /bot_data_volume
ENTRYPOINT ["./out/anonchat-bot"]