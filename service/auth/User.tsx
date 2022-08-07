import { Response } from '../../models/resp/base'
import { UserInfo } from '../../models/resp/user-info'
import { GetUserInfo as GetUserInfoApi } from '../../data/api/UserInfo'
import { GetAuthRequestInit } from '../../data/auth/Auth'

export async function GetUserInfo(): Promise<Response<UserInfo>> {
  const init = GetAuthRequestInit()
  return await GetUserInfoApi(init)
}