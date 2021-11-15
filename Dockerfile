FROM python:3.9.6-slim-bullseye

# Setup debian
RUN apt-get update && \
    # Install dependencies
    apt-get install -y --no-install-recommends \
    curl \
    wget \
    ca-certificates \
    git \
    # Clean cache
    && apt-get clean && apt-get -y autoremove && rm -rf /var/lib/apt/lists/*

WORKDIR /tmp/dev

RUN curl https://getmic.ro | bash && \
    mv micro /bin && \
    rm -rf /tmp/dev

WORKDIR /workspace

# Setup development dependencies
WORKDIR /development

COPY ./scripts/requirements.txt ./scripts/
RUN pip install --no-cache-dir -r ./scripts/requirements.txt

COPY scripts ./scripts
RUN invoke tool.install --include dev --yes

# Keep container alive
CMD exec /bin/bash -c "trap : TERM INT; sleep infinity & wait"
