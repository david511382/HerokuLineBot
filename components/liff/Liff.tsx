import { useEffect } from 'react'
import { NextRouter } from 'next/router'
import { LiffExtendableFunctions } from '@line/liff/dist/lib/init/definition/LiffExtension';
import { LiffCore } from '@line/liff/dist/lib/liff';
type LiffT = LiffCore & LiffExtendableFunctions;
export type LiffType = LiffT & { [key in Exclude<string, keyof LiffT>]: any };

interface Props {
    liffID : string
    successCallback?: ((liff:LiffType) => void) | undefined
    successLoginCallback?: ((liff:LiffType) => void) | undefined
    errorCallback?: ((error: Error) => void) | undefined
    autoRedirectRouter?: NextRouter
    isAutoLogin?: boolean
}

export default function Liff(
    {
        liffID,
        successCallback,
        errorCallback,
        autoRedirectRouter,
        isAutoLogin,
        successLoginCallback,
    }: Props
) {
    const SESSION_LIFF_REDIRECT_KEY = 'liffLoginRedirect'

    useEffect(() => {
        initializeLiff(
            liffID,
            (liff)=>{
                if (successCallback){
                    successCallback(liff)
                }

                if (liff.isLoggedIn()){
                    if (successLoginCallback){
                        successLoginCallback(liff)
                    }

                    const router = autoRedirectRouter
                    if (router){
                        const { pid } = router.query
                        const isFirstRedirect = pid?.includes('liff.state')
                        if (isFirstRedirect)
                            return
                    
                        const liffLoginRedirect = sessionStorage.getItem(SESSION_LIFF_REDIRECT_KEY)
                        if (!liffLoginRedirect)
                            return
                        sessionStorage.removeItem(SESSION_LIFF_REDIRECT_KEY)
    
                        const redirectUrl = new URL(liffLoginRedirect)
                        router.push(redirectUrl)
                        return
                    }
                } else if (isAutoLogin){
                    sessionStorage.setItem(SESSION_LIFF_REDIRECT_KEY, location.href);
                    liff.login({ redirectUri:  location.href });
                }
            },
            errorCallback,
        )
    },[
        liffID,
        successCallback,
        errorCallback,
        autoRedirectRouter,
        isAutoLogin,
        successLoginCallback,
    ])

    return (null);
}

export const LoadLiffID = () =>{
    return process.env.LIFF_ID
}

async function initializeLiff(
    liffID : string,
    successCallback?: ((liff:LiffType) => void) | undefined,
    errorCallback?: ((error: Error) => void) | undefined
){
    const liff = await importLiff()
    liff.init(
            {
                liffId:liffID
            },
            ()=>{
                if (successCallback){
                    successCallback(liff)
                }
            },
            errorCallback,
        )
    return liff
}

export async function importLiff() {
    return (await import('@line/liff')).default
}