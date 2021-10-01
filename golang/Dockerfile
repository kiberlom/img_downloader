FROM golang:1.16.5 as build

# EXPOSE 8000

RUN mkdir /app
WORKDIR /app
COPY . .
RUN go mod tidy
WORKDIR /app/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o main .


# CMD ./main

FROM alpine:latest
RUN mkdir /app
WORKDIR /app
COPY --from=build /app/cmd/main .
COPY --from=build /app/internal/webserver/templates/  ./templates
CMD ["./main"]

