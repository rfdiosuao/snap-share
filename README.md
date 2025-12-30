# Snap & Share (闪传) ⚡️

[![Go Version](https://img.shields.io/badge/go-1.23%2B-blue)](https://golang.org/)
[![Vue Version](https://img.shields.io/badge/vue-3.x-green)](https://vuejs.org/)
[![License](https://img.shields.io/badge/license-MIT-yellow)](LICENSE)

> **即传即扫，阅后即焚**。一个基于 Go + Vue 3 的极简文件传输系统。
> 这是一个“保姆级”教程，旨在帮助任何人（哪怕是 6 岁的小朋友，如果在父母陪同下）也能成功部署。

---

## 📚 目录

1.  [✨ 核心特性](#-核心特性)
2.  [🛠 第一步：准备工具（环境安装）](#-第一步准备工具环境安装)
3.  [📦 第二步：下载代码](#-第二步下载代码)
4.  [⚙️ 第三步：关键配置（必看！）](#️-第三步关键配置必看)
5.  [🚀 第四步：一键启动（Windows/Linux/Mac）](#-第四步一键启动windowslinuxmac)
6.  [❓ 常见问题 (FAQ)](#-常见问题-faq)

---

## ✨ 核心特性

- **极速传输**：支持大文件拖拽上传，速度飞快。
- **扫码即下**：自动生成二维码，手机扫一扫就能下载。
- **阅后即焚**：文件默认 1 小时后自动删除，或者下载 5 次后自动消失。
- **极简部署**：只需一个端口，就能搞定所有功能。

---

## 🛠 第一步：准备工具（环境安装）

在开始之前，你需要给电脑安装两个“超能力”工具：**Go** 和 **Node.js**。

### 1. 安装 Go 语言 (后端引擎)
*   **下载地址**：[https://go.dev/dl/](https://go.dev/dl/)
*   **安装方法**：下载对应系统的安装包（Windows 选 `.msi`，Mac 选 `.pkg`），一路点击“下一步”直到完成。
*   **验证安装**：打开终端（Windows 按 `Win+R` 输入 `cmd`），输入：
    ```bash
    go version
    ```
    如果显示类似 `go version go1.23.x ...`，说明安装成功！

### 2. 安装 Node.js (前端工厂)
*   **下载地址**：[https://nodejs.org/](https://nodejs.org/)
*   **版本选择**：推荐下载 **LTS (长期支持版)**。
*   **安装方法**：同样是一路“下一步”。
*   **验证安装**：在终端输入：
    ```bash
    node -v
    npm -v
    ```
    如果有版本号跳出来，说明准备就绪！

### 3. 配置 Go 加速器 (中国大陆用户必做 🇨🇳)
为了让下载依赖包像火箭一样快，请在终端执行以下命令：
```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

---

## 📦 第二步：下载代码

如果你会用 Git，直接 Clone；如果不会，直接点击 GitHub 页面上的 **"Code" -> "Download ZIP"**，然后解压到一个文件夹里。

假设你解压到了 `D:\snap-share`。

---

## ⚙️ 第三步：关键配置（必看！）

这是最关键的一步！为了让你的手机能扫描二维码下载，我们需要告诉程序你的电脑 IP 是什么。

1.  进入 `backend` 文件夹。
2.  找到 `config.json` 文件，用记事本或代码编辑器打开。
3.  找到 `"base_url"` 这一行。

**如何查看本机 IP？**
*   **Windows**: 打开终端，输入 `ipconfig`，找到“IPv4 地址”（通常是 `192.168.x.x`）。
*   **Mac/Linux**: 打开终端，输入 `ifconfig` 或 `ip a`。

**修改配置：**
```json
{
  "server": {
    "port": ":8080",
    "base_url": "http://192.168.1.100:8080"  <-- 把这里改成你的 IP！
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
> **注意**：不要忘了后面的 `:8080` 端口号！

---

## 🚀 第四步：一键启动（Windows/Linux/Mac）

我们采用最简单的“生产环境部署”模式，分为三个小动作：

### 动作 1：构建前端页面
打开终端，进入 `frontend` 文件夹：
```bash
cd frontend
npm install   # 安装依赖（第一次运行需要，可能要等几分钟）
npm run build # 开始构建
```
完成后，你会发现 `frontend` 目录下多了一个 `dist` 文件夹。

### 动作 2：搬运素材
将刚刚生成的 `dist` 文件夹，**整个复制**到 `backend` 文件夹里。
确保目录结构是这样的：
```
backend/
├── main.go
├── config.json
├── dist/       <-- 这里面有 index.html
└── ...
```

### 动作 3：启动引擎
打开终端，进入 `backend` 文件夹：
```bash
cd backend
go mod tidy     # 下载 Go 依赖（第一次运行需要）
go run main.go  # 启动！
```

当看到 `Server starting on :8080` 时，恭喜你，成功了！🎉

👉 **打开浏览器访问**：[http://localhost:8080](http://localhost:8080)

---

## ❓ 常见问题 (FAQ)

**Q: 手机扫码后打不开？**
A: 请检查：
1.  手机和电脑是否连接了**同一个 Wi-Fi**？
2.  `config.json` 里的 IP 地址填对了吗？
3.  电脑防火墙是否拦截了 8080 端口？（尝试临时关闭防火墙测试）

**Q: 怎么修改文件保存时间？**
A: 修改 `config.json` 里的 `file_ttl_minutes`（分钟数）。

**Q: 我想限制下载次数？**
A: 修改 `default_download_limit`，设为 `0` 表示不限制。

**Q: 下次启动还要这么麻烦吗？**
A: 不需要！只要不修改前端代码，下次只需要执行 **动作 3**（`go run main.go`）即可。

---

## 📄 开源协议

本项目采用 [MIT License](LICENSE) 开源。
