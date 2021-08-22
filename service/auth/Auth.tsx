import { GetToken } from './Token'

export function GetAuthRequestInit(init?: RequestInit | undefined) :RequestInit  {
  let token = GetToken()
  if (!init){
    init = {}
  }
  
  if (init.headers){
    let h = init.headers as Record<string, string>
    h["Authorization"] = token
  }else{
    init.headers = [["Authorization",token]]
  }

  return init
}