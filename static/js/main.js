$(document).ready(function () {
    $.get("/notes", function (data) {
        $.each(data, function (i, v) {
            addNoteToList(v, false)
        });
    });

    if (window.WebSocket === undefined) {
        $("#container").append("Your browser does not support WebSockets");
        return;
    } else {
        initWS();
    }

    function initWS() {
        var socket = new WebSocket("ws://localhost:8080/expire")
        socket.onmessage = function (response) {
            var data = jQuery.parseJSON(response.data);
            $.each(data, function (i, note) {
                var deleteNote = confirm("Note '" + note.title + "' has expired! Delete it?")
                if (deleteNote) {
                    $.ajax({
                        url: '/note/' + note.id,
                        type: 'delete',
                        success: function () {
                            $('li').filter('[data-id="' + note.id + '"]').remove();
                        },
                        error: function () {
                            alert("error deleting note")
                        }
                    });
                    return
                }

                note.canceled = 1;

                $.ajax({
                    url: '/note/' + note.id,
                    type: 'put',
                    dataType: 'json',
                    contentType: 'application/json',
                    success: function (note) {
                        $('a.note-link').filter('[data-id="' + note.id + '"]').text(note.title).trigger('click');
                        $("#editModal").modal('hide');
                    },
                    data: JSON.stringify(note)
                });
            });
        };
        return socket;
    }
});

$(document).on('click', '.note-link', function (e) {
    e.preventDefault();
    var id = $(this).attr('data-id');

    $.get("/note/" + id, function (note) {
        var elem = $("#note-content");
        elem.html("");
        elem.removeClass("d-none");
        elem.append("<h2>" + note.title + "</h2>")
            .append("<p class='text-left'>" + note.body + "</p>")
            .append("<p class='text-left'>Expires: " + note.expires_at + "</p>")
            .append("<p class='text-left'>Expires: " + note.expires_at + "</p>")
            .append("<button type='button' data-toggle='modal' data-target='#editModal' action='edit' data-id='" + note.id + "' class='btn btn-success btn-sm note-button'>[edit]</button>")
            .append("<button type='button' action='delete' data-id='" + note.id + "' class='btn btn-danger btn-sm ml-1 delete-note-link'>[delete]</button>");

        $('#edit-form #nf-title').val(note.title);
        $('#edit-form #nf-body').val(note.body);
        $('#edit-form #nf-id').val(note.id);
        $('#edit-form #nf-canceled').val(note.canceled);
    })
});

$(document).on('click', '#edit-btn', function (e) {
    var id = $('#edit-form #nf-id').val();
    var note = {
        body: $('#edit-form #nf-body').val(),
        title: $('#edit-form #nf-title').val(),
        canceled: $('#edit-form #nf-canceled').val()
    };

    $.ajax({
        url: '/note/' + id,
        type: 'put',
        dataType: 'json',
        contentType: 'application/json',
        success: function (note) {
            $('a.note-link').filter('[data-id="' + id + '"]').text(note.title).trigger('click');
            $("#editModal").modal('hide');
        },
        data: JSON.stringify(note)
    });
});

$('#add-btn').click(function (e) {
    e.preventDefault();

    var note = {
        title: $('#add-form #nf-title').val(),
        body: $('#add-form #nf-body').val()
    };

    $.ajax({
        url: '/note',
        type: 'post',
        dataType: 'json',
        contentType: 'application/json',
        success: function (note) {
            addNoteToList(note, true);
            $("#addModal").modal('hide');
            $('#add-form').trigger('reset');
        },
        data: JSON.stringify(note)
    });
});

$(document).on('click', '.delete-note-link', function (e) {
    var ask = confirm("Deleted note can't be restored. Are you sure?")
    if (!ask) {
        return
    }

    var id = $(this).attr('data-id');

    $.ajax({
        url: '/note/' + id,
        type: 'delete',
        success: function () {
            $('li').filter('[data-id="' + id + '"]').remove();
            $("#note-content").html("");
            $("#note-content").addClass("d-none");
        },
        error: function () {
            alert("error deleting note")
        }
    });
});

function addNoteToList(note, prepend) {
    var n = '<li class="list-group-item" data-id="' + note.id + '">' +
        '<a href="#" data-id="' + note.id + '" class="note-link float-left">' + note.title + '</a>' +
        '<a href="#" data-id="' + note.id + '" class="delete-note-link float-right badge badge-danger">delete</a>' +
        '</li>';

    if (prepend) {
        $("#notes-list").prepend(n);
        return
    }

    $("#notes-list").append(n);
}