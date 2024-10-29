# 1. 빌드 단계
FROM golang:1.23.2 AS builder

ARG APP_DIR
ARG APP_PORT
ARG DEBIAN_FRONTEND=noninteractive

ENV TZ=Asia/Seoul

EXPOSE ${APP_PORT}

# 작업 디렉토리 설정
WORKDIR ${APP_DIR}

# go.mod와 go.sum을 복사하여 의존성 설치
COPY go.mod go.sum ./
RUN go mod download

# 소스 코드를 복사하고 빌드
COPY . .
RUN go build -o app -buildvcs=false

# 애플리케이션 실행
CMD sh -c "/app/app || (echo 'Application failed, keeping container alive' && sleep infinity)"
