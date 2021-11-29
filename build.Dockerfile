ARG PYTHON_VERSION
ARG GO_VERSION

FROM python:${PYTHON_VERSION}-slim-buster

RUN apt-get update; \
    apt-get install -y --no-install-recommends \
        g++ \
        git \
        curl \
        make \
        pkg-config \
        gnupg \
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
RUN pip install pystan==2.19.1.1
RUN pip install prophet==1.0.1

RUN curl -fsSL https://pkgs.tangram.dev/stable/debian/buster.gpg | apt-key add - \
 && curl -fsSL https://pkgs.tangram.dev/stable/debian/buster.list | tee /etc/apt/sources.list.d/tangram.list \
 && apt-get update \
 && apt-get install tangram

COPY . .
RUN make
CMD ["./reckon"]
