/**
 * Created by peter on 4/02/17.
 */

var kb = new function() {

    var self = this;
    var sessionObj = utility.getObject("session");
    var session = sessionObj.session;

    // manage the entity list on the main screen of the entity dialog system

    this.entity_list = [];
    this.entity = null;
    this.field = null;
    this.files = null;
    this.upload_guid = null;
    this.delete_name = null;

    /////////////////////////////////////////////////////////////////////////////

    // helper: convert an entity list to an html tab for display
    function entity_list_to_html() {
        var t_head = "<thead><tr><th>NAME</th><th></th></tr></thead>";
        var t_body = "<tbody>";
        $.each(self.entity_list, function(i, entity) {
            var e_id = utility.escapeHtml(entity.id);
            t_body += "<tr>";
            t_body += "<td>" + utility.escapeHtml(entity.name) + "</td>";
            t_body += "<td><div style='float: right;'>";
            t_body += "<a href='#' onclick=\"kb.edit_entity('" + e_id + "');\"><img src='images/edit.png' title='edit' alt='edit'></a href='#'>&nbsp;&nbsp;";
            t_body += "<a href='#' onclick=\"kb.delete_entity('" + e_id + "','" + entity.name + "');\"><img src='images/delete.png' title='delete' alt='delete'></a>&nbsp;&nbsp;";
            t_body += "<a href='#' onclick=\"kb.upload_entities('" + e_id + "');\"><img src='images/upload.png' title='upload' alt='upload'></a>";
            t_body += "</div></td>";
            t_body += "</tr>";
        });
        t_body += "</tbody>";
        $("#entityTable").html(t_head + t_body);
    }

    function setup_entities(data) {
        self.entity_list = [];
        if (data && data.length) {
            $.each(data, function(i, kb_entity) {
                var obj;
                if (kb_entity && kb_entity.json_data) {
                    obj = JSON.parse(kb_entity.json_data);
                    self.entity_list.push(obj);
                } else if (kb_entity && kb_entity.Json_data) {
                    obj = JSON.parse(kb_entity.Json_data);
                    self.entity_list.push(obj);
                }
            });
        }
        entity_list_to_html();
    }

    // helper: convert an field list to an html tab for display
    function field_list_to_html() {
        var t_head = "<thead><tr><th>field</th><th>semantic</th><th></th></tr></thead>";
        var t_body = "<tbody>";
        $.each(self.entity.field_list, function(i, field) {
            var e_id = utility.escapeHtml(field.name);
            var e_sem = utility.escapeHtml(field.semantic);
            t_body += "<tr>";
            t_body += "<td>" + e_id + "</td>";
            t_body += "<td>" + e_sem + "</td>";
            t_body += "<td><div style='float: right;'>";
            t_body += "<a href='#' onclick=\"kb.edit_field('" + e_id + "');\"><img src='images/edit.png' title='edit' alt='edit'></a href='#'>&nbsp;&nbsp;";
            t_body += "<a href='#' onclick=\"kb.delete_field('" + e_id + "');\"><img src='images/delete.png' title='delete' alt='delete'></a>";
            t_body += "</div></td>";
            t_body += "</tr>";
        });
        t_body += "</tbody>";
        $("#fieldTable").html(t_head + t_body);
    }

    /////////////////////////////////////////////////////////////////////////////

    // get a list of the entities and display them in the page
    this.list_kb_entities = function() {
        var filter_text = "null";

        $.ajax({
            url: "/kb-entity/get_list/" + encodeURIComponent(session) + "/schema/null/10/name/" + filter_text,
            type: 'GET',
            cache: false,
            dataType: 'json',
            success: function (data, textStatus, jqXHR) {
                setup_entities(data);
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
    this.edit_entity = function(id) {
        self.entity = utility.getObjectFromList(self.entity_list, "id", id);
        if (self.entity) {
            if (!self.entity.entity_list) {
                self.entity.entity_list = [];
            }
            $("#entity_name").val(self.entity.name);
            $("#aiml_text").val(self.entity.aiml);
            field_list_to_html();

            $("#createEntity").modal("show");
        }
    };

    // confirm delete
    function do_delete() {
        if (self.delete_name) {
            $.ajax({
                url: "/kb-entity/delete/" + encodeURIComponent(session) + "/schema/" +
                            encodeURIComponent(self.delete_name),
                type: 'DELETE',
                success: function(result) {
                    self.list_kb_entities();
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
    this.delete_entity = function(id,name) {
        self.delete_name = id;
        utility.askConfirmMessage("Are you sure you want to delete '" + name + "'?", do_delete);
    };

    // create a new entity
    this.create_new_entity = function() {
        self.entity = {"name": "", "field_list": []};

        $("#entity_name").val("");
        $("#aiml_text").val("");
        field_list_to_html();

        $("#createEntity").modal("show");
    };

    // remove an existing item
    this.delete_field = function(name) {
        self.entity.field_list = utility.removeObjectFromList(entity.field_list, "name", name);
        field_list_to_html();
    };

    // save an entity to the system
    this.save_entity = function() {
        self.entity.name = $("#entity_name").val();
        if (self.entity && self.entity.name.length > 0 && self.entity.field_list && self.entity.field_list.length > 0) {
            if (!self.entity.id) { // create a new guid
                self.entity.id = utility.guid();
            }
            var kb_entity = {  // warp it into a kb entity
                "id": self.entity.id,
                "json_data": JSON.stringify(self.entity),
                "topic": "schema"
            };
            var kb_entity_json = JSON.stringify(kb_entity); // to string
            $.post("/kb-entity/save/" + encodeURIComponent(session), kb_entity_json,
                function() {
                    $("#createEntity").modal("hide");
                    self.list_kb_entities();
                });
        } else {
            utility.showErrorMessage("invalid parameters for save");
        }
    };

    /////////////////////////////////////////////////////////////////////////////

    // edit an existing item
    this.edit_field = function(name) {
        self.field = utility.getObjectFromList(self.entity.field_list, "name", name);
        if (self.field) {
            $("#field_name").val(self.field.name);
            $("#field_semantic").val(self.field.semantic);
            $("#aiml_text").val(self.field.aiml);
            $("#createField").modal("show");
        }
    };

    this.create_new_field = function() {
        self.field = null;
        $("#field_name").val("");
        $("#field_semantic").val("");
        $("#aiml_text").val("");
        $("#createField").modal("show");
    };

    this.save_field = function() {
        var name = $("#field_name").val();
        var semantic = $("#field_semantic").val();
        var aiml = $("#aiml_text").val();
        if (name.length > 0 && semantic.length > 0) {
            self.field = utility.getObjectFromList(self.entity.field_list, "name", name);
            if (!self.field) {
                self.entity.field_list.push({"name": name, "semantic": semantic, "aiml": aiml});
            } else {
                self.field.name = name;
                self.field.semantic = semantic;
                self.field.aiml = aiml;
            }
            $("#createField").modal("hide");
            field_list_to_html();
        } else {
            utility.showErrorMessage("you must provide a name and a semantic for a field")
        }
    };

    /////////////////////////////////////////////////////////////////////////////

    // show upload dialog for an entity
    this.upload_entities = function(id) {
        self.upload_guid = id;
        self.files = null;
        $("#uploadInstances").modal("show");
    };

    // prepare files for upload (onchange callback)
    this.prepareUpload = function(event) {
        self.files = event.target.files;
    };

    // private helper
    function upload_done() {
        $("#busy").modal("hide");
    }

    // do upload (click upload button)
    this.upload_file = function() {
        var val = $("#documentUpload").val();
        if (val.length == 0 || self.files == null) {
            utility.showErrorMessage("please select a file to upload first");
        } else {
            var data = new FormData();
            $.each(self.files, function(key, value) {
                data.append(key, value);
            });
            $("#uploadInstances").modal("hide");
            $("#busy").modal("show");
            $.ajax({
                url: '/kb-entity/upload/' + encodeURIComponent(session) + '/' + encodeURIComponent(self.upload_guid),
                type: 'POST',
                data: data,
                cache: false,
                dataType: 'multipart/form-data',
                processData: false, // Don't process the files
                contentType: false, // Set content type to false as jQuery will tell the server its a query string request
                success: function(data, textStatus, jqXHR) {
                    upload_done();
                },
                error: function(jqXHR, textStatus, errorThrown) {
                    upload_done();
                }
            });
        }
    };

};
