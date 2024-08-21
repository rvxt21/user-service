FROM golang

WORKDIR /app

COPY . .

RUN go mod download

WORKDIR /app/cmd/api

RUN go build -o api .

CMD [ "./api" ]