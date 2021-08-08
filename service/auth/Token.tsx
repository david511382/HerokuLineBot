const TOKEN_COOKIE_NAME = "token";

export function GetToken() : string{
  let cookies = document.cookie.split('; ')
  let cookieMap = new Map()
  cookies.forEach(
    (value: string, index: number, array: string[]) => {
      let keyValues = value.split('=')
      if (keyValues.length==0){
        return
      }
      let key = keyValues[0]
      if (keyValues.length>1){
        value = keyValues[1]
      }else{
        value = ""
      }
      cookieMap.set(key,value)
    })
  
  return cookieMap.get(TOKEN_COOKIE_NAME)
}

export function SetIDToken(idToken:string) {
  // Set a cookie
  document.cookie = TOKEN_COOKIE_NAME + '=' + idToken + ";path=/";
}