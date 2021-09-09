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
COPY --from=build /app/cmd/main .
CMD ["./main"]

