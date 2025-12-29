# Snap & Share (闪传) ⚡️

[![Go Version](https://img.shields.io/badge/go-1.23%2B-blue)](https://golang.org/)
[![Vue Version](https://img.shields.io/badge/vue-3.x-green)](https://vuejs.org/)
[![License](https://img.shields.io/badge/license-MIT-yellow)](LICENSE)

> **即传即扫，阅后即焚**。一个基于 Go + Vue 3 的极简文件传输系统。
> Zero-config file sharing for local networks.

## ✨ 核心特性

- 🚀 **极速传输**：基于 Go Gin 框架的高性能后端，支持大文件流式上传。
- 📱 **扫码即下**：自动生成二维码，手机无需安装 App 即可下载。
- 🔒 **阅后即焚**：
  - **时间限制**：默认 1 小时后自动销毁文件。
  - **次数限制**：支持配置下载次数（如：限制 5 次下载后自动删除）。
- 🎨 **极简体验**：Vue 3 打造的丝滑拖拽上传界面。
- 🛠 **高度可配**：通过配置文件自定义端口、存储路径、过期策略。

## 🏗 架构说明

- **Backend**: Go (Gin Web Framework)
- **Frontend**: Vue 3 + Vite
- **Storage**: 本地文件系统 (Ephemeral Storage)

## 🚀 快速开始

### 1. 环境准备
- Go 1.20+
- Node.js 16+

### 2. 启动后端

```bash
cd backend
go mod tidy
go run main.go
```

后端服务将在 `http://localhost:8080` 启动。

### 3. 启动前端

```bash
cd frontend
npm install
npm run dev
```

前端页面将在 `http://localhost:5173` 启动。

## ⚙️ 配置文件 (config.json)

在 `backend` 目录下创建 `config.json` 文件进行配置：

```json
{
  "server": {
    "port": ":8080",
    "base_url": "http://localhost:8080"
  },
  "storage": {
    "upload_dir": "./uploads",
    "max_file_size_mb": 100,
    "file_ttl_minutes": 60,
    "default_download_limit": 5
  }
}
```

| 字段 | 说明 | 默认值 |
| :--- | :--- | :--- |
| `port` | 后端监听端口 | `:8080` |
| `base_url` | 生成二维码的基础 URL (局域网请填 IP) | `http://localhost:8080` |
| `max_file_size_mb` | 最大文件大小 (MB) | `100` |
| `file_ttl_minutes` | 文件过期时间 (分钟) | `60` |
| `default_download_limit` | 最大下载次数 (0 为无限) | `5` |

## 📦 部署指南

### 编译后端
```bash
cd backend
go build -o snap-share-server main.go
```

### 构建前端
```bash
cd frontend
npm run build
```

将构建好的前端静态文件 (`dist` 目录) 部署到 Nginx 或集成到 Go 后端中即可。

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

## 📄 开源协议

本项目采用 [MIT License](LICENSE) 开源。
