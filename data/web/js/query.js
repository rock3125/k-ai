/**
 * Created by peter on 4/02/17.
 */

var query = new function() {

    var self = this;
    var sessionObj = utility.getObject("session");
    var session = sessionObj.session;
    var name = sessionObj.first_name + " " + sessionObj.surname;
    var max_len = 40;

    this.delete_item = null; // for removing factoids
    this.factoid_registry = {}; //  factoid id -> text


    this.formatResponse = function(result_list) {
        // any previous results?  merge them
        var list = [];
        var prev_result_list = utility.getObject("result_list");
        if (prev_result_list) {
            $.each(prev_result_list, function(i, item) {
                list.push(item);
            });
        }
        if (result_list) {
            $.each(result_list, function (i, item) {
                list.push(item);
            });
        }

        // limit list size?
        if (list.length > max_len) {
            list.splice(0, list.length - max_len);
        }
        utility.setObject("result_list", list);  // save them

        self.factoid_registry = {}; // reset

        if (list.length > 0) {
            var table_str = "<thead><tr><th></th><th>topic</th><th>time</th><th></th></tr></thead><tbody>";
            $.each(list, function (i, _item) {
                var item = list[list.length - (i+1)];
                if (item && item.text) {
                    table_str += "<tr><td>" + utility.escapeHtml(item.text) + "</td><td>" +
                        utility.escapeHtml(item.topic) + "</td><td>" +
                        utility.escapeHtml(item.timestamp) + "</td>";
                    if (item.text.trim().length > 0 && item.text.indexOf('ok, got that and stored') == -1) {
                        var item_str =item.text.replace("'s", "");
                        item_str = item_str.replace("'", "");
                        table_str += "<td><img src='images/view.png' title='view a parse tree of this text.' " +
                            "onclick='query.view_tree(\"" + item_str + "\")'>";
                        // does this item have a URL?  (ie. can be deleted)
                        if (item.url) {
                            self.factoid_registry[item.url] = item.text; // keep text with id
                            table_str += "<img src='images/delete.png' title='remove this factoid' " +
                                "onclick='query.delete(\"" + item.url + "\",\"" + item.type + "\")'>";
                        }
                        table_str += "</td></tr>";
                    } else {
                        table_str += "<td></td></tr>";
                    }
                }
            });
            table_str += "</tbody>";
            return table_str;
        }
        return "";
    };

    this.view_tree = function(text) {
        if (text.indexOf(">") == 0) {
            text = text.substr(1).trim();
        }
        if (text.trim().length > 0) {
            $("#parseTree").modal("show");
            $("#parseTreeImage").attr("src", "/sl/parse-to-png/" + encodeURIComponent(text));
        }
    };

    this.display_current = function() {
        $("#responseTable").html(self.formatResponse(null));
    };

    this.process_result = function (data) {
        if (data) {
            if (data.error) {
                utility.showErrorMessage(data.error);
            } else if (data.message) {
                $("#responseTable").html(self.formatResponse([{"text": data.message}]));
            } else if (data.result_list) {
                $("#responseTable").html(self.formatResponse(data.result_list));
            } else {
                utility.showErrorMessage("no response from server");
            }
        } else {
            utility.showErrorMessage("no response from server");
        }
    };

    // populate ask field through speech to text if available
    this.ask_stt = function() {
        if (!utility.speechToText(self.ask_stt_callback)) {
            utility.showErrorMessage("Sorry, this speech to text functionality is only supported in the latest Chrome browser.")
        }
    };

    this.ask_stt_callback = function(text) {
        if (text) {
            $("#txtQuery").val(text);
        }
    };

    // perform an "ask"
    this.ask = function() {
        var query_str = $("#txtQuery").val();
        if (query_str.length > 0) {

            var obj = {"result_list": [{"text": "> " + query_str, "topic": name, "timestamp": ""}]};

            $.ajax({
                url: '/ask/' + encodeURIComponent(session),
                type: 'POST',
                data: query_str,
                cache: false,
                dataType: 'json',  // return type
                success: function (data, textStatus, jqXHR) {
                    self.process_result(data);
                    self.process_result(obj);
                },
                error: function (jqXHR, textStatus, errorThrown) {
                    if (jqXHR && jqXHR.responseText) {
                        utility.showErrorMessage(jqXHR.responseText);
                    } else {
                        utility.showErrorMessage(textStatus);
                    }
                }
            });
        }
    };

    // populate ask field through speech to text if available
    this.teach_stt = function() {
        if (!utility.speechToText(self.teach_stt_callback)) {
            utility.showErrorMessage("Sorry, this speech to text functionality is only supported in the latest Chrome browser.")
        }
    };

    this.teach_stt_callback = function(text) {
        if (text) {
            $("#txtTeach").val(text);
        }
    };

    // perform a "teach"
    this.teach = function () {
        var teach_str = $("#txtTeach").val();
        if (teach_str.length > 0) {

            var obj = {"result_list": [{"text": "> " + teach_str, "topic": name, "timestamp": ""},
                                       {"text": " ", "topic": "", "timestamp": ""}]};

            $.ajax({
                url: '/teach/' + encodeURIComponent(session),
                type: 'POST',
                data: teach_str,
                cache: false,
                dataType: 'json',
                success: function (data, textStatus, jqXHR) {
                    self.process_result(data);
                    self.process_result(obj);
                },
                error: function (jqXHR, textStatus, errorThrown) {
                    if (jqXHR && jqXHR.responseText) {
                        utility.showErrorMessage(jqXHR.responseText);
                    } else {
                        utility.showErrorMessage(textStatus);
                    }
                }
            });
        }
    };

    // confirm delete
    function do_delete() {
        if (self.delete_item && self.delete_item.url && self.delete_item.type) {
            if (self.delete_item.type != "factoid") {
                $.ajax({
                    url: "/kb-entity/delete/" + encodeURIComponent(session) + "/" +
                            encodeURIComponent(self.delete_item.type) + "/" + encodeURIComponent(self.delete_item.url),
                    type: 'DELETE',
                    success: function (result) {
                        var obj = {
                            "result_list": [{
                                "text": "kb-entry '" + self.factoid_registry[self.delete_item.url] + "' removed.",
                                "topic": name,
                                "timestamp": ""
                            },
                                {"text": " ", "topic": "", "timestamp": ""}]
                        };
                        self.process_result(obj);
                    },
                    error: function (jqXHR, textStatus, errorThrown) {
                        if (jqXHR && jqXHR.responseText) {
                            utility.showErrorMessage(jqXHR.responseText);
                        } else {
                            utility.showErrorMessage(textStatus);
                        }
                    }
                });
            } else {
                $.ajax({
                    url: "/remove/factoid/" + encodeURIComponent(session) + "/" + encodeURIComponent(self.delete_item.url),
                    type: 'DELETE',
                    success: function (result) {
                        var obj = {
                            "result_list": [{
                                "text": "factoid '" + self.factoid_registry[self.delete_item.url] + "' removed.",
                                "topic": name,
                                "timestamp": ""
                            },
                                {"text": " ", "topic": "", "timestamp": ""}]
                        };
                        self.process_result(obj);
                    },
                    error: function (jqXHR, textStatus, errorThrown) {
                        if (jqXHR && jqXHR.responseText) {
                            utility.showErrorMessage(jqXHR.responseText);
                        } else {
                            utility.showErrorMessage(textStatus);
                        }
                    }
                });
            }
            utility.closeConfirmMessage();
        }
    }

    // remove an existing factoid
    this.delete = function(url, type) {
        self.delete_item = {"url": url, "type": type};
        utility.askConfirmMessage("Are you sure you want to delete factoid '" + self.factoid_registry[url] + "'?", do_delete);
    };

    // filter / search for an entity
    this.ask_keypress = function (event) {
        if (event && event.keyCode == 13) {
            self.ask();
            return false;
        }
    };

    // filter / search for an entity
    this.teach_keypress = function(event) {
        if (event && event.keyCode == 13) {
            self.teach();
            return false;
        }
    };

};

