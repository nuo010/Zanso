import AxiosUtil from '@/util/axios';


export function addMenu(data:any){
  return AxiosUtil({
      url:"/gateway/singleSystem/addMenu",
      method:"post",
      data
    }
  )
}
export function getMenuListTree(data:any){
  return AxiosUtil({
      url:"/gateway/singleSystem/getMenuListTree",
      method:"post",
      data
    }
  )
}



export function getTblSysMenuByParentIsZero(){
  return AxiosUtil({
      url:"/gateway/singleSystem/getTblSysMenuByParentIsZero",
      method:"get",
    }
  )
}
export function getUserList(data:any){
  return AxiosUtil({
      url:"/gateway/singleSystem/getUserList",
      method:"post",
      data
    }
  )
}
export function getUserById(params:any){
  return AxiosUtil({
      url:"/gateway/singleSystem/getUserById",
      method:"get",
      params
    }
  )
}

export function getSysRoleList(data:any){
  return AxiosUtil({
      url:"/gateway/singleSystem/getSysRoleList",
      method:"post",
      data
    }
  )
}
export function getMenuIdListByRoleId(params:any){
  return AxiosUtil({
      url:"/gateway/singleSystem/getMenuIdListByRoleId",
      method:"get",
      params
    }
  )
}
export function getTblSysMenuByRoleId(params:any){
  return AxiosUtil({
      url:"/gateway/singleSystem/getTblSysMenuByRoleId",
      method:"get",
      params
    }
  )
}
export function getSysRoleListByUserId(params:any){
  return AxiosUtil({
      url:"/gateway/singleSystem/getSysRoleListByUserId",
      method:"get",
      params
    }
  )
}
export function addRoleJoinUser(data:any){
  return AxiosUtil({
      url:"/gateway/singleSystem/addRoleJoinUser",
      method:"post",
      data
    }
  )
}

export function addRoleConnectMenu(data:any){
  return AxiosUtil({
      url:"/gateway/singleSystem/addRoleConnectMenu",
      method:"post",
      data
    }
  )
}



