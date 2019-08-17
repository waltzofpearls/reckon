FROM golang:1.12-stretch AS builder
WORKDIR /reckon/
COPY . .
RUN make

FROM python:3.6-slim-stretch
WORKDIR /reckon/
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY . .
COPY --from=builder /reckon/reckon .
CMD ["./reckon"]
