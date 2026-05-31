// utils/perf-monitor-spa.ts
import type { Router, RouteLocationNormalized, NavigationGuardNext } from 'vue-router'
import type { AxiosInstance } from 'axios'

interface PerfRecord {
    type: 'fp' | 'fcp' | 'first-screen' | 'route' | 'api' | 'component' | 'lcp' | 'cls' | 'inp' | 'runtime'
    data: Record<string, any>
    ts: number
    url: string
    env?: Record<string, any>
    sysName?: string
}

interface PerfMonitorOptions {
    reportUrl?: string
    reportInterval?: number
    maxQueueSize?: number
    debug?: boolean
    runtimeInterval?: number // 新增：运行时健康监控周期
    sysName?: string
}

class PerfMonitor {
    private router?: Router
    private reportUrl: string
    private reportInterval: number
    private queue: PerfRecord[] = []
    private firstPaint: number | null = null
    private firstContentfulPaint: number | null = null
    private maxQueueSize: number
    private debug: boolean
    private runtimeInterval: number
    private sysName?: string

    constructor(router?: Router, options: PerfMonitorOptions = {}) {
        this.router = router
        this.reportUrl = options.reportUrl || '/api/monitor/perf'
        this.reportInterval = options.reportInterval || 10000
        this.maxQueueSize = options.maxQueueSize || 500
        this.debug = options.debug ?? false
        this.runtimeInterval = options.runtimeInterval || 60000 // 默认1分钟采集一次
        this.sysName = options.sysName || 'default-system'

        this.observePerformance()
        this.observeWebVitals()
        if (router) this.observeRouter()
        this.observeRuntimeHealth()
        this.startReporter()
    }

    private getEnvInfo() {
        return {
            ua: navigator.userAgent,
            lang: navigator.language,
            network: (navigator as any).connection?.effectiveType || 'unknown',
            memory: (navigator as any).deviceMemory || 'unknown',
            cpu: (navigator as any).hardwareConcurrency || 'unknown'
        }
    }

    private log(type: PerfRecord['type'], data: Record<string, any>) {
        const record: PerfRecord = {
            type,
            data,
            ts: Date.now(),
            url: window.location.href,
            env: this.getEnvInfo(),
            sysName: this.sysName
        }
        if (this.queue.length >= this.maxQueueSize) this.queue.shift()
        this.queue.push(record)
        if (this.debug) console.log('[Perf]', record)
    }

    private observePerformance() {
        if ('PerformanceObserver' in window) {
            const paintObserver = new PerformanceObserver((list: PerformanceObserverEntryList) => {
                for (const entry of list.getEntries()) {
                    const paintEntry = entry as PerformanceEntry
                    if (paintEntry.name === 'first-paint') {
                        this.firstPaint = paintEntry.startTime
                        this.log('fp', { value: this.firstPaint })
                    }
                    if (paintEntry.name === 'first-contentful-paint') {
                        this.firstContentfulPaint = paintEntry.startTime
                        this.log('fcp', { value: this.firstContentfulPaint })
                    }
                }
            })
            paintObserver.observe({ type: 'paint', buffered: true })
        }

        window.addEventListener('load', () => {
            let firstScreenTime = 0
            if ('getEntriesByType' in performance) {
                const navEntries = performance.getEntriesByType('navigation') as PerformanceNavigationTiming[]
                if (navEntries.length > 0) {
                    firstScreenTime = navEntries[0].loadEventEnd
                } else {
                    firstScreenTime = performance.now()
                }
            }
            this.log('first-screen', { value: firstScreenTime })
        })
    }

    private observeWebVitals() {
        if (!('PerformanceObserver' in window)) return

        // LCP
        const lcpObserver = new PerformanceObserver((entryList) => {
            const entries = entryList.getEntries()
            const lastEntry = entries[entries.length - 1] as any
            this.log('lcp', { value: lastEntry.startTime })
        })
        lcpObserver.observe({ type: 'largest-contentful-paint', buffered: true })

        // CLS
        const clsObserver = new PerformanceObserver((entryList) => {
            for (const entry of entryList.getEntries() as any[]) {
                if (!entry.hadRecentInput) {
                    this.log('cls', { value: entry.value })
                }
            }
        })
        clsObserver.observe({ type: 'layout-shift', buffered: true })

        // INP
        try {
            const inpObserver = new PerformanceObserver((entryList) => {
                for (const entry of entryList.getEntries() as any[]) {
                    this.log('inp', { name: entry.name, duration: entry.duration })
                }
            })
            inpObserver.observe({ type: 'event', buffered: true })
        } catch (e) {
            if (this.debug) console.warn('[Perf] INP not supported in this browser')
        }
    }

    private observeRouter() {
        let routeStart = 0
        this.router?.beforeEach(
            (to: RouteLocationNormalized, from: RouteLocationNormalized, next: NavigationGuardNext) => {
                routeStart = performance.now()
                next()
            }
        )

        this.router?.afterEach(async (to: RouteLocationNormalized, from: RouteLocationNormalized) => {
            const routeDuration = performance.now() - routeStart
            this.log('route', {
                from: from.fullPath,
                to: to.fullPath,
                duration: routeDuration
            })

            // 懒加载组件首屏渲染耗时
            if (to.matched && to.matched.length > 0) {
                for (const record of to.matched) {
                    const component = record.components?.default
                    if (component && typeof component === 'function') {
                        const compStart = performance.now()
                        await (component as () => Promise<any>)()
                        const compDuration = performance.now() - compStart
                        this.log('component', {
                            route: to.fullPath,
                            component: record.name || 'anonymous',
                            duration: compDuration
                        })
                    }
                }
            }
        })
    }

    private observeRuntimeHealth() {
        // 定时采集
        setInterval(() => {
            // 内存
            const memory = (performance as any).memory
            if (memory) {
                this.log('runtime', {
                    usedJSHeapSize: memory.usedJSHeapSize,
                    totalJSHeapSize: memory.totalJSHeapSize,
                    jsHeapSizeLimit: memory.jsHeapSizeLimit
                })
            }

            // FPS
            let frameCount = 0
            const start = performance.now()
            const loop = () => {
                frameCount++
                const now = performance.now()
                if (now - start >= 1000) {
                    this.log('runtime', { fps: frameCount })
                } else {
                    requestAnimationFrame(loop)
                }
            }
            requestAnimationFrame(loop)

            // 白屏检测
            const app = document.getElementById('app')
            if (app && app.clientHeight === 0) {
                this.log('runtime', { whiteScreen: true })
            }
        }, this.runtimeInterval)

        // 长任务监控
        if ('PerformanceObserver' in window) {
            try {
                const longTaskObserver = new PerformanceObserver((entryList) => {
                    for (const entry of entryList.getEntries()) {
                        this.log('runtime', {
                            longTask: entry.duration,
                            startTime: entry.startTime
                        })
                    }
                })
                longTaskObserver.observe({ type: 'longtask', buffered: true })
            } catch (e) {
                if (this.debug) console.warn('[Perf] LongTask not supported')
            }
        }
    }

    static wrapFetch(monitor: PerfMonitor) {
        const originalFetch = window.fetch
        window.fetch = async (...args: Parameters<typeof fetch>): Promise<Response> => {
            const url = args[0] as string
            const start = performance.now()
            try {
                const res = await originalFetch(...args)
                const duration = performance.now() - start
                monitor.log('api', { url, duration, status: res.status })
                return res
            } catch (err: unknown) {
                const duration = performance.now() - start
                const errorMessage = (err as Error).message || String(err)
                monitor.log('api', { url, duration, error: errorMessage })
                throw err
            }
        }
    }

    static wrapAxios(instance: AxiosInstance, monitor: PerfMonitor){
        instance.interceptors.request.use((config) => {
            (config as any).metadata = { startTime: performance.now() }
            return config
        })

        instance.interceptors.response.use(
            (response) => {
                const startTime = (response.config as any).metadata?.startTime || 0
                const duration = performance.now() - startTime
                monitor.log('api', {
                    url: response.config.url,
                    method: response.config.method,
                    duration,
                    status: response.status
                })
                return response
            },
            (error: any) => {
                if (error.config && (error.config as any).metadata) {
                    const duration = performance.now() - ((error.config as any).metadata.startTime || 0)
                    monitor.log('api', {
                        url: error.config.url,
                        method: error.config.method,
                        duration,
                        error: error.message
                    })
                }
                return Promise.reject(error)
            }
        )
    }

    private startReporter() {
        const sendPayload = (payload: PerfRecord[]) => {
            const sendBeacon = () => {
                if (navigator.sendBeacon) {
                    const blob = new Blob([JSON.stringify(payload)], { type: 'application/json' })
                    return navigator.sendBeacon(this.reportUrl, blob)
                }
                return false
            }

            const fallbackFetch = () => {
                fetch(this.reportUrl, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(payload),
                    keepalive: true
                }).catch(err => this.debug && console.warn('[Perf] fallback fetch 上报失败', err))
            }

            try {
                const success = sendBeacon()
                if (!success) fallbackFetch()
            } catch (err) {
                this.debug && console.warn('[Perf] 上报异常，已捕获', err)
                this.queue.unshift(...payload)
            }
        }

        setInterval(() => {
            if (this.queue.length === 0) return
            const payload = this.queue.splice(0, this.queue.length)
            sendPayload(payload)
        }, this.reportInterval)
    }
}

export default PerfMonitor
