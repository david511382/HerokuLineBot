function GetConfigLiff(successFunc){
    $.ajax({
            method: "GET",
            url: "/api/config/liff",
        }).done(function( response ) {
            let liffID = response.liff_id;
            successFunc(liffID);
        });
}