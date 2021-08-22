import { GetStaticProps, InferGetStaticPropsType } from 'next'
import React, { useEffect,useState } from 'react'
import styles from './liff.module.css'
import {LiffType} from '../../components/liff/Liff';

export default function LiffPanel(props: InferGetStaticPropsType<GetStaticProps>) {
    const liffProps:LiffType  = props.liffProps
    const  [isLogin,setIsLogin]= useState<boolean>(liffProps.isLoggedIn())
    const IS_IN_CLIENT = liffProps.isInClient()
    useEffect(() => {
        displayIsInClientInfo(IS_IN_CLIENT);
    },[IS_IN_CLIENT])
    
    return (
        <div>
            <div id="liffAppContent" className={`${styles.div}`}>
                {/* <!-- ACTION BUTTONS --> */}
                <div className={`${styles.div} ${styles.buttonGroup}`}>
                    <div className={`${styles.div} ${styles.buttonRow}`}>
                        <button id="openWindowButton"
                            className={`${styles.button}`}
                            onClick={()=>{
                                liffProps.openWindow({
                                    url: 'https://line.me',
                                    external: true
                                });
                            }}>
                                Open External Window
                        </button>
                        <button id="closeWindowButton"
                            className={`${styles.button}`}
                            onClick={()=>{
                                if (!IS_IN_CLIENT) {
                                    sendAlertIfNotInClient();
                                } else {
                                    liffProps.closeWindow();
                                }
                            }}>
                                Close LIFF App
                        </button>
                    </div>
                    <div className={`${styles.div} ${styles.buttonRow}`}>
                        <button id="sendMessageButton"
                            className={`${styles.button}`}
                            onClick={()=>{
                                if (!IS_IN_CLIENT) {
                                    sendAlertIfNotInClient();
                                } else {
                                    liffProps.sendMessages([{
                                        'type': 'text',
                                        'text': "You've successfully sent a message! Hooray!"
                                    }]).then(function() {
                                        window.alert('Message sent');
                                    }).catch(function(error :any ) {
                                        window.alert('Error sending message: ' + error);
                                    });
                                }
                            }}>
                                Send Message
                        </button>
                        <button id="getAccessToken"
                            className={`${styles.button}`}
                            onClick={()=>{
                                if (!isLogin && !IS_IN_CLIENT) {
                                    alert('To get an access token, you need to be logged in. Please tap the "login" button below and try again.');
                                } else {
                                    const accessToken = liffProps.getAccessToken();
                                    const accessTokenField = document.getElementById('accessTokenField')
                                    if (accessTokenField)
                                        accessTokenField.textContent = accessToken;

                                    toggleAccessToken();
                                }
                            }}>
                                Get Access Token
                        </button>
                    </div>
                    <div className={`${styles.div} ${styles.buttonRow}`}>
                        <button id="getProfileButton"
                            className={`${styles.button}`}
                            onClick={()=>{
                                liffProps.getProfile().then(function(profile) {
                                    const userIdProfileField = document.getElementById('userIdProfileField')
                                    if (userIdProfileField)
                                        userIdProfileField.textContent = profile.userId;

                                    const displayNameField = document.getElementById('displayNameField')
                                    if (displayNameField)
                                        displayNameField.textContent = profile.displayName

                                    const statusMessageField = document.getElementById('statusMessageField')
                                    if (statusMessageField &&  profile.statusMessage)
                                        statusMessageField.textContent  = profile.statusMessage;
                                    
                                    const profilePictureDiv = document.getElementById('profilePictureDiv');
                                    if (profilePictureDiv) {
                                        if (profilePictureDiv.firstElementChild)
                                            profilePictureDiv.removeChild(profilePictureDiv.firstElementChild);
                                            
                                        const img = document.createElement('img');
                                        if (profile.pictureUrl){
                                            img.src = profile.pictureUrl;
                                        }
                                        img.alt = 'Profile Picture';
                                        img.className = 'img'
                                        profilePictureDiv.appendChild(img);
                                    }
                                    
                                    toggleProfileData();
                                }).catch(function(error) {
                                    window.alert('Error getting profile: ' + error);
                                });
                            }}>
                                Get Profile</button>
                        <button id="shareTargetPicker"
                            className={`${styles.button}`}
                            onClick={()=>{
                                if (liffProps.isApiAvailable('shareTargetPicker')) {
                                    liffProps.shareTargetPicker([{
                                        'type': 'text',
                                        'text': 'Hello, World!'
                                    }]).then(()=>{
                                        const shareTargetPickerMessage = document.getElementById('shareTargetPickerMessage')
                                        if (shareTargetPickerMessage)
                                            shareTargetPickerMessage.textContent = "Share target picker was launched."
                                    }).catch(function (res) {
                                        const shareTargetPickerMessage = document.getElementById('shareTargetPickerMessage')
                                        if (shareTargetPickerMessage)
                                            shareTargetPickerMessage.textContent = "Failed to launch share target picker.";
                                    });
                                } else {
                                    const shareTargetPickerMessage = document.getElementById('shareTargetPickerMessage')
                                        if (shareTargetPickerMessage)
                                            shareTargetPickerMessage.innerHTML = "<div>Share target picker unavailable.<div><div>This is possibly because you haven't enabled the share target picker on <a href='https://developers.line.biz/console/'>LINE Developers Console</a>.</div>";
                                }
                            }}>
                                Open Share Target Picker</button>
                    </div>
                </div>
                <div id="shareTargetPickerMessage" className={`${styles.div}`}></div>
                {/* <!-- ACCESS TOKEN DATA --> */}
                <div id="accessTokenData" className={`${styles.div} ${styles.hidden} ${styles.textLeft}`}>
                    <h2>Access Token</h2>
                    <a href="#" onClick={()=>{toggleAccessToken()}}>Close Access Token</a>
                    <table className={styles.table}>
                        <tbody>
                            <tr>
                                <th className={styles.th}>accessToken</th>
                                <td id="accessTokenField" className={styles.td}></td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                {/* <!-- PROFILE INFO --> */}
                <div id="profileInfo" className={`${styles.div} ${styles.hidden} ${styles.textLeft}`}>
                    <h2>Profile</h2>
                    <a href="#" onClick={()=>{toggleProfileData()}}>Close Profile</a>
                    <a href="#" >Close Access Token</a>
                    <div id="profilePictureDiv" className={`${styles.div}`}>
                    </div>
                    <table className={styles.table}>
                        <tbody>
                            <tr>
                                <th className={styles.th}>userId</th>
                                <td id="userIdProfileField" className={styles.td}></td>
                            </tr>
                            <tr>
                                <th className={styles.th}>displayName</th>
                                <td id="displayNameField" className={styles.td}></td>
                            </tr>
                            <tr>
                                <th className={styles.th}>statusMessage</th>
                                <td id="statusMessageField" className={styles.td}></td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                {/* <!-- LIFF DATA --> */}
                <div id="liffData" className={`${styles.div}`}>
                    <h2 id="liffDataHeader" className={styles.textLeft}>LIFF Data</h2>
                    <table className={styles.table}>
                        <tbody>
                            <tr>
                                <th className={styles.th}>OS</th>
                                <td id="deviceOS" className={`${styles.textLeft} ${styles.td}`}>{liffProps.getOS()}</td>
                            </tr>
                            <tr>
                                <th className={styles.th}>Language</th>
                                <td id="browserLanguage" className={`${styles.textLeft} ${styles.td}`}>{liffProps.getLanguage()}</td>
                            </tr>
                            <tr>
                                <th className={styles.th}>LIFF SDK Version</th>
                                <td id="sdkVersion" className={`${styles.textLeft} ${styles.td}`}>{liffProps.getVersion()}</td>
                            </tr>
                            <tr>
                                <th className={styles.th}>LINE Version</th>
                                <td id="lineVersion" className={`${styles.textLeft} ${styles.td}`}>{liffProps.getLineVersion()}</td>
                            </tr>
                            <tr>
                                <th className={styles.th}>isInClient</th>
                                <td id="isInClient" className={`${styles.textLeft} ${styles.td}`}>{(IS_IN_CLIENT)?"true":"false"}</td>
                            </tr>
                            <tr>
                                <th className={styles.th}>isLoggedIn</th>
                                <td id="isLoggedIn" className={`${styles.textLeft} ${styles.td}`}>{(isLogin)?"true":"false"}</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                {/* <!-- LOGIN LOGOUT BUTTONS --> */}
                <div className={`${styles.div} ${styles.buttonGroup}`}>
                    {/* <button className={styles.button" onClick={sendMessage}>send message</button>  */}
                    <button id="liffLoginButton"
                        className={`${styles.button}`}
                        onClick={()=>{
                            if (!isLogin) {
                                // set `redirectUri` to redirect the user to a URL other than the front page of your LIFF app.
                                liffProps.login();
                                setIsLogin(liffProps.isLoggedIn())
                            }
                        }}
                        disabled={isLogin}>
                            Log in</button>
                    <button id="liffLogoutButton"
                        className={`${styles.button}`}
                        onClick={()=>{
                            if (isLogin) {
                                liffProps.logout();
                                setIsLogin(liffProps.isLoggedIn())
                                // delete a cookie
                                // document.cookie = tokenCookieName + '=';
                    
                                // initializeApp();
                                // liffOnChangeUserRegisters.forEach(register => register(liff.isLoggedIn()));
                            }
                        }}
                        disabled={!isLogin}>
                            Log out</button>
                </div>
                <div id="statusMessage" className={`${styles.div}`}>
                    <div id="isInClientMessage" className={`${styles.div}`}>
                        {
                            (IS_IN_CLIENT)? 
                                'You are opening the app in the in-app browser of LINE.':
                                'You are opening the app in an external browser.'
                        }
                    </div>
                    <div id="apiReferenceMessage" className={`${styles.div}`}>
                        <p>Available LIFF methods vary depending on the browser you use to open the LIFF app.</p>
                        <p>Please refer to the <a href="https://developers.line.biz/en/reference/liff/#initialize-liff-app">API reference page</a> for more information.</p>
                    </div>
                </div>
            </div>

            {/* <!-- LIFF ID ERROR --> */}
            <div id="liffIdErrorMessage" className={`${styles.div} ${styles.hidden}`}>
                <p>You have not assigned any value for LIFF ID.</p>
                <p>If you are running the app using Node.js, please set the LIFF ID as an environment variable in your Heroku account follwing the below steps: </p>
                <code id="code-block">
                    <ol>
                        <li>{"Go to 'Dashboard' in your Heroku account."}</li>
                        <li>{"Click on the app you just created."}</li>
                        <li>{"Click on 'Settings' and toggle 'Reveal Config Vars'."}</li>
                        <li>{"Set 'MY_LIFF_ID' as the key and the LIFF ID as the value."}</li>
                        <li>{"Your app should be up and running. Enter the URL of your app in a web browser."}</li>
                    </ol>
                </code>
                <p>If you are using any other platform, please add your LIFF ID in the <code>index.html</code> file.</p>
                <p>For more information about how to add your LIFF ID, see <a href="https://developers.line.biz/en/reference/liff/#initialize-liff-app">Initializing the LIFF app</a>.</p>
            </div>

            {/* <!-- LIFF INIT ERROR --> */}
            <div id="liffInitErrorMessage" className={`${styles.div} ${styles.hidden}`}>
                <p>{`Something went wrong with LIFF initialization.`}</p>
                <p>{`LIFF initialization can fail if a user clicks "Cancel" on the "Grant permission" screen, or if an error occurs in the process of `}<code>liff.init()</code>.</p>
            </div>
            
            {/* <!-- NODE.JS LIFF ID ERROR --> */}
            <div id="nodeLiffIdErrorMessage" className={`${styles.div} ${styles.hidden}`}>
                <p>Unable to receive the LIFF ID as an environment variable.</p>
            </div>
        </div>
    )
}

/**
* Toggle the login/logout buttons based on the isInClient status, and display a message accordingly
*/
function displayIsInClientInfo(isInClient :boolean) {
  if (isInClient) {
      document.getElementById('liffLoginButton')?.classList.toggle('hidden');
      document.getElementById('liffLogoutButton')?.classList.toggle('hidden');
  } else {
      document.getElementById('shareTargetPicker')?.classList.toggle('hidden');
  }
}

/**
* Alert the user if LIFF is opened in an external browser and unavailable buttons are tapped
*/
function sendAlertIfNotInClient() {
  alert('This button is unavailable as LIFF is currently being opened in an external browser.');
}


function toggleAccessToken() {
  toggleElement('accessTokenData');
}

function toggleProfileData() {
  toggleElement('profileInfo');
}

function toggleElement(elementId:string) {
  const elem = document.getElementById(elementId);
  if (elem && elem.offsetWidth > 0 && elem.offsetHeight > 0) {
      elem.style.display = 'none';
  } else if (elem){
      elem.style.display = 'block';
  }
}
