# バイナリ作成用コンテナ
FROM golang:1.18.5-bullseye as deploy-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-w -s" -o app


# デプロイ用コンテナ
FROM debian:bullseye-slim as deploy

RUN apt-get update

COPY --from=deploy-builder /app/app .

CMD ["./app"]


# 開発用ホットリロード環境
FROM golang:1.18.5 as dev
WORKDIR /app
RUN go install github.com/cosmtrek/air@latest
CMD ["air"]

# コンテナベースイメージ選定基準 cf.『実用Go言語』14.2