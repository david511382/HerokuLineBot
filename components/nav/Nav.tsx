import React, { useEffect,useState } from 'react'
import styles from './nav.module.css'
import {GetUserInfo} from '../../service/auth/User'
import { RoleID,  UserInfo as UserInfoResp } from '../../models/resp/user-info'
import {Link,HashRouter, Switch} from "react-router-dom";

export enum PageEnum  {
  LiffDebug = "Liff Debug",
}

interface Props  {
  defaultPage?: PageEnum
  children?: JSX.Element | (JSX.Element | undefined)[] 
}

type Page ={
  Name: PageEnum
  Path: string
  Selected: boolean
}

const PagePathMap : Map<PageEnum,string> = new Map([
  [
    PageEnum.LiffDebug, "/debug",
  ],
]);  

const RolePagesMap : Map<RoleID,PageEnum[]> = new Map([
  [
    RoleID.ROLE_ADMIN, 
    [
      PageEnum.LiffDebug,
    ],
  ],
]);  

export default function Nav(
  {
    children,
    defaultPage
  }: Props
) {
  const  [userInfo,setUserInfo]= useState<UserInfoResp>()
  const  [navs,setNavs]= useState<Page[]>([])
  
  useEffect(() => {
    if (!defaultPage){
      let currentHash = window.location.hash.replaceAll("#","")
      PagePathMap.forEach((path,page)=>{
        if (defaultPage)
          return
        defaultPage = (currentHash === path)?
          page:
          undefined
      })
    }
    
    GetUserInfo()
      .then((userInfo)=>{
        if (!userInfo)
          return

        setUserInfo(userInfo)
        
        const roleID = userInfo.role_id
        let pages = RolePagesMap.get(roleID)
        if (!pages){
          pages = []
        }

        let pageViews :Page[] = []
        pages.forEach((page)=>pageViews.push({
          Name: page,
          Path: getPath(page),
          Selected: (defaultPage && page === defaultPage)? true : false,
        }))
        
        setNavs(pageViews)
      })
  },[
    defaultPage
  ])
  
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
                    className={`${styles.li} ${(page.Selected)?`${styles.selected}`:""}`}
                    key={page.Name}>
                    <Link
                      to={page.Path}
                      onClick={()=>{
                        navs.forEach((page)=>page.Selected=false)
                        navs[i].Selected = true
                        setNavs([...navs])
                      }}>
                      <a>{page.Name}</a>
                    </Link>
                  </li>
                ) 
              })
            }</ul>
          </nav>

          <Switch>
            {children}
          </Switch>
        </HashRouter>
      }
    </div>
  )
}

export function getPath(page:PageEnum) :string {
  const path = PagePathMap.get(page)
  if (path)
    return path

  return '/'
}