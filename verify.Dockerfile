ARG PYTHON_VERSION

FROM python:${PYTHON_VERSION}-slim-buster
RUN apt-get update; \
    apt-get install -y --no-install-recommends \
        g++ \
        ; \
    rm -rf /var/lib/apt/lists/*
COPY dist/reckon_v0.0.0-SNAPSHOT-da0554f_linux_amd64.tar.gz .
RUN tar xvf reckon_v0.0.0-SNAPSHOT-da0554f_linux_amd64.tar.gz \
 && mv reckon_v0.0.0-SNAPSHOT-da0554f_linux_amd64 reckon \
 && cd reckon \
 && pip install --no-cache-dir -r ./model/requirements.txt
WORKDIR /reckon
CMD ["./reckon"]
