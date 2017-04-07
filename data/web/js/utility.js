/**
 * Created by peter on 4/02/17.
 */

var utility = new function() {

    var self = this;
    var recording = false;
    var recognition = null;

    this.localMap = {};   // in case there is no local storage

    // tell the browser to goto loc
    this.goto = function(loc) {
        window.location = loc;
    };

    // signout, remove session and conversation list
    this.signout = function() {
        var sessionObj = self.getObject("session");
        if (sessionObj && sessionObj.session) {  // tell server to signout
            $.ajax({
                url: '/user/signout/' + encodeURIComponent(sessionObj.session),
                type: 'GET',
                cache: false,
                dataType: 'json'  // return type
            });
        }
        self.setObject("session", null);
        self.setObject("result_list", null);
        console.log("signing out");
        self.goto("/");
    };

    // make sure html doesn't do anything fishy
    this.escapeHtml = function(text) {
        'use strict';
        if (text && text.replace) {
            return text.replace(/[\"&'\/<>]/g, function (a) {
                return {
                    '"': '&quot;', '&': '&amp;', "'": '&#39;',
                    '/': '&#47;', '<': '&lt;', '>': '&gt;'
                }[a];
            });
        } else {
            return text
        }
    };

    // take a list of objects and remove the item that has { "name": id }
    // return all other objects that don't have this id
    this.removeObjectFromList = function(list, name, id) {
        var newList = [];
        $.each(list, function (i, obj) {
            if (obj && obj[name] != id) {
                newList.push(obj);
            }
        });
        return newList;
    };

    // get an object from a list by id,  name is the field {"name": id}
    this.getObjectFromList = function(list, name, id) {
        var selectObj = null;
        $.each(list, function (i, obj) {
            if (obj && obj[name] == id && selectObj == null) {
                selectObj = obj;
            }
        });
        return selectObj;
    };

    // create a random guid
    this.guid = function () {
        function s4() {
            return Math.floor((1 + Math.random()) * 0x10000)
                .toString(16)
                .substring(1);
        }

        return s4() + s4() + '-' + s4() + '-' + s4() + '-' + s4() + '-' + s4() + s4() + s4();
    };

    // requires errorMessage dialog and lblErrorMessage label to be present on the page
    this.showErrorMessage = function(msg) {
        if (msg.length > 0) {
            $("#lblErrorMessage").html(msg);
            $("#errorMessage").modal("show");
        }
    };

    this.closeConfirmMessage = function () {
        $("#confirmMessage").modal("hide");
    };

    // show the confirm dialog (confirmMessage) with a message (lblConfirmMessage)
    this.askConfirmMessage = function(msg, callbackYes) {
        if (msg && callbackYes) {
            $("#lblConfirmMessage").html(msg);
            $("#confirmAck").click(callbackYes);
            $("#confirmCancel1").click(self.closeConfirmMessage);
            $("#confirmCancel2").click(self.closeConfirmMessage);
            $("#confirmMessage").modal("show");
        }
    };

    /////////////////////////////////////////////////////////////////////////////

    // regex for checking email address
    this.validateEmail = function(email) {
        if (email && email.length > 0 && email.length < 50) {
            var re = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
            return re.test(email);
        }
        return false;
    };

    /////////////////////////////////////////////////////////////////////////////

    // check this person has provided a name for using K/AI
    // if not - force the user to the login.html form
    this.checkSession = function() {
        var sessionObj = self.getObject("session");
        if (!sessionObj || !sessionObj.session) {
            self.goto("login.html");
            return false;
        }
        return true;
    };

    /////////////////////////////////////////////////////////////////////////////

    // get a value from local-storage
    this.getValue = function (name) {
        if (typeof(Storage) !== "undefined") {
            if (name) {
                return localStorage.getItem(name);
            } else {
                return null;
            }
        } else {
            if (name) {
                return self.localMap[name];
            } else {
                return null;
            }
        }
    };

    // get a value from local-storage
    this.localStorageSupported = function (name) {
        return (typeof(Storage) !== "undefined");
    };

    // set a value into local-storage
    this.setValue = function (name, value) {
        if (typeof(Storage) !== "undefined") {
            if (name) {
                if (value == null) {
                    localStorage.removeItem(name);
                } else {
                    return localStorage.setItem(name, value);
                }
            } else {
                return null;
            }
        } else {
            if (name) {
                if (value == null) {
                    delete self.localMap[name];
                } else {
                    self.localMap[name] = value;
                    return value;
                }
            } else {
                return null;
            }
        }
        return null;
    };

    // set a value into local-storage
    this.getObject = function (name) {
        if (typeof(Storage) !== "undefined") {
            if (name) {
                return JSON.parse(localStorage.getItem(name));
            } else {
                return null;
            }
        } else {
            alert('Sorry! No Web Storage support...\nPlease upgrade your browser');
        }
    };

    // set a value into local-storage
    this.setObject = function (name, obj) {
        if (typeof(Storage) !== "undefined") {
            if (name) {
                if (obj == null) {
                    localStorage.removeItem(name);
                } else {
                    return localStorage.setItem(name, JSON.stringify(obj));
                }
            } else {
                return null;
            }
        } else {
            alert('Sorry! No Web Storage support...\nPlease upgrade your browser');
        }
    };

    // read some text out aloud (if supported, returns false if not supported)
    this.speakText = function(text) {
        // this is google speech
        var synth = window.speechSynthesis;
        if ( synth ) {
            var utterance = new SpeechSynthesisUtterance(text);
            utterance.lang = 'en-GB';
            synth.speak(utterance);
            return true;
        } else {
            return false;
        }
    };

    // start a speech to text session using Google STT
    this.speechToText = function(target_callback) {
        if (target_callback && ('webkitSpeechRecognition' in window)) {
            var final_transcript = '';
            if (recording) {
                recording = false;
                if (recognition) {
                    recognition.stop();
                }
            } else {
                // google speech to text
                recognition = new webkitSpeechRecognition();
                recognition.continuous = true;
                recognition.interimResults = false;
                recognition.onresult = function (event) {
                    for (var i = event.resultIndex; i < event.results.length; ++i) {
                        if (event.results[i].isFinal) {
                            final_transcript += event.results[i][0].transcript;
                        }
                    }
                };
                recognition.onend = function (event) {
                    recording = false;
                    if (target_callback) {
                        target_callback(final_transcript);
                    }
                };
                recognition.start();
                recording = true;
            }
            return true;
        } else {
            return false;
        }
    };


};

