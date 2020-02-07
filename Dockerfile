FROM golang:1.11-stretch

RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 

RUN go get -u "golang.org/x/net/context"
RUN go get -u "golang.org/x/oauth2"
RUN go get -u "golang.org/x/oauth2/google"
RUN go get -u "google.golang.org/api/gmail/v1"
RUN go get -u "github.com/sendgrid/sendgrid-go"


RUN go build ./main.go

EXPOSE 8080

CMD ["/app/main"]