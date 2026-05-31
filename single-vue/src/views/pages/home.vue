<template>
  <div style="min-height: 100vh;">
    <h1>首页</h1>
    <div ref="chartRef" style="width: 600px; height: 400px;"></div>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, onBeforeUnmount, ref } from 'vue'
import * as echarts from 'echarts'

const chartRef = ref<HTMLDivElement | null>(null)
let chartInstance: echarts.ECharts | null = null

onMounted(() => {
  if (!chartRef.value) return

  chartInstance = echarts.init(chartRef.value)

  chartInstance.setOption({
    title: { text: 'ECharts 入门示例' },
    tooltip: {},
    xAxis: {
      data: ['衬衫', '羊毛衫', '雪纺衫', '裤子', '高跟鞋', '袜子']
    },
    yAxis: {},
    series: [
      {
        name: '销量',
        type: 'bar',
        data: [5, 20, 36, 10, 10, 20]
      }
    ]
  })

  window.addEventListener('resize', resizeChart)
})

function resizeChart() {
  chartInstance?.resize()
}

onBeforeUnmount(() => {
  window.removeEventListener('resize', resizeChart)
  chartInstance?.dispose()
})

</script>
