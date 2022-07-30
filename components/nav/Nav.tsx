import React, { useEffect, useState } from 'react'
import styles from './nav.module.css'
import { GetUserInfo } from '../../service/auth/User'
import { RoleID, UserInfo as UserInfoResp } from '../../models/resp/user-info'
import { Response as Resp } from '../../models/resp/base'
import {
  Link,
  HashRouter,
  Switch,
  Route,
  Redirect,
} from "react-router-dom";

export enum PageEnum {
  LiffDebug = "Liff Debug",
  RentalCourts = "租場狀況",
}

interface Props {
  defaultPage?: PageEnum
  children?: JSX.Element | (JSX.Element | undefined)[]
}

type Page = {
  Name: PageEnum
  Path: string
  Selected: boolean
}

const PagePathMap: Map<PageEnum, string> = new Map([
  [PageEnum.LiffDebug, "/debug"],
  [PageEnum.RentalCourts, "/rental-courts"],
]);

const RolePagesMap: Map<RoleID, PageEnum[]> = new Map([
  [
    RoleID.ROLE_ADMIN,
    [
      PageEnum.LiffDebug,
      PageEnum.RentalCourts,
    ],
  ],
]);

export default function Nav(
  {
    children,
    defaultPage
  }: Props
) {
  const [userInfo, setUserInfo] = useState<UserInfoResp>()
  const [navs, setNavs] = useState<Page[]>([])
  let showPage: PageEnum
  useEffect(() => {
    if (!showPage) {
      let currentHash = window.location.hash.replaceAll("#", "")
      PagePathMap.forEach((path, page) => {
        if (showPage)
          return
        if (currentHash === path)
          showPage = page
      })
    }

    let roleID = RoleID.ROLE_GUEST
    GetUserInfo()
      .then((resp) => {
        if (!resp)
          return
        setUserInfo(resp.data)

        roleID = resp.data.role_id
      })
      .catch(() => { })
      .finally(() => {
        let pages = RolePagesMap.get(roleID)

        if (!pages) {
          pages = []
        }

        let pageViews: Page[] = []
        pages.forEach((page) => pageViews.push({
          Name: page,
          Path: GetPath(page),
          Selected: (showPage && page === showPage) ? true : false,
        }))

        setNavs(pageViews)
      })
  }, [])

  return (
    <div>
      使用者:{userInfo?.username}

      {
        navs.length > 0 &&
        <HashRouter>
          <nav className={styles.nav}>
            <ul className={styles.ul}>{
              navs.map((page, i) => {
                return (
                  <li
                    className={`${styles.li} ${(page.Selected) ? `${styles.selected}` : ""}`}
                    key={page.Name}>
                    <Link
                      to={page.Path}
                      onClick={() => {
                        navs.forEach((page) => page.Selected = false)
                        navs[i].Selected = true
                        setNavs([...navs])
                      }}>
                      <p>{page.Name}</p>
                    </Link>
                  </li>
                )
              })
            }</ul>
          </nav>

          <Switch>
            {children}
            {
              defaultPage &&
              <Route
                path="/"
                render={() => {
                  return (
                    defaultPage &&
                    <Redirect to={GetPath(defaultPage)} />
                  )
                }}
              />
            }
          </Switch>
        </HashRouter>
      }
    </div>
  )
}

export function GetPath(page: PageEnum): string {
  const path = PagePathMap.get(page)
  if (path)
    return path

  return '/'
}