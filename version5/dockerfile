FROM golang:1.8

# ENV GOPATH $HOME/workspace/go/
# RUN go env
# ADD . .
# COPY . /go
# RUN cd go && go build -o myapp .

# CMD ["myapp"]  
WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]