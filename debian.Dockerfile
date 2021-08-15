ARG PYTHON_VERSION

FROM python:${PYTHON_VERSION}-slim-buster
RUN apt-get update; \
    apt-get install -y --no-install-recommends \
        curl \
        g++ \
        ; \
    rm -rf /var/lib/apt/lists/*
ARG APP
ARG VERSION
ARG OS
ARG ARCH
RUN curl -L https://github.com/waltzofpearls/${APP}/releases/download/v${VERSION}/${APP}_${VERSION}_${OS}_${ARCH}.tar.gz | tar xvz \
 && mv ${APP}_${VERSION}_${OS}_${ARCH} reckon \
 && cd reckon \
 && pip install --no-cache-dir -r ./model/requirements.txt
WORKDIR /reckon
CMD ["./reckon"]
