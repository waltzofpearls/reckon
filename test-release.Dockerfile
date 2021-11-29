ARG PYTHON_VERSION

FROM python:${PYTHON_VERSION}-slim-buster

RUN apt-get update; \
    apt-get install -y --no-install-recommends \
        curl \
        g++ \
        gnupg \
        ; \
    rm -rf /var/lib/apt/lists/*

ARG APP
ARG VERSION
ARG OS
ARG ARCH

ADD ./model/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
RUN pip install pystan==2.19.1.1
RUN pip install prophet==1.0.1

RUN curl -fsSL https://pkgs.tangram.dev/stable/debian/buster.gpg | apt-key add - \
 && curl -fsSL https://pkgs.tangram.dev/stable/debian/buster.list | tee /etc/apt/sources.list.d/tangram.list \
 && apt-get update \
 && apt-get install tangram

ADD dist/${APP}_${VERSION}_${OS}_${ARCH}.tar.gz .
RUN mv ${APP}_${VERSION}_${OS}_${ARCH} reckon
WORKDIR /reckon
CMD ["./reckon"]
