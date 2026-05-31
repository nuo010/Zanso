<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

// 个人信息数据
const userInfo = reactive({
  name: '张三',
  avatar: 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png',
  userType: '内部用户',
  userAccount: 'zhangsan001',
  organization: '广州交投公路建设有限公司',
  department: '技术部',
  position: '高级工程师',
  phone: '138****8888',
  email: 'zhangsan@company.com',
  roles: ['系统管理员', '项目负责人', '技术专家'],
  lastLoginTime: '2024-01-15 09:30:25'
})

// 应用中心数据
const appCenter = ref([
  {
    id: 1,
    name: '系统监控',
    icon: 'Monitor',
    description: '实时监控系统运行状态',
    color: '#409EFF',
    url: '/monitor'
  },
  {
    id: 2,
    name: '用户管理',
    icon: 'User',
    description: '管理系统用户和权限',
    color: '#67C23A',
    url: '/user'
  },
  {
    id: 3,
    name: '数据分析',
    icon: 'TrendCharts',
    description: '数据统计和分析报表',
    color: '#E6A23C',
    url: '/analytics'
  },
  {
    id: 4,
    name: '文件管理',
    icon: 'Folder',
    description: '文件上传下载和管理',
    color: '#F56C6C',
    url: '/files'
  },
  {
    id: 5,
    name: '消息中心',
    icon: 'Message',
    description: '系统消息和通知管理',
    color: '#909399',
    url: '/messages'
  },
  {
    id: 6,
    name: '设置中心',
    icon: 'Setting',
    description: '系统配置和参数设置',
    color: '#606266',
    url: '/settings'
  }
])

// 业务中心数据
const businessCenter = ref([
  {
    id: 1,
    name: '项目管理',
    icon: 'Document',
    description: '项目创建、分配和跟踪',
    color: '#409EFF',
    count: 12,
    url: '/projects'
  },
  {
    id: 2,
    name: '任务管理',
    icon: 'List',
    description: '任务分配、进度跟踪',
    color: '#67C23A',
    count: 28,
    url: '/tasks'
  },
  {
    id: 3,
    name: '客户管理',
    icon: 'UserFilled',
    description: '客户信息维护和跟进',
    color: '#E6A23C',
    count: 156,
    url: '/customers'
  },
  {
    id: 4,
    name: '财务管理',
    icon: 'Money',
    description: '财务数据统计和分析',
    color: '#F56C6C',
    count: 8,
    url: '/finance'
  },
  {
    id: 5,
    name: '报表中心',
    icon: 'DataAnalysis',
    description: '各类业务报表生成',
    color: '#909399',
    count: 24,
    url: '/reports'
  },
  {
    id: 6,
    name: '审批流程',
    icon: 'Check',
    description: '各类审批流程管理',
    color: '#606266',
    count: 15,
    url: '/approval'
  }
])

// 消息通知
const messages = ref([
  {
    id: 1,
    type: 'system',
    title: '系统维护通知',
    content: '系统将于今晚22:00-24:00进行维护升级',
    time: '2小时前',
    isRead: false,
    priority: 'high'
  },
  {
    id: 2,
    type: 'app',
    title: '新版本发布',
    content: '项目管理模块V2.1版本已发布',
    time: '昨天',
    isRead: false,
    priority: 'medium'
  },
  {
    id: 3,
    type: 'system',
    title: '安全提醒',
    content: '检测到异常登录行为，请及时修改密码',
    time: '2天前',
    isRead: true,
    priority: 'high'
  },
  {
    id: 4,
    type: 'app',
    title: '功能更新',
    content: '数据分析模块新增了实时报表功能',
    time: '3天前',
    isRead: true,
    priority: 'low'
  }
])

// 弹窗控制
const businessDialogVisible = ref(false)

// 用户权限和收藏
const userPermissions = ref(['项目管理', '任务管理', '客户管理', '财务管理', '报表中心', '审批流程', '人力资源', '采购管理'])
const userFavorites = ref(['项目管理', '任务管理', '客户管理'])

// 所有业务系统数据
const allBusinessSystems = ref([
  {
    id: 1,
    name: '项目管理',
    description: '项目创建、分配和跟踪',
    icon: 'Document',
    color: '#409EFF',
    status: 'online',
    tags: ['核心', '项目管理']
  },
  {
    id: 2,
    name: '任务管理',
    description: '任务分配、进度跟踪',
    icon: 'List',
    color: '#67C23A',
    status: 'online',
    tags: ['核心', '任务']
  },
  {
    id: 3,
    name: '客户管理',
    description: '客户信息维护和跟进',
    icon: 'User',
    color: '#E6A23C',
    status: 'online',
    tags: ['核心', '客户']
  },
  {
    id: 4,
    name: '财务管理',
    description: '财务数据统计和分析',
    icon: 'TrendCharts',
    color: '#F56C6C',
    status: 'online',
    tags: ['核心', '财务']
  },
  {
    id: 5,
    name: '报表中心',
    description: '各类业务报表生成',
    icon: 'DataAnalysis',
    color: '#909399',
    status: 'online',
    tags: ['报表', '分析']
  },
  {
    id: 6,
    name: '审批流程',
    description: '各类审批流程管理',
    icon: 'CircleCheck',
    color: '#606266',
    status: 'online',
    tags: ['流程', '审批']
  },
  {
    id: 7,
    name: '人力资源',
    description: '员工信息管理和考勤',
    icon: 'UserFilled',
    color: '#409EFF',
    status: 'online',
    tags: ['人事', '考勤']
  },
  {
    id: 8,
    name: '采购管理',
    description: '采购流程和供应商管理',
    icon: 'ShoppingCart',
    color: '#67C23A',
    status: 'online',
    tags: ['采购', '供应商']
  },
  {
    id: 9,
    name: '库存管理',
    description: '库存监控和出入库管理',
    icon: 'Box',
    color: '#E6A23C',
    status: 'maintenance',
    tags: ['库存', '仓储']
  },
  {
    id: 10,
    name: '销售管理',
    description: '销售订单和客户关系管理',
    icon: 'TrendCharts',
    color: '#F56C6C',
    status: 'online',
    tags: ['销售', '订单']
  },
  {
    id: 11,
    name: '质量管理',
    description: '质量检测和标准管理',
    icon: 'Medal',
    color: '#909399',
    status: 'online',
    tags: ['质量', '检测']
  },
  {
    id: 12,
    name: '设备管理',
    description: '设备维护和保养管理',
    icon: 'Setting',
    color: '#606266',
    status: 'offline',
    tags: ['设备', '维护']
  }
])

// 处理应用点击
const handleAppClick = (app: any) => {
  console.log('点击应用:', app.name)
  // 这里可以添加路由跳转逻辑
}

// 处理业务模块点击
const handleBusinessClick = (business: any) => {
  console.log('点击业务模块:', business.name)
  // 这里可以添加路由跳转逻辑
}

// 处理业务中心查看全部点击
const handleBusinessViewAll = () => {
  businessDialogVisible.value = true
}

// 处理弹窗关闭
const handleCloseBusinessDialog = () => {
  businessDialogVisible.value = false
}

// 检查用户是否有权限
const hasPermission = (systemName: string) => {
  return userPermissions.value.includes(systemName)
}

// 检查是否已收藏
const isFavorite = (systemName: string) => {
  return userFavorites.value.includes(systemName)
}

// 切换收藏状态
const toggleFavorite = (systemName: string) => {
  const index = userFavorites.value.indexOf(systemName)
  if (index > -1) {
    userFavorites.value.splice(index, 1)
    ElMessage.success('已取消收藏')
  } else {
    userFavorites.value.push(systemName)
    ElMessage.success('已添加到收藏')
  }
}

// 申请权限
const applyPermission = (systemName: string) => {
  ElMessageBox.confirm(
    `您确定要申请 ${systemName} 系统的访问权限吗？`,
    '权限申请',
    {
      confirmButtonText: '确定申请',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(() => {
    ElMessage.success('权限申请已提交，请等待管理员审核')
  }).catch(() => {
    ElMessage.info('已取消申请')
  })
}

// 处理系统点击
const handleSystemClick = (system: any) => {
  if (hasPermission(system.name)) {
    console.log('进入系统:', system.name)
    // 这里可以添加路由跳转逻辑
    businessDialogVisible.value = false
  } else {
    applyPermission(system.name)
  }
}

// 处理消息点击
const handleMessageClick = (message: any) => {
  console.log('点击消息:', message.title)
  // 标记为已读
  message.isRead = true
  // 这里可以添加查看详情的逻辑
}

// 获取未读消息数量
const unreadMessageCount = computed(() => {
  return messages.value.filter(m => !m.isRead).length
})

// 获取消息类型图标
const getMessageIcon = (type: string) => {
  return type === 'system' ? 'Bell' : 'Message'
}

// 获取优先级颜色
const getPriorityColor = (priority: string) => {
  switch (priority) {
    case 'high': return '#F56C6C'
    case 'medium': return '#E6A23C'
    case 'low': return '#909399'
    default: return '#409EFF'
  }
}

</script>

<template>
  <div class="portal-container">
    <!-- 顶部导航栏 -->
    <div class="top-navbar">
      <div class="navbar-content">
        <div class="logo-section">
          <div class="logo">
            <el-icon size="32"><Platform /></el-icon>
            <span class="logo-text">企业门户</span>
          </div>
        </div>
        <div class="navbar-right">
          <div class="weather-info">
            <el-icon size="20"><Sunny /></el-icon>
            <span>北京 晴 22°C</span>
          </div>
          <div class="user-profile">
            <el-avatar :size="40" :src="userInfo.avatar" />
            <div class="user-details">
              <span class="user-name">{{ userInfo.name }}</span>
              <span class="user-role">{{ userInfo.position }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="main-content">

      <!-- 业务中心 -->
      <div class="section-card business-section">
        <div class="section-header">
          <div class="header-left">
            <div class="section-icon business-icon">
              <el-icon><Briefcase /></el-icon>
            </div>
            <div>
              <h3>业务中心</h3>
              <p>核心业务功能模块</p>
            </div>
          </div>
          <el-button type="primary" size="small" class="more-btn" @click="handleBusinessViewAll">
            查看全部
            <el-icon><ArrowRight /></el-icon>
          </el-button>
        </div>
        <div class="business-grid">
          <div 
            v-for="business in businessCenter" 
            :key="business.id"
            class="business-card"
            @click="handleBusinessClick(business)"
          >
            <div class="business-icon-wrapper" :style="{ background: `linear-gradient(135deg, ${business.color} 0%, ${business.color}CC 100%)` }">
              <el-icon size="16"><component :is="business.icon" /></el-icon>
            </div>
            <div class="business-info">
              <h4>{{ business.name }}</h4>
              <p>{{ business.description }}</p>
            </div>
            <div class="business-arrow">
              <el-icon><ArrowRight /></el-icon>
            </div>
          </div>
        </div>
      </div>

      <!-- 应用中心 -->
      <div class="section-card app-section">
        <div class="section-header">
          <div class="header-left">
            <div class="section-icon app-icon">
              <el-icon><Grid /></el-icon>
            </div>
            <div>
              <h3>应用中心</h3>
              <p>快速访问常用应用</p>
            </div>
          </div>
          <el-button type="primary" size="small" class="more-btn">
            查看全部
            <el-icon><ArrowRight /></el-icon>
          </el-button>
        </div>
        <div class="app-grid">
          <div 
            v-for="app in appCenter" 
            :key="app.id"
            class="app-card"
            @click="handleAppClick(app)"
          >
            <div class="app-icon-wrapper" :style="{ background: `linear-gradient(135deg, ${app.color} 0%, ${app.color}CC 100%)` }">
              <el-icon size="16"><component :is="app.icon" /></el-icon>
            </div>
            <div class="app-info">
              <h4>{{ app.name }}</h4>
              <p>{{ app.description }}</p>
            </div>
            <div class="app-arrow">
              <el-icon><ArrowRight /></el-icon>
            </div>
          </div>
        </div>
      </div>

      <!-- 个人信息卡片 -->
      <div class="section-card personal-section">
        <div class="section-header">
          <div class="header-left">
            <div class="section-icon personal-icon">
              <el-icon><User /></el-icon>
            </div>
            <div>
              <h3>个人信息</h3>
              <p>账户详情</p>
            </div>
          </div>
        </div>
        <div class="personal-info">
          <!-- 主要信息区域 -->
          <div class="main-info">
            <div class="user-basic">
              <div class="user-name">
                <el-icon><User /></el-icon>
                <span>{{ userInfo.name }}</span>
              </div>
              <div class="user-account">
                <el-icon><Key /></el-icon>
                <span>{{ userInfo.userAccount }}</span>
              </div>
            </div>
            <div class="user-roles">
              <div class="roles-label">
                <el-icon><Star /></el-icon>
                <span>角色</span>
              </div>
              <div class="roles-list">
                <span v-for="role in userInfo.roles" :key="role" class="role-tag">{{ role }}</span>
              </div>
            </div>
          </div>
          
          <!-- 详细信息区域 -->
          <div class="detail-info">
            <div class="info-row">
              <div class="info-item">
                <el-icon><UserFilled /></el-icon>
                <span class="info-label">用户类型</span>
                <span class="info-value">{{ userInfo.userType }}</span>
              </div>
              <div class="info-item">
                <el-icon><Position /></el-icon>
                <span class="info-label">职务</span>
                <span class="info-value">{{ userInfo.position }}</span>
              </div>
            </div>
            <div class="info-row">
              <div class="info-item">
                <el-icon><OfficeBuilding /></el-icon>
                <span class="info-label">所属机构</span>
                <span class="info-value">{{ userInfo.organization }}</span>
              </div>
              <div class="info-item">
                <el-icon><Briefcase /></el-icon>
                <span class="info-label">所属部门</span>
                <span class="info-value">{{ userInfo.department }}</span>
              </div>
            </div>
            <div class="info-row">
              <div class="info-item">
                <el-icon><Phone /></el-icon>
                <span class="info-label">手机号</span>
                <span class="info-value">{{ userInfo.phone }}</span>
              </div>
              <div class="info-item">
                <el-icon><Message /></el-icon>
                <span class="info-label">邮箱</span>
                <span class="info-value">{{ userInfo.email }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 消息通知 -->
      <div class="section-card message-section">
        <div class="section-header">
          <div class="header-left">
            <div class="section-icon message-icon">
              <el-icon><Bell /></el-icon>
            </div>
            <div>
              <h3>消息通知</h3>
              <p>系统公告 & 应用信息</p>
            </div>
          </div>
          <div class="message-badge" v-if="unreadMessageCount > 0">
            {{ unreadMessageCount }}
          </div>
        </div>
        <div class="message-list">
          <div 
            v-for="message in messages" 
            :key="message.id"
            class="message-item"
            :class="{ 'unread': !message.isRead }"
            @click="handleMessageClick(message)"
          >
            <div class="message-icon-wrapper" :style="{ backgroundColor: getPriorityColor(message.priority) }">
              <el-icon><component :is="getMessageIcon(message.type)" /></el-icon>
            </div>
            <div class="message-content">
              <div class="message-header">
                <span class="message-title">{{ message.title }}</span>
                <span class="message-type">{{ message.type === 'system' ? '系统' : '应用' }}</span>
              </div>
              <p class="message-desc">{{ message.content }}</p>
              <div class="message-footer">
                <span class="message-time">{{ message.time }}</span>
                <div class="message-priority" :style="{ backgroundColor: getPriorityColor(message.priority) }">
                  {{ message.priority === 'high' ? '高' : message.priority === 'medium' ? '中' : '低' }}
                </div>
              </div>
            </div>
            <div class="message-arrow" v-if="!message.isRead">
              <el-icon><CircleCheck /></el-icon>
            </div>
          </div>
        </div>
      </div>

    </div>

    <!-- 业务系统弹窗 -->
    <el-dialog
      v-model="businessDialogVisible"
      title="业务系统"
      width="80%"
      :before-close="handleCloseBusinessDialog"
    >
      <div class="business-dialog-content">
        <div class="dialog-header">
          <h3>所有业务系统</h3>
          <p>点击系统卡片进入对应业务模块</p>
        </div>
        <div class="business-dialog-grid">
          <div 
            v-for="system in allBusinessSystems" 
            :key="system.id"
            class="business-dialog-card"
            :class="{ 'no-permission': !hasPermission(system.name) }"
          >
            <div class="system-icon" :style="{ background: `linear-gradient(135deg, ${system.color} 0%, ${system.color}CC 100%)` }">
              <el-icon size="24"><component :is="system.icon" /></el-icon>
            </div>
            <div class="system-info">
              <h4>{{ system.name }}</h4>
              <p>{{ system.description }}</p>
            </div>
            <div class="system-actions">
              <div class="action-buttons">
                <el-button 
                  v-if="hasPermission(system.name)"
                  type="primary" 
                  size="small"
                  @click="handleSystemClick(system)"
                >
                  进入系统
                </el-button>
                <el-button 
                  v-else
                  type="warning" 
                  size="small"
                  @click="applyPermission(system.name)"
                >
                  申请权限
                </el-button>
                <el-button 
                  v-if="hasPermission(system.name)"
                  :type="isFavorite(system.name) ? 'success' : 'default'"
                  size="small"
                  :icon="isFavorite(system.name) ? 'StarFilled' : 'Star'"
                  @click.stop="toggleFavorite(system.name)"
                >
                  {{ isFavorite(system.name) ? '已收藏' : '收藏' }}
                </el-button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<style scoped>
.portal-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #1e3c72 0%, #2a5298 50%, #1e3c72 100%);
  position: relative;
  overflow-x: hidden;
}

.portal-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: 
    radial-gradient(circle at 20% 80%, rgba(120, 119, 198, 0.3) 0%, transparent 50%),
    radial-gradient(circle at 80% 20%, rgba(255, 255, 255, 0.1) 0%, transparent 50%),
    radial-gradient(circle at 40% 40%, rgba(120, 119, 198, 0.2) 0%, transparent 50%);
  pointer-events: none;
}

/* 顶部导航栏 */
.top-navbar {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  padding: 0 24px;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.1);
}

.navbar-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 64px;
  max-width: 1400px;
  margin: 0 auto;
}

.logo-section {
  display: flex;
  align-items: center;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  color: #1e3c72;
  font-weight: 700;
  font-size: 20px;
}

.logo-text {
  background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.navbar-right {
  display: flex;
  align-items: center;
  gap: 24px;
}

.weather-info {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #1e3c72;
  font-weight: 500;
  font-size: 14px;
}

.user-profile {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 16px;
  background: rgba(30, 60, 114, 0.1);
  border-radius: 24px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.user-profile:hover {
  background: rgba(30, 60, 114, 0.15);
  transform: translateY(-1px);
}

.user-details {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.user-name {
  font-weight: 600;
  color: #1e3c72;
  font-size: 14px;
}

.user-role {
  font-size: 12px;
  color: #666;
}

/* 主要内容区域 */
.main-content {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 16px;
  padding: 16px 24px;
  max-width: 1600px;
  margin: 0 auto;
  position: relative;
  z-index: 1;
  height: calc(100vh - 80px);
  overflow: hidden;
}


/* 卡片样式 */
.section-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  padding: 16px;
  margin-bottom: 0;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: all 0.3s ease;
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* 业务中心 */
.business-section {
  grid-column: 1;
  grid-row: 1;
}

/* 应用中心 */
.app-section {
  grid-column: 1;
  grid-row: 2;
}

/* 个人信息 */
.personal-section {
  grid-column: 2;
  grid-row: 1;
}

/* 消息通知 */
.message-section {
  grid-column: 2;
  grid-row: 2;
}

.message-icon {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.section-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.15);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #f0f2f5;
  flex-shrink: 0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.section-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 20px;
}

.app-icon {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.business-icon {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.personal-icon {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.quick-icon {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.recent-icon {
  background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
}

.section-header h3 {
  margin: 0 0 4px 0;
  color: #1e3c72;
  font-size: 20px;
  font-weight: 700;
}

.section-header p {
  margin: 0;
  color: #666;
  font-size: 14px;
  font-weight: 500;
}

.more-btn {
  background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
  border: none;
  color: white;
  font-weight: 600;
  padding: 8px 16px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 6px;
  transition: all 0.3s ease;
}

.more-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(30, 60, 114, 0.3);
}

/* 应用网格 */
.app-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 6px;
  flex: 1;
  overflow-y: auto;
}

.app-card {
  display: flex;
  align-items: center;
  padding: 8px;
  border-radius: 8px;
  background: #fafbfc;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid transparent;
  position: relative;
  overflow: hidden;
  min-height: 50px;
}

.app-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(30, 60, 114, 0.05) 0%, rgba(42, 82, 152, 0.05) 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.app-card:hover::before {
  opacity: 1;
}

.app-card:hover {
  background: #f0f9ff;
  border-color: #1e3c72;
  transform: translateY(-4px);
  box-shadow: 0 8px 25px rgba(30, 60, 114, 0.15);
}

.app-icon-wrapper {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  margin-right: 8px;
  position: relative;
  z-index: 1;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  flex-shrink: 0;
}

.app-info {
  flex: 1;
  position: relative;
  z-index: 1;
  min-width: 0;
}

.app-info h4 {
  margin: 0 0 1px 0;
  color: #1e3c72;
  font-size: 10px;
  font-weight: 700;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.app-info p {
  margin: 0;
  color: #666;
  font-size: 8px;
  line-height: 1.2;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.app-arrow {
  color: #1e3c72;
  opacity: 0;
  transition: all 0.3s ease;
  position: relative;
  z-index: 1;
  margin-left: 6px;
  font-size: 12px;
}

.app-card:hover .app-arrow {
  opacity: 1;
  transform: translateX(4px);
}

/* 业务网格 */
.business-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 6px;
  flex: 1;
  overflow-y: auto;
}

.business-card {
  display: flex;
  align-items: center;
  padding: 8px;
  border-radius: 8px;
  background: #fafbfc;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid transparent;
  position: relative;
  overflow: hidden;
  min-height: 50px;
}

.business-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(30, 60, 114, 0.05) 0%, rgba(42, 82, 152, 0.05) 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.business-card:hover::before {
  opacity: 1;
}

.business-card:hover {
  background: #f0f9ff;
  border-color: #1e3c72;
  transform: translateY(-4px);
  box-shadow: 0 8px 25px rgba(30, 60, 114, 0.15);
}

.business-icon-wrapper {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  margin-right: 8px;
  position: relative;
  z-index: 1;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.12);
  flex-shrink: 0;
}

.business-info {
  flex: 1;
  position: relative;
  z-index: 1;
  min-width: 0;
}

.business-info h4 {
  margin: 0 0 1px 0;
  color: #1e3c72;
  font-size: 10px;
  font-weight: 700;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.business-info p {
  margin: 0;
  color: #666;
  font-size: 8px;
  line-height: 1.2;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.business-arrow {
  color: #1e3c72;
  opacity: 0;
  transition: all 0.3s ease;
  position: relative;
  z-index: 1;
  margin-left: 6px;
  font-size: 12px;
}

.business-card:hover .business-arrow {
  opacity: 1;
  transform: translateX(4px);
}

/* 个人信息 */
.personal-info {
  display: flex;
  flex-direction: column;
  gap: 12px;
  flex: 1;
  overflow-y: auto;
}

/* 主要信息区域 */
.main-info {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.user-basic {
  display: flex;
  gap: 16px;
  align-items: center;
}

.user-name, .user-account {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  border-radius: 6px;
  background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
  color: white;
  font-weight: 600;
  font-size: 12px;
}

.user-name .el-icon, .user-account .el-icon {
  font-size: 14px;
}

.user-roles {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.roles-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 10px;
  color: #666;
  font-weight: 500;
}

.roles-label .el-icon {
  font-size: 12px;
  color: #1e3c72;
}

.roles-list {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.role-tag {
  padding: 2px 8px;
  border-radius: 12px;
  background: #e8f4fd;
  color: #1e3c72;
  font-size: 8px;
  font-weight: 500;
  border: 1px solid #b3d8ff;
}

/* 详细信息区域 */
.detail-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.info-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 8px;
  border-radius: 4px;
  background: #f8f9fa;
  border: 1px solid #f0f2f5;
  min-height: 24px;
}

.info-item .el-icon {
  width: 14px;
  height: 14px;
  color: #1e3c72;
  flex-shrink: 0;
}

.info-label {
  font-size: 8px;
  color: #999;
  font-weight: 500;
  min-width: 50px;
}

.info-value {
  font-size: 9px;
  color: #1e3c72;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
}

/* 消息通知 */
.message-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
  overflow-y: auto;
}

.message-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px;
  border-radius: 4px;
  background: #f8f9fa;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid transparent;
  position: relative;
  min-height: 32px;
}

.message-item.unread {
  background: #f0f9ff;
  border-color: #409EFF;
  box-shadow: 0 1px 3px rgba(64, 158, 255, 0.1);
}

.message-item:hover {
  background: #e9ecef;
  border-color: #1e3c72;
  transform: translateY(-1px);
  box-shadow: 0 2px 6px rgba(30, 60, 114, 0.15);
}

.message-item.unread:hover {
  background: #e6f7ff;
  border-color: #409EFF;
}

.message-icon-wrapper {
  width: 16px;
  height: 16px;
  border-radius: 3px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 8px;
  flex-shrink: 0;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.message-content {
  display: flex;
  flex-direction: column;
  gap: 1px;
  flex: 1;
  min-width: 0;
}

.message-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 4px;
}

.message-title {
  font-size: 9px;
  color: #1e3c72;
  font-weight: 700;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
  line-height: 1.1;
}

.message-type {
  font-size: 7px;
  color: #666;
  background: #e9ecef;
  padding: 1px 3px;
  border-radius: 2px;
  font-weight: 500;
  flex-shrink: 0;
}

.message-desc {
  font-size: 8px;
  color: #666;
  line-height: 1.2;
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.message-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 4px;
  margin-top: 1px;
}

.message-time {
  font-size: 7px;
  color: #999;
  font-weight: 500;
}

.message-priority {
  font-size: 6px;
  color: white;
  padding: 1px 3px;
  border-radius: 2px;
  font-weight: 600;
  flex-shrink: 0;
}

.message-arrow {
  color: #67C23A;
  font-size: 8px;
  opacity: 0;
  transition: all 0.3s ease;
  flex-shrink: 0;
}

.message-item:hover .message-arrow {
  opacity: 1;
  transform: scale(1.1);
}

.message-badge {
  background: linear-gradient(135deg, #F56C6C 0%, #f56c6c 100%);
  color: white;
  padding: 1px 4px;
  border-radius: 6px;
  font-size: 8px;
  font-weight: 700;
  min-width: 12px;
  height: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 1px 2px rgba(245, 108, 108, 0.3);
}


/* 响应式设计 */
@media (max-width: 1200px) {
  .main-content {
    grid-template-columns: 1fr;
    height: auto;
    overflow: visible;
  }
  
  .section-card {
    height: auto;
  }
}

@media (max-width: 768px) {
  .main-content {
    padding: 12px 16px;
    gap: 12px;
  }
  
  .top-navbar {
    padding: 0 16px;
  }
  
  .navbar-right {
    gap: 12px;
  }
  
  .user-profile {
    padding: 4px 8px;
  }
  
  .app-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .business-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .quick-actions {
    grid-template-columns: 1fr;
  }
  
  .section-card {
    padding: 12px;
  }
}

/* 业务系统弹窗样式 */
.business-dialog-content {
  padding: 20px 0;
}

.dialog-header {
  text-align: center;
  margin-bottom: 30px;
}

.dialog-header h3 {
  font-size: 24px;
  color: #1e3c72;
  margin: 0 0 8px 0;
  font-weight: 700;
}

.dialog-header p {
  font-size: 14px;
  color: #666;
  margin: 0;
}

.business-dialog-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  max-height: 60vh;
  overflow-y: auto;
}

.business-dialog-card {
  display: flex;
  align-items: center;
  padding: 16px;
  border-radius: 12px;
  background: #fafbfc;
  border: 1px solid #e8f4fd;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
  min-height: 80px;
}

.business-dialog-card.no-permission {
  background: #f5f5f5;
  border-color: #d9d9d9;
  opacity: 0.8;
}

.business-dialog-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(30, 60, 114, 0.05) 0%, rgba(42, 82, 152, 0.05) 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.business-dialog-card:hover::before {
  opacity: 1;
}

.business-dialog-card:hover {
  background: #f0f9ff;
  border-color: #1e3c72;
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(30, 60, 114, 0.15);
}

.system-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  margin-right: 16px;
  position: relative;
  z-index: 1;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  flex-shrink: 0;
}

.system-info {
  flex: 1;
  position: relative;
  z-index: 1;
  min-width: 0;
  margin-right: 12px;
}

.system-info h4 {
  margin: 0 0 4px 0;
  color: #1e3c72;
  font-size: 16px;
  font-weight: 700;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.system-info p {
  margin: 0;
  color: #666;
  font-size: 12px;
  line-height: 1.4;
  word-wrap: break-word;
}


.system-actions {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 8px;
  position: relative;
  z-index: 1;
  flex-shrink: 0;
}

.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 6px;
  width: 100px;
}

</style>
