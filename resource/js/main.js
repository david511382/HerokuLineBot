
var initFuncs = [];

window.onload = function() {
    initFuncs.forEach(initFunc => initFunc());
};

function loadScript(scriptSrc, script){
    var scripts = $('head > script')
    var isExist = false;
    for (let i = 0; i < scripts.length; i++){
        let script = scripts[i];
        srcHtml = script.outerHTML;
         if (~srcHtml.indexOf(scriptSrc)){
            isExist = true
            break;
         }
    }
    if (!isExist){
        // Load the script
        $('head')[0].appendChild(script);
    }
}


function loadStyle(styleSrc, style){
    var styles = $('head > link')
    var isExist = false;
    for (let i = 0; i < styles.length; i++){
        let script = styles[i];
        srcHtml = script.outerHTML;
         if (~srcHtml.indexOf(styleSrc)){
            isExist = true
            break;
         }
    }
    if (!isExist){
        // Load the style
        $('head')[0].appendChild(style);
    }
}