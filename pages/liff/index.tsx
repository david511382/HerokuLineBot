import { InferGetStaticPropsType,GetStaticPropsContext} from 'next'
import Head from 'next/head'
import {GetLiffID} from '../../data/env/Liff';
import Liff,{LiffType} from '../../components/liff/Liff';
import {SetIDToken} from '../../service/auth/Token';
import React, { useState } from 'react'
import LiffPanel from '../../components/liff/LiffPanel'
import Nav from '../../components/nav/Nav'
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
        <Nav>
          <Route exact path="/debug" component={
            ()=>  {
              return LiffPanel(
                {
                  liffProps:liffProps,
                }
              )
            }
          }
          />
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
