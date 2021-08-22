import { InferGetStaticPropsType,GetStaticPropsContext} from 'next'
import Head from 'next/head'
import {GetLiffID} from '../../data/env/Liff';
import Liff,{LiffType} from '../../components/liff/Liff';
import {SetIDToken} from '../../service/auth/Token';
import React, { useState } from 'react'
import LiffPanel from '../../hashPages/liffPanenl/LiffPanel'
import Nav,{PageEnum,GetPath} from '../../components/nav/Nav'
import RentalCourt from '../../hashPages/clubRentalCourt/RentalCourt'
import {Route} from "react-router-dom";

export default function Page({liffID}: InferGetStaticPropsType<typeof getStaticProps>) {
  const  [liffProps,setLiff]= useState<LiffType>()
  
  if (!liffID){
    return (null)
  }

  return (
    <div>
      <Head>
        <title>羽球</title>
        <meta httpEquiv="Content-Type" content="text/html; charset=utf-8"/>
        <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <Liff
        liffID={liffID}
        isAutoLogin={true}
        successLoginCallback={(liff:LiffType)=>{
          setLiff(liff)

          const idToken = liff.getIDToken();
          if (!idToken)
            return
          SetIDToken(idToken)
        }}
        errorCallback={(err:any)=>{
            console.log(err)}
        }
      />
      {
        liffProps &&
        <Nav defaultPage={PageEnum.RentalCourts}>
          <Route exact path={GetPath(PageEnum.LiffDebug)}
            component={
              ()=>  {
                return LiffPanel(
                  {
                    liffProps:liffProps,
                  }
                )
              }
            }
          />
          <Route exact path={GetPath(PageEnum.RentalCourts)} component={RentalCourt}></Route>
        </Nav>
      }
    </div>
  )
}

export const getStaticProps = (
  context: GetStaticPropsContext
) => {
  return {
      props: {
        liffID: GetLiffID(),
      },
  }
}