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
- 📦 **单端口部署**：后端自动集成前端静态资源，对外仅需开放一个端口。

## 🏗 架构说明

- **Backend**: Go (Gin Web Framework)
- **Frontend**: Vue 3 + Vite
- **Storage**: 本地文件系统 (Ephemeral Storage)

## 🚀 快速开始

### 1. 环境准备
- Go 1.20+
- Node.js 16+

### 2. 生产环境部署 (推荐)

这是最简单、最稳定的运行方式，适合局域网或服务器部署。

#### 第一步：构建前端
```bash
cd frontend
npm install
npm run build
```
构建完成后，会在 `frontend` 目录下生成一个 `dist` 文件夹。

#### 第二步：整合资源
将生成的 `dist` 文件夹**完整复制**到 `backend` 目录下。
目录结构应如下所示：
```
backend/
├── main.go
├── config.json
├── uploads/
└── dist/       <-- 前端构建产物
    ├── index.html
    └── assets/
```

#### 第三步：启动服务
```bash
cd backend
go mod tidy
go run main.go
# 或者编译后运行: go build -o server && ./server
```

此时，访问 `http://localhost:8080` 即可看到完整应用。

### 3. 开发环境运行 (调试用)

如果你需要修改代码，可以分别启动前后端：

*   **后端**: `cd backend && go run main.go` (运行在 :8080)
*   **前端**: `cd frontend && npm run dev` (运行在 :5173，自动代理 API 到后端)

## ⚙️ 首次启动配置 (Configuration)

在 `backend` 目录下创建或修改 `config.json` 文件。

> ⚠️ **重要提示**：为了让手机能扫描二维码下载，你**必须**修改 `base_url`。

```json
{
  "server": {
    "port": ":8080",
    "base_url": "http://192.168.1.100:8080"  // <--- 修改这里！
  },
  "storage": {
    "upload_dir": "./uploads",
    "max_file_size_mb": 100,
    "file_ttl_minutes": 60,
    "default_download_limit": 5,
    "static_dir": "./dist"
  }
}
```

### 关键配置项说明

| 字段 | 说明 | 推荐值 |
| :--- | :--- | :--- |
| `server.base_url` | **核心配置**。生成二维码时使用的基础 URL。**必须修改为局域网 IP 或公网域名**，否则手机扫码后无法访问。 | `http://<你的IP>:8080` |
| `server.port` | 后端监听端口。 | `:8080` |
| `storage.static_dir` | 前端静态文件目录。如果配置为空，则仅作为 API 服务器运行。 | `./dist` |
| `storage.default_download_limit` | 最大下载次数。超过次数后文件立即销毁。设为 0 则无限制。 | `5` |
| `storage.file_ttl_minutes` | 文件自动过期时间 (分钟)。 | `60` |

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

## 📄 开源协议

本项目采用 [MIT License](LICENSE) 开源。
