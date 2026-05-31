import { useCookies } from '@vueuse/integrations/useCookies';
import { tokenKey } from "@/util/constants";




const cookie = useCookies();



export function getToken() {
  return cookie.get(tokenKey);
}

export function setToken(token: string) {
  // 设置ck过期时间为3天
  let date = new Date().getTime();
  let expiresTime = new Date(date + 60 * 1000 * 60 * 24 * 30);
  return cookie.set(tokenKey, token, {
    expires: expiresTime,
  });
}

export function removeToken() {
  return cookie.remove(tokenKey);
}
