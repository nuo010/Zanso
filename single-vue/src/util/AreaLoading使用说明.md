# 区域Loading使用说明

## 概述
区域Loading组件用于在`<el-main>`区域内显示加载状态，而不是全屏loading。提供了灵活的手动控制方法。

## 基本用法

### 1. 导入方法
```typescript
import { 
  showAreaLoading, 
  hideAreaLoading, 
  forceHideAreaLoading, 
  setAreaLoadingText,
  isAreaLoading,
  areaLoadingText
} from '@/util/routeLoading'
```

### 2. 基本控制方法

#### 显示Loading
```typescript
// 显示默认loading
showAreaLoading()

// 显示带自定义文本的loading
showAreaLoading('正在加载数据...')
```

#### 隐藏Loading
```typescript
// 隐藏loading（支持计数，需要调用相同次数才会隐藏）
hideAreaLoading()

// 强制隐藏loading（立即隐藏，重置计数）
forceHideAreaLoading()
```

#### 更改Loading文本
```typescript
// 更改当前显示的loading文本
setAreaLoadingText('新的加载文本...')
```

### 3. 状态获取
```typescript
// 获取当前loading状态
console.log(isAreaLoading.value) // true/false

// 获取当前loading文本
console.log(areaLoadingText.value) // 当前显示的文本
```

## 使用场景示例

### 1. 异步操作
```typescript
const handleAsyncOperation = async () => {
  showAreaLoading('正在处理数据...')
  
  try {
    const result = await someAsyncFunction()
    console.log('操作完成:', result)
  } catch (error) {
    console.error('操作失败:', error)
  } finally {
    hideAreaLoading()
  }
}
```

### 2. API请求
```typescript
const fetchData = async () => {
  showAreaLoading('正在请求数据...')
  
  try {
    const response = await fetch('/api/data')
    setAreaLoadingText('数据处理中...')
    
    const data = await response.json()
    return data
  } finally {
    hideAreaLoading()
  }
}
```

### 3. 多个并发操作
```typescript
const handleMultipleOperations = async () => {
  // 第一个操作
  showAreaLoading('正在执行操作1...')
  await operation1()
  
  // 第二个操作（会更新文本，但不会重复显示loading）
  setAreaLoadingText('正在执行操作2...')
  await operation2()
  
  // 第三个操作
  setAreaLoadingText('正在执行操作3...')
  await operation3()
  
  // 所有操作完成
  hideAreaLoading()
}
```

### 4. 带超时的操作
```typescript
const handleWithTimeout = async () => {
  showAreaLoading('正在处理，请稍候...')
  
  const timeout = setTimeout(() => {
    setAreaLoadingText('处理时间较长，请耐心等待...')
  }, 5000)
  
  try {
    await longRunningOperation()
  } finally {
    clearTimeout(timeout)
    hideAreaLoading()
  }
}
```

## 高级用法

### 1. 使用管理实例
```typescript
import { areaLoadingManager } from '@/util/routeLoading'

// 直接使用管理实例
areaLoadingManager.show('自定义文本')
areaLoadingManager.hide()
areaLoadingManager.setText('新文本')
console.log(areaLoadingManager.isLoading)
```

### 2. 在Vue组件中使用
```vue
<template>
  <div>
    <el-button @click="startLoading">开始加载</el-button>
    <el-button @click="stopLoading">停止加载</el-button>
  </div>
</template>

<script setup>
import { showAreaLoading, hideAreaLoading } from '@/util/routeLoading'

const startLoading = () => {
  showAreaLoading('组件中的loading...')
}

const stopLoading = () => {
  hideAreaLoading()
}
</script>
```

## 注意事项

1. **计数机制**: `showAreaLoading()`和`hideAreaLoading()`使用计数机制，调用多少次show就需要调用多少次hide才会真正隐藏loading。

2. **强制隐藏**: 如果需要立即隐藏loading，使用`forceHideAreaLoading()`方法。

3. **文本更新**: 在loading显示期间，可以随时使用`setAreaLoadingText()`更新显示文本。

4. **状态响应**: `isAreaLoading`和`areaLoadingText`是响应式的，可以在Vue组件中直接使用。

5. **性能考虑**: loading组件使用了绝对定位，不会影响页面布局，性能开销很小。

## 样式自定义

如果需要自定义loading样式，可以修改`src/components/AreaLoading/AreaLoading.vue`文件中的CSS样式。

主要样式类：
- `.area-loading`: 主容器
- `.area-loading-mask`: 遮罩层
- `.area-loading-content`: 内容区域
- `.spinner-circle`: 旋转动画
- `.loading-text`: 文本样式
