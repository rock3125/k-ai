/**
 * Created by peter on 8/02/17.
 */

var namefrm = new function() {

    var self = this;

    // do the signin
    this.signin = function() {
        var email = $("#txtUsername").val().toLowerCase().trim();
        var password = $("#txtPassword").val().trim();
        if (password.length > 0 && utility.validateEmail(email)) {

            var obj = {'email': email, 'password_hash': password };

            $.ajax({
                url: "/user/signin",
                type: 'POST',
                data: JSON.stringify(obj),
                dataType: 'json',
                success: function (data, textStatus, jqXHR) {
                    self.session_done(data);
                },
                error: function (jqXHR, textStatus, errorThrown) {
                    if (jqXHR && jqXHR.responseText) {
                        self.session_done(JSON.parse(jqXHR.responseText));
                    } else {
                        self.session_done(JSON.parse(textStatus));
                    }
                }
            });

        } else {
            utility.showErrorMessage("Please provide us with a valid email address and password to proceed.");
        }
    };

    // show the create user password
    this.showCreateAccount = function() {
        // clear existing fields
        $("#reg_email").val("");
        $("#reg_firstname").val("");
        $("#reg_lastname").val("");
        $("#reg_password").val("");
        $("#reg_password_confirm").val("");

        $("#createAccountDlg").modal("show");
        window.setTimeout( function() {
                $("#reg_firstname").focus();
            }, 500);
        return false;
    };

    // callback from server
    this.session_done = function(data) {
        if (data && data.error) {
            utility.showErrorMessage(data.error);
        } else if (data && data.session) {
            utility.setObject("session", data);
            $("#createAccountDlg").modal("hide");
            utility.goto("query.html");
        }
    };

    // create the user from ui call
    this.createUser = function() {
        var email = $("#reg_email").val().toLowerCase().trim();
        var first_name = $("#reg_firstname").val().trim();
        var surname = $("#reg_lastname").val().trim();
        var password_1 = $("#reg_password").val().trim();
        var password_2 = $("#reg_password_confirm").val().trim();

        if (password_1.length < 8) {
            utility.showErrorMessage("Password length insufficient (8 characters minimum)");
            return;
        }
        if (password_1 != password_2) {
            utility.showErrorMessage("Passwords do not match.");
            return;
        }
        if (first_name.length == 0 || surname.length == 0) {
            utility.showErrorMessage("Please provide a first-name and surname.");
            return;
        }
        if (!utility.validateEmail(email)) {
            utility.showErrorMessage("Please provide a valid email address.");
            return;
        }
        var obj = {'email': email, 'first_name': first_name, 'surname': surname, 'password_hash': password_1 };
        // save it
        $.ajax({
            url: "/user/create",
            type: 'POST',
            data: JSON.stringify(obj),
            dataType: 'json',
            success: function (data, textStatus, jqXHR) {
                self.session_done(data);
            },
            error: function (jqXHR, textStatus, errorThrown) {
                if (jqXHR && jqXHR.responseText) {
                    self.session_done(JSON.parse(jqXHR.responseText));
                } else {
                    self.session_done(JSON.parse(textStatus));
                }
            }
        });
    };

    // enter key
    this.signin_keypress = function(event) {
        if (event && event.keyCode == 13) {
            self.signin();
            return false;
        }
    };

    // check local storage is supported - alert the user if its not
    if (!utility.localStorageSupported()) {
        utility.showErrorMessage("WARNING: your browser does not support all the required features (local-storage), this will affect the way this system performs.");
    }

};

