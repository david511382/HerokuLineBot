import React, { useEffect,useState } from 'react'
import {
  HashRouter,
  Switch,
  Route,
  Redirect,
} from "react-router-dom";

export enum PageEnum  {
  Home = "Home",
}

interface Props  {
  defaultPage?: PageEnum
  children?: JSX.Element | (JSX.Element | undefined)[] 
}

const PagePathMap : Map<PageEnum,string> = new Map([
  [
    PageEnum.Home, "/臭豆腐",
  ],
]);  

export default function Nav(
  props : Props
) {
  const  [isDomMount,setDomMount]= useState(false)

  let defaultPage = PageEnum.Home
  if (props.defaultPage){
    defaultPage = props.defaultPage
  }

  useEffect(() => {
    setDomMount(true)
  },[])

  return (
    <div>
      {
        isDomMount &&
        <HashRouter>
          <Switch>
            {props.children}
            {
              props.defaultPage &&
              <Route
                  path="/"
                  render={() => {
                      return (
                        <Redirect to={getPath(defaultPage)} />
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

export function getPath(page:PageEnum) :string {
  const path = PagePathMap.get(page)
  if (path)
    return path

  return '/'
}