/**
 * Created by peter on 7/02/17.
 */

var entities = new function() {

    var self = this;
    var sessionObj = utility.getObject("session");
    var session = sessionObj.session;

    this.entity_list = [];

    // for removing entities
    this.delete_name = null;

    // helper: convert an field list to an html tab for display
    function entity_list_to_html() {
        var t_head = "<thead><tr><th>NAME</th><th>SEMANTIC</th><th></th></tr></thead>";
        var t_body = "<tbody>";
        $.each(self.entity_list, function(i, entity) {
            var e_id = utility.escapeHtml(entity.name);
            var e_sem = utility.escapeHtml(entity.semantic);
            t_body += "<tr>";
            t_body += "<td>" + e_id + "</td>";
            t_body += "<td>" + e_sem + "</td>";
            t_body += "<td><div style='float: right;'>";
            t_body += "<a href='#' onclick=\"entities.edit_entity('" + e_id + "');\"><img src='images/edit.png' title='edit' alt='edit'></a href='#'>&nbsp;&nbsp;";
            t_body += "<a href='#' onclick=\"entities.delete_entity('" + e_id + "');\"><img src='images/delete.png' title='delete' alt='delete'></a>";
            t_body += "</div></td>";
            t_body += "</tr>";
        });
        t_body += "</tbody>";
        $("#entityTable").html(t_head + t_body);
    }


    /////////////////////////////////////////////////////////////////////////////

    this.process_result = function (data) {
        if (data) {
            if (data.error) {
                utility.showErrorMessage(data.error);
            } else if (typeof data === 'string' || data instanceof String) {
                if (data != "ok") {
                    utility.showErrorMessage(data);
                }
            } else {
                utility.showErrorMessage("no response from server");
            }
        } else {
            utility.showErrorMessage("no response from server");
        }
    };

    /////////////////////////////////////////////////////////////////////////////

    // find from UI
    this.find = function() {
        var name = $("#txtFindEntity").val();
        if (name.length > 0) {
            $.get("/entities/find/" + encodeURIComponent(session) + "/" + encodeURIComponent(name),
                function (data) {
                    if (typeof data === 'string' || data instanceof String) {
                        self.process_result(data);
                    } else if (data && data.length >= 0) {
                        self.entity_list = data;
                        entity_list_to_html();
                    }
                });
        }
    };

    // onkeypress from UI
    this.find_keypress = function(event) {
        if (event && event.keyCode == 13) {
            self.find();
            return false;
        }
    };

    // ui save
    this.save_entity = function() {
        var name = $("#entity_name").val();
        var semantic = $("#entity_semantic").val();
        if (name.length > 0 && semantic.length > 0) {
            $.get("/entities/save/" + encodeURIComponent(session) + "/" + encodeURIComponent(name) + "/" + encodeURIComponent(semantic),
                function (data) {
                    $("#createEntity").modal("hide");
                    self.find();
                });
        } else {
            utility.showErrorMessage("You must supply a name and a semantic when saving an entity.");
        }
    };

    this.edit_entity = function(name) {
        var entity = utility.getObjectFromList(self.entity_list, "name", name);
        if (entity) {
            $("#entity_name").val(entity.name);
            $("#entity_semantic").val(entity.semantic);
            $("#createEntityTitle").html("EDIT ENTITY");
            $("#createEntity").modal("show");
        }
    };

    // create a new entity
    this.create_new_entity = function () {
        $("#entity_name").val("");
        $("#entity_semantic").val("");
        $("#createEntity").modal("show");
    };

    // confirm delete
    function do_delete() {
        if (self.delete_name) {
            $.ajax({
                url: "/entities/delete/" + encodeURIComponent(session) + "/" + encodeURIComponent(self.delete_name),
                type: 'DELETE',
                success: function(data) {
                    if (data && data.responseText) {
                        self.process_result(JSON.parse(data.responseText));
                    } else if (data && data.message) {
                        self.process_result(data.message);
                    }
                    utility.closeConfirmMessage();
                    self.find();
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
    }

    this.delete_entity = function(name) {
        self.delete_name = name;
        utility.askConfirmMessage("Are you sure you want to delete '" + name + "'?", do_delete);
    };

};

