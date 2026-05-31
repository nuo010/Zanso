// @/util/routeLoading.ts
import { ref } from 'vue'

// 全局loading状态
const isRouteLoading = ref(false)
const loadingText = ref('页面加载中...')

// 创建loading管理类
class RouteLoadingManager {
  private static instance: RouteLoadingManager
  private loadingCount = 0

  static getInstance() {
    if (!RouteLoadingManager.instance) {
      RouteLoadingManager.instance = new RouteLoadingManager()
    }
    return RouteLoadingManager.instance
  }

  // 显示loading
  show(text?: string) {
    this.loadingCount++
    isRouteLoading.value = true
    if (text) {
      loadingText.value = text
    }
  }

  // 隐藏loading
  hide() {
    if (this.loadingCount > 0) {
      this.loadingCount--
    }
    if (this.loadingCount === 0) {
      isRouteLoading.value = false
    }
  }

  // 强制隐藏loading（重置计数）
  forceHide() {
    this.loadingCount = 0
    isRouteLoading.value = false
  }

  // 设置loading文本
  setText(text: string) {
    loadingText.value = text
  }

  // 获取当前状态
  get isLoading() {
    return isRouteLoading.value
  }

  get text() {
    return loadingText.value
  }
}

// 创建全局实例
const loadingManager = RouteLoadingManager.getInstance()

// 导出便捷方法
export const showAreaLoading = (text?: string) => {
  loadingManager.show(text)
}

export const hideAreaLoading = () => {
  loadingManager.hide()
}

export const forceHideAreaLoading = () => {
  loadingManager.forceHide()
}

export const setAreaLoadingText = (text: string) => {
  loadingManager.setText(text)
}

// 导出状态
export const isAreaLoading = isRouteLoading
export const areaLoadingText = loadingText

// 导出管理实例
export const areaLoadingManager = loadingManager

// 兼容旧的useRouteLoading方法
export const useRouteLoading = () => {
  return {
    isRouteLoading,
    loadingText,
    startLoading: showAreaLoading,
    stopLoading: hideAreaLoading,
    setLoadingText: setAreaLoadingText
  }
}
