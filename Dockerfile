FROM golang:1.19-alpine
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o urlshortner main.go 
EXPOSE 8000
ENTRYPOINT [ "./urlshortner" ]