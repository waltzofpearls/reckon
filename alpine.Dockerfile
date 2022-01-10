ARG PYTHON_VERSION
FROM python:${PYTHON_VERSION}-alpine3.13

RUN apk add -U --no-cache \
        curl \
        gcc \
        g++ \
        make \
        musl-dev \
        zlib-dev \
        jpeg-dev

ARG APP VERSION OS ARCH
RUN curl -L https://github.com/waltzofpearls/${APP}/releases/download/v${VERSION}/${APP}_${VERSION}_${OS}_${ARCH}.tar.gz | tar xvz \
 && mv ${APP}_${VERSION}_${OS}_${ARCH} reckon \
 && cd reckon \
 && pip install --no-cache-dir -r ./model/requirements.txt
RUN pip install pystan==2.19.1.1
RUN pip install prophet==1.0.1

RUN curl -fsSL https://pkgs.tangram.dev/stable/alpine/tangram.rsa | tee /etc/apk/keys/tangram.rsa \
 && echo "https://pkgs.tangram.dev/stable/alpine" | tee /etc/apk/repositories \
 && apk add tangram

WORKDIR /reckon
CMD ["./reckon"]
