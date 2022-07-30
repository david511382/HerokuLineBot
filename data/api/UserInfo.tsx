import { Response } from '../../models/resp/base'
import { UserInfo as Data } from '../../models/resp/user-info'
import { GetBackendUrl } from '../env/Http'

export async function GetUserInfo(init?: RequestInit | undefined): Promise<Response<Data>> {
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