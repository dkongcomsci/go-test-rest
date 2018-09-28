FROM golang:1.8

USER nobody

RUN mkdir -p /go/src/github.com/dkongcomsci/go-test-rest
WORKDIR /go/src/github.com/dkongcomsci/go-test-rest

COPY . /go/src/github.com/dkongcomsci/go-test-rest
RUN go-wrapper download && go-wrapper install

CMD ["go-wrapper", "run"]
