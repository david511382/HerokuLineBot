import { UserInfo } from '../../models/resp/user-info'

export async function GetUserInfo() :Promise<UserInfo> {
  return await fetch('/api/user-info')
    .then((response) => {
      if (!response.ok) {
        throw new Error(response.statusText)
      }
      return response.json()
    })
}
