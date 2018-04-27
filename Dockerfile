FROM golang:alpine AS build-env
ADD . /go/src/github.com/hibooboo2/strawpoll/
RUN ls /go/src/github.com/hibooboo2/strawpoll/
RUN cd /go/src/github.com/hibooboo2/strawpoll/ && go build -o goapp

FROM alpine
EXPOSE 8080
WORKDIR /app
COPY --from=build-env /go/src/github.com/hibooboo2/strawpoll/goapp /app/
ENTRYPOINT ./goapp
