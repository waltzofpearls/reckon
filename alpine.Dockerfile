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
ARG APP
ARG VERSION
ARG OS
ARG ARCH
RUN curl -L https://github.com/waltzofpearls/${APP}/releases/download/v${VERSION}/${APP}_${VERSION}_${OS}_${ARCH}.tar.gz | tar xvz \
 && mv ${APP}_${VERSION}_${OS}_${ARCH} reckon \
 && cd reckon \
 && pip install --no-cache-dir -r ./model/requirements.txt
RUN pip install pystan==2.19.1.1
RUN pip install prophet==1.0.1
WORKDIR /reckon
CMD ["./reckon"]
