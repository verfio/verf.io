FROM golang:1.11-stretch

RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 


RUN go build ./main.go

EXPOSE 8080

CMD ["/app/main"]