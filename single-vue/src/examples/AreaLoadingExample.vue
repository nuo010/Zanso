<!-- @/examples/AreaLoadingExample.vue -->
<template>
  <div class="loading-example">
    <h2>区域Loading使用示例</h2>
    
    <div class="button-group">
      <el-button type="primary" @click="showDefaultLoading">
        显示默认Loading
      </el-button>
      
      <el-button type="success" @click="showCustomLoading">
        显示自定义文本Loading
      </el-button>
      
      <el-button type="warning" @click="showMultipleLoading">
        显示多个Loading（测试计数）
      </el-button>
      
      <el-button type="info" @click="changeLoadingText">
        更改Loading文本
      </el-button>
      
      <el-button type="danger" @click="hideLoading">
        隐藏Loading
      </el-button>
      
      <el-button @click="forceHideLoading">
        强制隐藏Loading
      </el-button>
    </div>

    <div class="demo-section">
      <h3>模拟异步操作</h3>
      <el-button type="primary" @click="simulateAsyncOperation">
        模拟异步操作（3秒）
      </el-button>
    </div>

    <div class="demo-section">
      <h3>模拟API请求</h3>
      <el-button type="success" @click="simulateApiRequest">
        模拟API请求
      </el-button>
    </div>

    <div class="status-info">
      <p>当前Loading状态: <span :class="{ 'loading-active': isAreaLoading }">{{ isAreaLoading ? '显示中' : '隐藏中' }}</span></p>
      <p>当前Loading文本: {{ areaLoadingText }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { 
  showAreaLoading, 
  hideAreaLoading, 
  forceHideAreaLoading, 
  setAreaLoadingText,
  isAreaLoading,
  areaLoadingText
} from '@/util/routeLoading'

// 模拟异步操作
const simulateAsyncOperation = async () => {
  showAreaLoading('正在处理数据...')
  
  try {
    // 模拟3秒的异步操作
    await new Promise(resolve => setTimeout(resolve, 3000))
    console.log('异步操作完成')
  } finally {
    hideAreaLoading()
  }
}

// 模拟API请求
const simulateApiRequest = async () => {
  showAreaLoading('正在请求数据...')
  
  try {
    // 模拟API请求
    const response = await fetch('/api/data')
    if (response.ok) {
      setAreaLoadingText('数据处理中...')
      await new Promise(resolve => setTimeout(resolve, 1000))
      console.log('API请求成功')
    }
  } catch (error) {
    setAreaLoadingText('请求失败，重试中...')
    await new Promise(resolve => setTimeout(resolve, 1000))
    console.error('API请求失败:', error)
  } finally {
    hideAreaLoading()
  }
}

// 显示默认Loading
const showDefaultLoading = () => {
  showAreaLoading()
}

// 显示自定义文本Loading
const showCustomLoading = () => {
  showAreaLoading('自定义加载文本...')
}

// 显示多个Loading（测试计数功能）
const showMultipleLoading = () => {
  showAreaLoading('第一个Loading')
  setTimeout(() => {
    showAreaLoading('第二个Loading')
  }, 500)
  setTimeout(() => {
    showAreaLoading('第三个Loading')
  }, 1000)
}

// 更改Loading文本
const changeLoadingText = () => {
  setAreaLoadingText('文本已更改: ' + new Date().toLocaleTimeString())
}

// 隐藏Loading
const hideLoading = () => {
  hideAreaLoading()
}

// 强制隐藏Loading
const forceHideLoading = () => {
  forceHideAreaLoading()
}
</script>

<style scoped>
.loading-example {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.button-group {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin: 20px 0;
}

.demo-section {
  margin: 30px 0;
  padding: 20px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  background: #fafafa;
}

.status-info {
  margin-top: 30px;
  padding: 15px;
  background: #f0f9ff;
  border: 1px solid #b3d8ff;
  border-radius: 8px;
}

.status-info p {
  margin: 5px 0;
}

.loading-active {
  color: #409eff;
  font-weight: bold;
}
</style>
