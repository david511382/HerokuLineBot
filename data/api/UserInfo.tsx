import { UserInfo } from '../../models/resp/user-info'
import { GetBackendUrl } from '../env/Http'

export async function GetUserInfo(init?: RequestInit | undefined) :Promise<UserInfo> {
  return await fetch(
      `${GetBackendUrl()}/api/user-info`,
      init,
    ).then((response) => {
      if (!response.ok) {
        throw new Error(response.statusText)
      }
      return response.json()
    })
}