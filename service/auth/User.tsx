import { Response } from '../../models/resp/base'
import { UserInfo } from '../../models/resp/user-info'
import { GetToken } from './Token'
import { GetUserInfo as GetUserInfoApi } from '../../data/api/UserInfo'

export async function GetUserInfo(): Promise<Response<UserInfo>> {
  let token = GetToken()
  return await GetUserInfoApi(
    {
      headers: [
        ["Authorization", token],
      ]
    }
  )
}