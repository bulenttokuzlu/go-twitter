FROM golang:latest 
RUN mkdir /app 
RUN mkdir /log
ADD . /app/ 
WORKDIR /app
RUN go build streaming.go  
CMD ["/app/streaming"]