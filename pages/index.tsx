import Head from 'next/head'
import Image from 'next/image'
import Nav,{PageEnum} from '../components/nav/Nav'
import HomePage from '../components/home/Home'
import {Route} from "react-router-dom";
import React, { useEffect,useState } from 'react'

export default function Home() {  
  const  [isDomMount,setDomMount]= useState(false)

  useEffect(() => {
    setDomMount(true)
  },[])


  return (
    <div>
      <Head>
        <title>大台中臭豆腐批發</title>
        <meta name="description" content="大台中臭豆腐商行。臭豆腐、麻辣臭豆腐、炭烤臭豆腐台中工廠各縣市批發"/>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Nav
        defaultPage={PageEnum.Home}>
        <Route path="/臭豆腐" component={HomePage}/>
      </Nav>

      {
        !isDomMount &&
        <HomePage></HomePage>
      }

    </div>
  )
}
