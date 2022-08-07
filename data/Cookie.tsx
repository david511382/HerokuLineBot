export function GetCookies(): Map<any, any> {
    let cookies = document.cookie.split('; ')
    let cookieMap = new Map()
    cookies.forEach(
        (value: string, index: number, array: string[]) => {
            let keyValues = value.split('=')
            if (keyValues.length == 0) {
                return
            }
            let key = keyValues[0]
            if (keyValues.length > 1) {
                value = keyValues[1]
            } else {
                value = ""
            }
            cookieMap.set(key, value)
        })

    return cookieMap
}

export function GetCookieOnKey(key: string): string {
    const RAW_COOKIES_STR = document.cookie
    const KEY_START_INDEX = RAW_COOKIES_STR.indexOf(key)
    const RAW_KEY_COOKIES_STR = RAW_COOKIES_STR.substr(KEY_START_INDEX)
    const EQUAL_SYMBOL_INDEX = RAW_KEY_COOKIES_STR.indexOf('=')
    const RAW_EQUAL_SYMBOL_COOKIES_STR = RAW_KEY_COOKIES_STR.substr(EQUAL_SYMBOL_INDEX)
    const VALUE_END_INDEX = RAW_EQUAL_SYMBOL_COOKIES_STR.indexOf(';')
    return RAW_EQUAL_SYMBOL_COOKIES_STR.substr(0, VALUE_END_INDEX - 1)
}