FROM golang:1.12-alpine AS builder
WORKDIR /reckon/
RUN apk --no-cache add make
COPY . .
RUN make

FROM python:3.6-slim-stretch
WORKDIR /reckon/
COPY requirements.txt .
RUN pip install -r requirements.txt
COPY . .
COPY --from=builder /reckon/reckon .
CMD ["./reckon"]
