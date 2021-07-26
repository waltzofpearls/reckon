ARG PYTHON_VERSION
ARG GO_VERSION

FROM python:${PYTHON_VERSION}-slim-buster
RUN apt-get update; \
    apt-get install -y --no-install-recommends \
        g++ \
        curl \
        make \
        pkg-config \
        ; \
    rm -rf /var/lib/apt/lists/*
ARG GO_VERSION
RUN curl -O https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz \
 && tar xvf go${GO_VERSION}.linux-amd64.tar.gz \
 && mv go /usr/local
ENV GOROOT=/usr/local/go
ENV PATH=$GOROOT/bin:$PATH
WORKDIR /reckon/
COPY ./model/requirements.txt ./model/requirements.txt
RUN pip install --no-cache-dir -r ./model/requirements.txt
COPY . .
RUN make
CMD ["./reckon"]
