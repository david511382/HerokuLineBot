const SELECTED_COLOR = 'red';
const UNSELECTED_COLOR = 'green';
const ROLE_ADMIN = 1;
const ROLE_CADRE = 2;
const ROLE_MEMBER = 3;
const ROLE_GUEST = 4;
var lastToggleObj;

function InitNav(){
    var scriptSrc = '/js/api/nav.js';
    var script = document.createElement("SCRIPT");
    script.src = scriptSrc;
    script.type = 'text/javascript';
    script.charset = 'utf-8';
    script.onload = function() {
        onUserChange(true);
    };
    loadScript(scriptSrc, script);

    const styleSrc = '/css/nav.css'; 
    // Create new link Element
    var style = document.createElement('link');
    // set the attributes for link element 
    style.rel = 'stylesheet'; 
    style.type = 'text/css';
    style.href = styleSrc; 
    loadStyle(styleSrc, style);
    
    liffOnChangeUserRegisters.push(
        onUserChange,
    );
}

function loadNavHtml() {
    $('#nav').html(`
        使用者:<font id="name"></font></br>
        <nav class="nav">
            <ul></ul>
        </nav>
    `);
}

function onUserChange(isLogin) {
    loadNavHtml();
    if (isLogin){
        loadNav();
    }
}

function loadNav() {
    GetUserInfo(function(user){
            $("#name").text(user.username);

            var navs = [];
            switch (user.role_id){
                case ROLE_ADMIN:
                    navs.push({
                        Name:"Liff",
                        Func: "loadLiffHtml",
                    })
                default:
            }
            
            navs.forEach(nav => $('nav.nav ul').append(`
            <li onclick="navToggle(this);` + nav.Func + '">' +  
                nav.Name + `
            </li>`));

            if (navs.length>0){
                $('nav.nav ul li:first-child').click();
            }
        });
}

function setSelected(obj, isSelect){
    const selectedClass = "selected";
        if (isSelect){
        $(obj).addClass(selectedClass);
    }else{
        $(obj).removeClass(selectedClass);
    }
}

function navToggle(obj){
    if (lastToggleObj){
        setSelected(lastToggleObj,false);
    }
    setSelected(obj,true);
    lastToggleObj = obj;
}


