FROM golang:1.16

WORKDIR /go/src/recipe-count

ARG root_dir
COPY ${root_dir} .

ARG db_dir
COPY ${db_dir} ./${db_dir}

RUN go env -w GO111MODULE=auto
RUN go get -d -v -t ./...
RUN go install -v ./...

CMD ["recipe-count"]
