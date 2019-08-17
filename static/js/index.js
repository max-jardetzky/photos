$(document).ready(function() {
    $(".container").center();
    $("#submitSignIn").click(function() {
        userdata = JSON.stringify({
            username: $("#username").val().toLowerCase(),
            password: $("#password").val()
        });
        $.ajax({
            type: "POST",
            url: "/signin",
            data: userdata,
            success: function() {
                $("#signInText").text("signed in");
                $("#signInText").css("display", "block");
                $(".container").center();
            },
            statusCode: {
                400: function() {
                    $("#signInText").text("bad request");
                    $("#signInText").css("display", "block");
                    $(".container").center();
                },
                401: function() {
                    $("#signInText").text("invalid password");
                    $("#signInText").css("display", "block");
                    $(".container").center();
                },
                404: function() {
                    $("#signInText").text("user not found");
                    $("#signInText").css("display", "block");
                    $(".container").center();
                }
            }
        })
    });
    $("#submitSignUp").click(function() {
        userdata = JSON.stringify({
            username: $("#newUsername").val().toLowerCase(),
            password: $("#newPassword").val()
        });
        $.ajax({
            type: "POST",
            url: "/signup",
            data: userdata,
            success: function() {
                $("#signUpText").text("account created");
                $("#signUpText").css("display", "block");
                $(".container").center();
            },
            statusCode: {
                400: function() {
                    $("#signUpText").text("bad request");
                    $("#signUpText").css("display", "block");
                    $(".container").center();
                },
                409: function() {
                    $("#signUpText").text("username taken");
                    $("#signUpText").css("display", "block");
                    $(".container").center();
                }
            }
        })
    });
})

$(window).on('resize', function(){
    $(".container").center();
});

jQuery.fn.center = function () {
    this.css("position","absolute");
    this.css("top", Math.max(0, (($(window).height() - $(this).outerHeight()) / 2) + 
                                                $(window).scrollTop()) + "px");
    this.css("left", Math.max(0, (($(window).width() - $(this).outerWidth()) / 2) + 
                                                $(window).scrollLeft()) + "px");
    return this;
}