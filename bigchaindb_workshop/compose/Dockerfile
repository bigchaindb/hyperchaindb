FROM python:3

RUN mkdir -p /usr/src/app/bigchaindb_workshop
WORKDIR /usr/src/app/

COPY setup.py /usr/src/app/
COPY backend /usr/src/app/bigchaindb_workshop/backend
COPY exercises /usr/src/app/bigchaindb_workshop/exercises

RUN pip install --upgrade pip
RUN pip install --no-cache-dir -e .[dev]
