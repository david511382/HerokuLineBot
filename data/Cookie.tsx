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
    const RAW_KEY_COOKIES_STR = RAW_COOKIES_STR.substring(KEY_START_INDEX)
    const EQUAL_SYMBOL_INDEX = RAW_KEY_COOKIES_STR.indexOf('=')
    const VALUE_COOKIES_STR = RAW_KEY_COOKIES_STR.substring(EQUAL_SYMBOL_INDEX + 1)
    const VALUE_END_INDEX = VALUE_COOKIES_STR.indexOf(';')
    if (VALUE_END_INDEX === -1) {
        return VALUE_COOKIES_STR.substring(0)
    }
    return VALUE_COOKIES_STR.substring(0, VALUE_END_INDEX - 1)
}