<script setup>
import { ref } from 'vue'
import QrcodeVue from 'qrcode.vue'
import axios from 'axios'

const file = ref(null)
const uploading = ref(false)
const uploadProgress = ref(0)
const result = ref(null)
const error = ref(null)

const handleFileChange = (e) => {
  const selectedFile = e.target.files[0]
  if (selectedFile) {
    uploadFile(selectedFile)
  }
}

const handleDrop = (e) => {
  e.preventDefault()
  const droppedFile = e.dataTransfer.files[0]
  if (droppedFile) {
    uploadFile(droppedFile)
  }
}

const uploadFile = async (selectedFile) => {
  if (selectedFile.size > 100 * 1024 * 1024) {
    error.value = "文件大小不能超过 100MB"
    return
  }

  file.value = selectedFile
  uploading.value = true
  error.value = null
  result.value = null
  uploadProgress.value = 0

  const formData = new FormData()
  formData.append('file', selectedFile)

  try {
    const response = await axios.post('/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      onUploadProgress: (progressEvent) => {
        uploadProgress.value = Math.round((progressEvent.loaded * 100) / progressEvent.total)
      }
    })
    result.value = response.data
  } catch (err) {
    console.error(err)
    error.value = err.response?.data?.error || "上传失败，请重试"
  } finally {
    uploading.value = false
  }
}

const downloadQrCode = () => {
  const canvas = document.querySelector('.qr-container canvas')
  if (canvas) {
    const link = document.createElement('a')
    link.download = `qrcode-${file.value.name}.png`
    link.href = canvas.toDataURL()
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  }
}

const reset = () => {
  file.value = null
  result.value = null
  error.value = null
  uploadProgress.value = 0
}
</script>

<template>
  <div class="container">
    <header>
      <h1>Snap & Share</h1>
      <p>即传即扫，阅后即焚</p>
    </header>

    <main>
      <!-- 上传区域 -->
      <div 
        v-if="!result" 
        class="upload-zone"
        :class="{ 'dragging': false }" 
        @dragover.prevent
        @drop="handleDrop"
        @click="$refs.fileInput.click()"
      >
        <input 
          type="file" 
          ref="fileInput" 
          style="display: none" 
          @change="handleFileChange"
        >
        
        <div v-if="!uploading" class="upload-placeholder">
          <div class="icon">📁</div>
          <h3>点击或拖拽文件到此处</h3>
          <p>支持最大 100MB 文件</p>
        </div>

        <div v-else class="progress-container">
          <div class="spinner"></div>
          <p>上传中... {{ uploadProgress }}%</p>
        </div>
      </div>

      <!-- 结果展示 -->
      <div v-else class="result-zone">
        <div class="qr-container">
          <qrcode-vue :value="result.download_url" :size="200" level="H" />
        </div>
        
        <div class="info">
          <p class="filename">{{ file.name }}</p>
          <div class="actions">
            <a :href="result.download_url" target="_blank" class="download-btn">直接下载文件</a>
            <button @click="downloadQrCode" class="action-btn">保存二维码</button>
            <button @click="reset" class="reset-btn">上传新文件</button>
          </div>
        </div>
        
        <p class="expire-hint">文件将在 1 小时后自动销毁</p>
      </div>

      <!-- 错误提示 -->
      <div v-if="error" class="error-message">
        {{ error }}
      </div>
    </main>
  </div>
</template>

<style scoped>
.container {
  max-width: 600px;
  margin: 0 auto;
  padding: 2rem;
  font-family: 'Helvetica Neue', Arial, sans-serif;
  text-align: center;
  color: #333;
}

header h1 {
  font-size: 2.5rem;
  margin-bottom: 0.5rem;
  background: linear-gradient(45deg, #42b883, #35495e);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

header p {
  color: #666;
  margin-bottom: 3rem;
}

.upload-zone {
  border: 2px dashed #ddd;
  border-radius: 12px;
  padding: 4rem 2rem;
  cursor: pointer;
  transition: all 0.3s ease;
  background: #f9f9f9;
}

.upload-zone:hover {
  border-color: #42b883;
  background: #f0fdf4;
}

.icon {
  font-size: 4rem;
  margin-bottom: 1rem;
}

.progress-container {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.result-zone {
  animation: fadeIn 0.5s ease;
}

.qr-container {
  background: white;
  padding: 1rem;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  display: inline-block;
  margin-bottom: 2rem;
}

.info {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  align-items: center;
}

.filename {
  font-weight: bold;
  color: #555;
}

.actions {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
  justify-content: center;
}

.download-btn {
  background: #42b883;
  color: white;
  padding: 0.8rem 2rem;
  border-radius: 25px;
  text-decoration: none;
  font-weight: bold;
  transition: transform 0.2s;
}

.download-btn:hover {
  transform: scale(1.05);
}

.action-btn {
  background: #35495e;
  color: white;
  border: none;
  padding: 0.8rem 2rem;
  border-radius: 25px;
  cursor: pointer;
  font-weight: bold;
  transition: transform 0.2s;
}

.action-btn:hover {
  transform: scale(1.05);
  background: #2c3e50;
}

.reset-btn {
  background: transparent;
  border: 1px solid #ddd;
  padding: 0.8rem 2rem;
  border-radius: 25px;
  cursor: pointer;
  color: #666;
  font-weight: bold;
}

.reset-btn:hover {
  border-color: #999;
  color: #333;
}

.expire-hint {
  margin-top: 2rem;
  font-size: 0.8rem;
  color: #999;
}

.error-message {
  margin-top: 1rem;
  color: #e74c3c;
  background: #fdeaea;
  padding: 0.5rem;
  border-radius: 4px;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
