## Overview

This project is a Go service that exposes a Teams Notification API. It follows a conventional Golang mono repo

## Layout

```CSS
repo-root/
├── cmd/                  # 各個可執行服務的入口
│   ├── service-a/
│   │   └── main.go
│   ├── service-b/
│   │   └── main.go
│   └── worker-job/
│       └── main.go
│
├── internal/             # 內部共用，不對外公開
│   ├── pkg1/
│   ├── pkg2/
│   └── platform/         # 基礎設施（db, cache, kafka, tracing...）
│
├── pkg/                  # 可重用、對外公開的 library
│   ├── logger/
│   ├── middleware/
│   └── auth/
│
├── api/                  # API 定義（gRPC proto / OpenAPI spec）
│   ├── proto/
│   └── openapi/
│
├── configs/              # 設定檔（YAML/JSON/env）
│   ├── service-a.yaml
│   ├── service-b.yaml
│   └── common.yaml
│
├── deployments/          # K8s YAML / Helm Charts / Terraform
│   ├── service-a/
│   ├── service-b/
│   └── infra/
│
├── scripts/              # CI/CD 腳本、工具
│   ├── build.sh
│   ├── lint.sh
│   └── gen.sh
│
├── test/                 # 整合測試 / e2e 測試
│   ├── service-a/
│   └── service-b/
│
├── go.mod
├── go.sum
└── Makefile
```

