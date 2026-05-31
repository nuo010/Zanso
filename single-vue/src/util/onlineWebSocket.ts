/**
 * 用户实时在线统计 WebSocket
 * 连接 wss 服务，接收并解析在线人数，支持断线重连
 */

const WS_BASE = 'wss://cfsingle.fun/server/singleWebSocket/onlineUser';

export type OnlineCountCallback = (count: number) => void;
export type StatusCallback = (status: 'connecting' | 'connected' | 'closed' | 'error') => void;

/** 解析服务端消息中的在线人数，兼容多种格式 */
function parseOnlineCount(data: string | unknown): number | null {
  if (typeof data === 'number' && Number.isFinite(data)) return data;
  if (typeof data !== 'string') return null;
  const num = Number(data);
  if (Number.isFinite(num)) return num;
  try {
    const json = JSON.parse(data) as Record<string, unknown>;
    const count =
      typeof json.onlineCount === 'number' ? json.onlineCount
      : typeof json.count === 'number' ? json.count
      : typeof json.online === 'number' ? json.online
      : null;
    return count;
  } catch {
    return null;
  }
}

class OnlineWebSocketManager {
  private ws: WebSocket | null = null;
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  private readonly reconnectDelay = 3000;
  private readonly maxReconnectDelay = 30000;
  private currentDelay = 3000;
  private countCallbacks: Set<OnlineCountCallback> = new Set();
  private statusCallbacks: Set<StatusCallback> = new Set();
  private _lastCount = 0;
  /** 连接时使用的 userId，重连时会复用 */
  private _userId: string = '';
  /** 心跳定时器：每 5 秒向服务端发送 ping */
  private heartbeatTimer: ReturnType<typeof setInterval> | null = null;
  private readonly heartbeatInterval = 5000;

  get lastCount(): number {
    return this._lastCount;
  }

  /** 订阅在线人数变化 */
  onCount(cb: OnlineCountCallback): () => void {
    this.countCallbacks.add(cb);
    cb(this._lastCount);
    return () => this.countCallbacks.delete(cb);
  }

  /** 订阅连接状态变化 */
  onStatus(cb: StatusCallback): () => void {
    this.statusCallbacks.add(cb);
    return () => this.statusCallbacks.delete(cb);
  }

  private emitCount(count: number) {
    this._lastCount = count;
    this.countCallbacks.forEach((cb) => cb(count));
  }

  private emitStatus(status: 'connecting' | 'connected' | 'closed' | 'error') {
    this.statusCallbacks.forEach((cb) => cb(status));
  }

  /**
   * 建立连接，仅在有 userId 时连接（只统计登录用户）
   * @param userId 用户 ID，会拼接到 wss://.../singleWebSocket/ 后面；为空则不连接
   */
  connect(userId?: string) {
    if (userId !== undefined && userId !== '') this._userId = String(userId);
    if (!this._userId) return;
    if (this.ws?.readyState === WebSocket.OPEN) return;
    const url = `${WS_BASE}/${encodeURIComponent(this._userId)}`;
    this.emitStatus('connecting');
    try {
      this.ws = new WebSocket(url);
    } catch (e) {
      this.emitStatus('error');
      this.scheduleReconnect();
      return;
    }

    this.ws.onopen = () => {
      this.currentDelay = this.reconnectDelay;
      this.emitStatus('connected');
      this.startHeartbeat();
    };

    this.ws.onmessage = (event) => {
      const data = event.data;
      if (typeof data === 'string' && data.trim().toLowerCase() === 'ping') {
        this.ws?.send('pong');
        return;
      }
      const count = parseOnlineCount(data);
      if (count !== null) this.emitCount(count);
    };

    this.ws.onerror = () => {
      this.emitStatus('error');
    };

    this.ws.onclose = () => {
      this.stopHeartbeat();
      this.ws = null;
      this.emitStatus('closed');
      this.scheduleReconnect();
    };
  }

  private startHeartbeat() {
    this.stopHeartbeat();
    this.heartbeatTimer = setInterval(() => {
      if (this.ws?.readyState === WebSocket.OPEN) {
        this.ws.send('ping');
      }
    }, this.heartbeatInterval);
  }

  private stopHeartbeat() {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer);
      this.heartbeatTimer = null;
    }
  }

  private scheduleReconnect() {
    if (this.reconnectTimer) return;
    this.reconnectTimer = setTimeout(() => {
      this.reconnectTimer = null;
      this.currentDelay = Math.min(
        this.currentDelay * 1.5,
        this.maxReconnectDelay
      );
      this.connect();
    }, this.currentDelay);
  }

  disconnect() {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
    this.stopHeartbeat();
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    this.emitStatus('closed');
  }
}

export const onlineWebSocket = new OnlineWebSocketManager();
