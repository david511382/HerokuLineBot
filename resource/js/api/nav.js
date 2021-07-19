function GetUserInfo(successFunc) {
    $.ajax({
            method: "GET",
            url: "/api/user-info",
        }).done(successFunc);
}