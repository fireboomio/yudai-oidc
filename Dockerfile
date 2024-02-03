FROM golang:alpine as builder
WORKDIR /yudai
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/yudai

FROM alpine
WORKDIR /yudai
COPY --from=builder /app/yudai .
EXPOSE 9825
CMD ["./yudai"]
