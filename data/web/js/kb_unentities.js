/**
 * Created by peter on 4/02/17.
 */

var unkb = new function() {

    var self = this;
    var sessionObj = utility.getObject("session");
    var session = sessionObj.session;

    // manage the entity list on the main screen of the entity dialog system

    this.un_entity_list = [];
    this.un_entity = null;
    this.delete_name = null;

    /////////////////////////////////////////////////////////////////////////////

    // helper: convert an entity list to an html tab for display
    function un_entity_list_to_html() {
        var t_head = "<thead><tr><th>NAME</th><th></th></tr></thead>";
        var t_body = "<tbody>";
        $.each(self.un_entity_list, function(i, entity) {
            var e_id = utility.escapeHtml(entity.name);
            t_body += "<tr>";
            t_body += "<td>" + e_id + "</td>";
            t_body += "<td><div style='float: right;'>";
            t_body += "<a href='#' onclick=\"unkb.edit_un_entity('" + e_id + "');\"><img src='images/edit.png' title='edit' alt='edit'></a href='#'>&nbsp;&nbsp;";
            t_body += "<a href='#' onclick=\"unkb.delete_un_entity('" + e_id + "');\"><img src='images/delete.png' title='delete' alt='delete'></a>&nbsp;&nbsp;";
            t_body += "</div></td>";
            t_body += "</tr>";
        });
        t_body += "</tbody>";
        $("#unEntityTable").html(t_head + t_body);
    }

    function setup_un_entities(data) {
        self.un_entity_list = [];
        if (data && data.length) {
            $.each(data, function(i, entity) {
                if (entity && entity.name) {
                    self.un_entity_list.push(entity);
                }
            });
        }
        un_entity_list_to_html();
    }

    /////////////////////////////////////////////////////////////////////////////

    // get a list of the entities and display them in the page
    this.list_un_entities = function() {
        var filter_text = "null";
        $.ajax({
            url: "/topic/get_list/" + encodeURIComponent(session) + "/null/10/" + filter_text,
            type: 'GET',
            cache: false,
            dataType: 'json',
            success: function (data, textStatus, jqXHR) {
                setup_un_entities(data);
            },
            error: function (jqXHR, textStatus, errorThrown) {
                if (jqXHR && jqXHR.responseText) {
                    utility.showErrorMessage(jqXHR.responseText);
                } else {
                    utility.showErrorMessage(textStatus);
                }
            }
        });
    };

    /////////////////////////////////////////////////////////////////////////////

    // edit an existing item
    this.edit_un_entity = function(name) {
        self.un_entity = utility.getObjectFromList(self.un_entity_list, "name", name);
        if (self.un_entity) {
            $("#un_entity_name").val(self.un_entity.name);
            $("#un_entity_text").val(self.un_entity.text);
            $("#createUnstructuredEntity").modal("show");
        }
    };

    // confirm delete
    function do_delete() {
        if (self.delete_name) {
            $.ajax({
                url: "/topic/delete/" + encodeURIComponent(session) + "/" + encodeURIComponent(self.delete_name),
                type: 'DELETE',
                success: function(result) {
                    self.list_un_entities();
                },
                error: function (jqXHR, textStatus, errorThrown) {
                    if (jqXHR && jqXHR.responseText) {
                        utility.showErrorMessage(jqXHR.responseText);
                    } else {
                        utility.showErrorMessage(textStatus);
                    }
                }
            });
            utility.closeConfirmMessage();
        }
    }

    // remove an existing item
    this.delete_un_entity = function(name) {
        self.delete_name = name;
        utility.askConfirmMessage("Are you sure you want to delete '" + name + "'?", do_delete);
    };

    // create a new entity
    this.create_new_un_entity = function() {
        self.un_entity = {"name": "", "text": ""};
        $("#un_entity_name").val("");
        $("#un_entity_text").val("");
        $("#createUnstructuredEntity").modal("show");
    };

    // save an entity to the system
    this.save_un_entity = function() {
        self.un_entity.name = $("#un_entity_name").val();
        self.un_entity.text = $("#un_entity_text").val();
        if (self.un_entity && self.un_entity.name.length > 0 && self.un_entity.text.length > 0) {
            if (!self.un_entity.id) { // create a new guid
                self.un_entity.id = utility.guid();
            }
            $.post("/topic/save/" + encodeURIComponent(session) + "/" + encodeURIComponent(self.un_entity.name),
                self.un_entity.text,
                function() {
                    $("#createUnstructuredEntity").modal("hide");
                    self.list_un_entities();
                });
        } else {
            utility.showErrorMessage("invalid parameters for save");
        }
    };

};
