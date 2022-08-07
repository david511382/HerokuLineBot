import { GetCookieOnKey } from '../Cookie'

const TOKEN_COOKIE_NAME = "token";

export function GetToken(): string {
  return GetCookieOnKey(TOKEN_COOKIE_NAME)
}

export function SetIDToken(idToken: string) {
  // Set a cookie
  document.cookie = TOKEN_COOKIE_NAME + '=' + idToken + ";path=/";
}