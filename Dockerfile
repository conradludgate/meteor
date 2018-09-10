FROM golang:latest
RUN go get -v github.com/conradludgate/meteor
RUN go install github.com/conradludgate/meteor
RUN meteor