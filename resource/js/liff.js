const tokenCookieName = "token";
// get liff id
var liffOnChangeUserRegisters = [];

initFuncs.push(
    function() {
        var scriptSrc = 'https://static.line-scdn.net/liff/edge/2/sdk.js';
        var script = document.createElement("SCRIPT");
        script.src = scriptSrc;
        script.type = 'text/javascript';
        script.async = true;
        script.charset = 'utf-8';
        script.onload = function() {
            var scriptSrc = '/js/api/liff.js';
            var script = document.createElement("SCRIPT");
            script.src = scriptSrc;
            script.type = 'text/javascript';
            script.charset = 'utf-8';
            script.async = true;
            script.onload = function() {
                GetConfigLiff(function( liffID ) {
                    initializeLiff(liffID);
                });
            };
            loadScript(scriptSrc, script);
        };
        loadScript(scriptSrc, script);
    }
);

function loadLiffHtml() {
    $('#main').html(`
    <div id="liffAppContent">
        <!-- ACTION BUTTONS -->
        <div class="buttonGroup">
            <div class="buttonRow">
                <button id="openWindowButton">Open External Window</button>
                <button id="closeWindowButton">Close LIFF App</button>
            </div>
            <div class="buttonRow">
                <button id="sendMessageButton">Send Message</button>
                <button id="getAccessToken">Get Access Token</button>
            </div>
            <div class="buttonRow">
                <button id="getProfileButton">Get Profile</button>
                <button id="shareTargetPicker">Open Share Target Picker</button>
            </div>
        </div>
        <div id="shareTargetPickerMessage"></div>
        <!-- ACCESS TOKEN DATA -->
        <div id="accessTokenData" class="hidden textLeft">
            <h2>Access Token</h2>
            <a href="#" onclick="toggleAccessToken()">Close Access Token</a>
            <table>
                <tr>
                    <th>accessToken</th>
                    <td id="accessTokenField"></td>
                </tr>
            </table>
        </div>
        <!-- PROFILE INFO -->
        <div id="profileInfo" class="hidden textLeft">
            <h2>Profile</h2>
            <a href="#" onclick="toggleProfileData()">Close Profile</a>
            <div id="profilePictureDiv">
            </div>
            <table>
                <tr>
                    <th>userId</th>
                    <td id="userIdProfileField"></td>
                </tr>
                <tr>
                    <th>displayName</th>
                    <td id="displayNameField"></td>
                </tr>
                <tr>
                    <th>statusMessage</th>
                    <td id="statusMessageField"></td>
                </tr>
            </table>
        </div>
        <!-- LIFF DATA -->
        <div id="liffData">
            <h2 id="liffDataHeader" class="textLeft">LIFF Data</h2>
            <table>
                <tr>
                    <th>OS</th>
                    <td id="deviceOS" class="textLeft"></td>
                </tr>
                <tr>
                    <th>Language</th>
                    <td id="browserLanguage" class="textLeft"></td>
                </tr>
                <tr>
                    <th>LIFF SDK Version</th>
                    <td id="sdkVersion" class="textLeft"></td>
                </tr>
                <tr>
                    <th>LINE Version</th>
                    <td id="lineVersion" class="textLeft"></td>
                </tr>
                <tr>
                    <th>isInClient</th>
                    <td id="isInClient" class="textLeft"></td>
                </tr>
                <tr>
                    <th>isLoggedIn</th>
                    <td id="isLoggedIn" class="textLeft"></td>
                </tr>
            </table>
        </div>
        <!-- LOGIN LOGOUT BUTTONS -->
        <div class="buttonGroup">
            <button id="liffLoginButton">Log in</button>
            <button id="liffLogoutButton">Log out</button>
        </div>
        <div id="statusMessage">
            <div id="isInClientMessage"></div>
            <div id="apiReferenceMessage">
                <p>Available LIFF methods vary depending on the browser you use to open the LIFF app.</p>
                <p>Please refer to the <a href="https://developers.line.biz/en/reference/liff/#initialize-liff-app">API reference page</a> for more information.</p>
            </div>
        </div>
    </div>
    <!-- LIFF ID ERROR -->
    <div id="liffIdErrorMessage" class="hidden">
        <p>You have not assigned any value for LIFF ID.</p>
        <p>If you are running the app using Node.js, please set the LIFF ID as an environment variable in your Heroku account follwing the below steps: </p>
        <code id="code-block">
            <ol>
                <li>Go to 'Dashboard' in your Heroku account.</li>
                <li>Click on the app you just created.</li>
                <li>Click on 'Settings' and toggle 'Reveal Config Vars'.</li>
                <li>Set 'MY_LIFF_ID' as the key and the LIFF ID as the value.</li>
                <li>Your app should be up and running. Enter the URL of your app in a web browser.</li>
            </ol>
        </code>
        <p>If you are using any other platform, please add your LIFF ID in the <code>index.html</code> file.</p>
        <p>For more information about how to add your LIFF ID, see <a href="https://developers.line.biz/en/reference/liff/#initialize-liff-app">Initializing the LIFF app</a>.</p>
    </div>
    <!-- LIFF INIT ERROR -->
    <div id="liffInitErrorMessage" class="hidden">
        <p>Something went wrong with LIFF initialization.</p>
        <p>LIFF initialization can fail if a user clicks "Cancel" on the "Grant permission" screen, or if an error occurs in the process of <code>liff.init()</code>.</p>
    </div>
    <!-- NODE.JS LIFF ID ERROR -->
    <div id="nodeLiffIdErrorMessage" class="hidden">
        <p>Unable to receive the LIFF ID as an environment variable.</p>
    </div>
    `);
}

/**
* Initialize LIFF
* @param {string} myLiffId The LIFF ID of the selected element
*/
function initializeLiff(myLiffId) {
    liff
        .init({
            liffId: myLiffId
        })
        .then(() => {
            if (liff.isLoggedIn()) {
                const idToken = liff.getIDToken();
                // Set a cookie
                document.cookie = tokenCookieName + '=' + idToken + ";path=/";

                var scriptSrc = '/js/nav.js';
                var script = document.createElement("SCRIPT");
                script.src = scriptSrc;
                script.async = true;
                script.type = 'text/javascript';
                script.charset = 'utf-8';
                script.onload = function() {
                    InitNav();
                };
                loadScript(scriptSrc, script);

                const styleSrc = '/css/liff.css'; 
                // Create new link Element
                var style = document.createElement('link');
                // set the attributes for link element 
                style.rel = 'stylesheet'; 
                style.type = 'text/css';
                style.href = styleSrc; 
                loadStyle(styleSrc, style);
                loadLiffHtml();
                initializeApp()
            }else{
                sessionStorage.setItem('liffLoginRedirect', location.href);

                liff.login();
            }
        })
        .catch((err) => {
            console.log(err);
        });
}

/**
 * Initialize the app by calling functions handling individual app components
 */
function initializeApp() {
    displayLiffData();
    displayIsInClientInfo();
    registerButtonHandlers();

    // check if the user is logged in/out, and disable inappropriate button
    if (liff.isLoggedIn()) {
        document.getElementById('liffLoginButton').disabled = true;
        document.getElementById('liffLogoutButton').disabled = false;
    } else {
        document.getElementById('liffLoginButton').disabled = false;
        document.getElementById('liffLogoutButton').disabled = true;
    }
}

/**
* Display data generated by invoking LIFF methods
*/
function displayLiffData() {
    document.getElementById('browserLanguage').textContent = liff.getLanguage();
    document.getElementById('sdkVersion').textContent = liff.getVersion();
    document.getElementById('lineVersion').textContent = liff.getLineVersion();
    document.getElementById('isInClient').textContent = liff.isInClient();
    document.getElementById('isLoggedIn').textContent = liff.isLoggedIn();
    document.getElementById('deviceOS').textContent = liff.getOS();
}

/**
* Toggle the login/logout buttons based on the isInClient status, and display a message accordingly
*/
function displayIsInClientInfo() {
    if (liff.isInClient()) {
        document.getElementById('liffLoginButton').classList.toggle('hidden');
        document.getElementById('liffLogoutButton').classList.toggle('hidden');
        document.getElementById('isInClientMessage').textContent = 'You are opening the app in the in-app browser of LINE.';
    } else {
        document.getElementById('isInClientMessage').textContent = 'You are opening the app in an external browser.';
        document.getElementById('shareTargetPicker').classList.toggle('hidden');
    }
}

/**
* Register event handlers for the buttons displayed in the app
*/
function registerButtonHandlers() {
    // openWindow call
    document.getElementById('openWindowButton').addEventListener('click', function() {
        liff.openWindow({
            url: 'https://line.me',
            external: true
        });
    });

    // closeWindow call
    document.getElementById('closeWindowButton').addEventListener('click', function() {
        if (!liff.isInClient()) {
            sendAlertIfNotInClient();
        } else {
            liff.closeWindow();
        }
    });

    // sendMessages call
    document.getElementById('sendMessageButton').addEventListener('click', function() {
        if (!liff.isInClient()) {
            sendAlertIfNotInClient();
        } else {
            liff.sendMessages([{
                'type': 'text',
                'text': "You've successfully sent a message! Hooray!"
            }]).then(function() {
                window.alert('Message sent');
            }).catch(function(error) {
                window.alert('Error sending message: ' + error);
            });
        }
    });

    // get access token
    document.getElementById('getAccessToken').addEventListener('click', function() {
        if (!liff.isLoggedIn() && !liff.isInClient()) {
            alert('To get an access token, you need to be logged in. Please tap the "login" button below and try again.');
        } else {
            const accessToken = liff.getAccessToken();
            document.getElementById('accessTokenField').textContent = accessToken;
            toggleAccessToken();
        }
    });

    // get profile call
    document.getElementById('getProfileButton').addEventListener('click', function() {
        liff.getProfile().then(function(profile) {
            document.getElementById('userIdProfileField').textContent = profile.userId;
            document.getElementById('displayNameField').textContent = profile.displayName;

            const profilePictureDiv = document.getElementById('profilePictureDiv');
            if (profilePictureDiv.firstElementChild) {
                profilePictureDiv.removeChild(profilePictureDiv.firstElementChild);
            }
            const img = document.createElement('img');
            img.src = profile.pictureUrl;
            img.alt = 'Profile Picture';
            profilePictureDiv.appendChild(img);

            document.getElementById('statusMessageField').textContent = profile.statusMessage;
            toggleProfileData();
        }).catch(function(error) {
            window.alert('Error getting profile: ' + error);
        });
    });

    document.getElementById('shareTargetPicker').addEventListener('click', function () {
        if (liff.isApiAvailable('shareTargetPicker')) {
            liff.shareTargetPicker([{
                'type': 'text',
                'text': 'Hello, World!'
            }]).then(
                document.getElementById('shareTargetPickerMessage').textContent = "Share target picker was launched."
            ).catch(function (res) {
                document.getElementById('shareTargetPickerMessage').textContent = "Failed to launch share target picker.";
            });
        } else {
            document.getElementById('shareTargetPickerMessage').innerHTML = "<div>Share target picker unavailable.<div><div>This is possibly because you haven't enabled the share target picker on <a href='https://developers.line.biz/console/'>LINE Developers Console</a>.</div>";
        }
    });

    // login call, only when external browser is used
    document.getElementById('liffLoginButton').addEventListener('click', function() {
        if (!liff.isLoggedIn()) {
            // set `redirectUri` to redirect the user to a URL other than the front page of your LIFF app.
            liff.login();
        }
    });

    // logout call only when external browse
    document.getElementById('liffLogoutButton').addEventListener('click', function() {
        if (liff.isLoggedIn()) {
            liff.logout();
            
            // delete a cookie
            document.cookie = tokenCookieName + '=';

            initializeApp();
            liffOnChangeUserRegisters.forEach(register => register(liff.isLoggedIn()));
        }
    });
}

/**
* Alert the user if LIFF is opened in an external browser and unavailable buttons are tapped
*/
function sendAlertIfNotInClient() {
    alert('This button is unavailable as LIFF is currently being opened in an external browser.');
}

/**
* Toggle access token data field
*/
function toggleAccessToken() {
    toggleElement('accessTokenData');
}

/**
* Toggle profile info field
*/
function toggleProfileData() {
    toggleElement('profileInfo');
}

/**
* Toggle specified element
* @param {string} elementId The ID of the selected element
*/
function toggleElement(elementId) {
    const elem = document.getElementById(elementId);
    if (elem.offsetWidth > 0 && elem.offsetHeight > 0) {
        elem.style.display = 'none';
    } else {
        elem.style.display = 'block';
    }
}
