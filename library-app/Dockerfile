FROM golang:1.17
LABEL PROJECT="library_application"

COPY . /go/src/app/
WORKDIR /go/src/app/

RUN go get -d

RUN go install

ENTRYPOINT ["/bin/bash"]