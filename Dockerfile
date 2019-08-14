FROM golang:1.12-alpine AS builder
WORKDIR /go/src/github.com/waltzofpearls/reckon/
RUN apk --no-cache add make
COPY . .
RUN make

FROM python:3.6-slim-stretch
WORKDIR /reckon/
COPY requirements.txt .
RUN apt-get update \
 && apt-get install -y gcc g++ \
 && pip install pystan \
 && pip install -r requirements.txt
COPY . .
COPY --from=builder /go/src/github.com/waltzofpearls/reckon/reckon .
CMD ["./reckon"]
