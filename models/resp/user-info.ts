export interface UserInfo {
  username: string
  role_id: RoleID
}

export enum RoleID  {
  ROLE_ADMIN = 1,
  ROLE_CADRE = 2,
  ROLE_MEMBER = 3,
  ROLE_GUEST = 4,
}
