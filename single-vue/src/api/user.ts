import AxiosUtil from '@/util/axios';


export function doLogin(userName: string, passWord: string){
  return AxiosUtil({
      url:"/singleAuth/login",
      method:"post",
      data:{
          userName,
          passWord
      }
    }
  )
}

export function registerEmail(userName: string, passWord: string, email?: string, code?: string){
  return AxiosUtil({
      url:"/singleAuth/registerEmail",
      method:"post",
      data:{
          userName,
          passWord,
          email,
          code
      }
    }
  )
}

export function sendEmailCode(email: string){
  return AxiosUtil({
      url:"/singleMessagePush/sendEmailCode",
      method:"post",
      data:{
          email
      }
    }
  )
}

export function getUser(){
  return AxiosUtil({
      url:"/singleAuth/getUser",
      method:"get"
    }
  )
}
export function getMenuTree(params:any){
  return AxiosUtil({
      url:"/singleSystem/getMenuTree",
      params,
      method:"get"
    }
  )
}

export function logout(){
  return AxiosUtil({
      url:"/singleAuth/logout",
      method:"post"
    }
  )
}
export function getLginOperationLog(data:any){
  return AxiosUtil({
      url:"/singleSystem/getLginOperationLog",
      method:"post",
      data
    }
  )
}
export function getMenuList(data:any){
  return AxiosUtil({
      url:"/singleSystem/getMenuList",
      method:"post",
      data
    }
  )
}
export function changePassword(data:any){
    return AxiosUtil({
            url:"/singleAuth/changePassword",
            method:"post",
            data
        }
    )
}
