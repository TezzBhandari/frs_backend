FROM golang:1.22.5-alpine3.20

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o frs-api-server

CMD [ "./frs-api-server" ]