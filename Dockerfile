FROM public.ecr.aws/docker/library/golang:alpine3.18 AS builder

ENV SRC=/go/src/github.com/mdwn/ghstatus
RUN mkdir -p $SRC
WORKDIR $SRC
COPY . .
RUN mkdir -p ./out
RUN go build -o ./out ./...

FROM public.ecr.aws/docker/library/alpine:3.18.0 AS runner

COPY --from=builder /go/src/github.com/mdwn/ghstatus/out/ghstatus /opt/ghstatus/ghstatus
ENTRYPOINT /opt/ghstatus/ghstatus
