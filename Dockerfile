FROM golang:1.19-alpine

WORKDIR /
COPY . .

RUN go mod tidy
RUN go build .

CMD ["./price-tracker"]