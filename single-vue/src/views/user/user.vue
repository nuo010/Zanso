<template>
  <div class="user-container">
    <div class="user-header">
      <h1 class="user-title">
        <el-icon class="title-icon"><User /></el-icon>
        个人中心
      </h1>
    </div>

    <div v-if="store.user" class="user-content">
      <el-card class="user-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <el-icon class="header-icon"><Avatar /></el-icon>
            <span>基本信息</span>
          </div>
        </template>
        
        <div class="user-info-grid">
          <div class="info-item">
            <div class="info-label">
              <el-icon><User /></el-icon>
              <span>用户名</span>
            </div>
            <div class="info-value">{{ store.user.userName || '未设置' }}</div>
          </div>

          <div class="info-item">
            <div class="info-label">
              <el-icon><UserFilled /></el-icon>
              <span>昵称</span>
            </div>
            <div class="info-value">{{ store.user.nickName || '未设置' }}</div>
          </div>
          <div class="info-item">
            <div class="info-label">
              <el-icon><Message /></el-icon>
              <span>邮箱</span>
            </div>
            <div class="info-value">{{ store.user.email || '未设置' }}</div>
          </div>
          <div class="info-item">
            <div class="info-label">
              <el-icon><Location /></el-icon>
              <span>登录IP地址</span>
            </div>
            <div class="info-value">{{ store.user.ipaddr || '未获取' }}</div>
          </div>
          <div class="info-item">
            <div class="info-label">
              <el-icon><Key /></el-icon>
              <span>用户ID</span>
            </div>
            <div class="info-value">{{ store.user.userId || '未设置' }}</div>
          </div>
          <div class="info-item">
            <div class="info-label">
              <el-icon><Key /></el-icon>
              <span>用户Token</span>
            </div>
            <div class="info-value">{{ getToken() }}</div>
          </div>

          

          

          

        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { userMainStore } from "@/store";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { 
  User, 
  Avatar, 
  Key, 
  UserFilled, 
  Message, 
  Location, 
  Setting, 
  Edit, 
  Download 
} from "@element-plus/icons-vue";
import {getToken} from "../../util/auth";

const store = userMainStore();
const router = useRouter();


</script>

<style scoped lang="scss">
.user-container {
  padding: 24px;
  height: 100%;
  background-attachment: fixed;
  overflow-y: auto;
}

.user-header {
  text-align: center;
  margin-bottom: 32px;
}

.user-title {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
}

.title-icon {
  font-size: 2.5rem;
}

.user-subtitle {
  font-size: 1.1rem;
  opacity: 0.9;
  margin: 0;
}

.user-content {
  max-width: 800px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.user-card {
  border-radius: 16px;
  overflow: hidden;
  backdrop-filter: blur(10px);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  font-size: 1.1rem;
  color: #2c3e50;
}

.header-icon {
  font-size: 1.2rem;
  color: #409eff;
}

.user-info-grid {
  display: grid;
  gap: 20px;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 12px;
  border-left: 4px solid #409eff;
  transition: all 0.3s ease;
}

.info-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  background: #e3f2fd;
}

.info-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #666;
  font-size: 0.9rem;
}

.info-label .el-icon {
  color: #409eff;
  font-size: 1rem;
}

.info-value {
  font-size: 1.1rem;
  color: #2c3e50;
  font-weight: 500;
  word-break: break-all;
}

.action-buttons {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
  justify-content: center;
}

.action-buttons .el-button {
  border-radius: 8px;
  padding: 12px 24px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.action-buttons .el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}

.empty-state .el-empty {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  padding: 40px;
  backdrop-filter: blur(10px);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .user-container {
    padding: 16px;
  }
  
  .user-title {
    font-size: 2rem;
  }
  
  .user-info-grid {
    grid-template-columns: 1fr;
  }
  
  .action-buttons {
    flex-direction: column;
    align-items: center;
  }
  
  .action-buttons .el-button {
    width: 100%;
    max-width: 200px;
  }
}
</style>