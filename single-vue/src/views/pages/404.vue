<template>
  <div>
    <el-result
      icon="warning"
      title="404"
    >
      <template #extra>
        <el-button @click="getHome()"> 回到首页 </el-button>
      </template>

    </el-result>
  </div>
</template>

<script lang="ts" setup>
import { useRouter } from 'vue-router';
let router  = useRouter()
const getHome = ()=>{
  router.push('/')
}

</script>

<style scoped>

</style>
