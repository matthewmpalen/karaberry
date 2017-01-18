$(document).ready(function() {
  var proto = location.protocol === 'https:' ? 'wss:' : 'ws:';
  var wsConn = new WebSocket(proto + '//' + location.host + '/ws');
  var $skipButton = $('#btnSkip');
  var $dataTable = $('#dt').DataTable();
  var $flashMessage = $('#flashMessage');
  var $history = $('#history');

  wsConn.onclose = function(evt) {
    $flashMessage.text('Websocket connection closed');
  };

  wsConn.onmessage = function(evt) {
    $flashMessage.text(evt.data);
  };

  $skipButton.click(function() {
    $.post("/skip", function(response) {
      console.log(response.message);
    });
  });

  $('#dt tbody').on('click', 'tr', function() {
    let songId = $('td:first', $(this)).text();
    let data = {songId: songId};

    $.post("/queue", data, function(response) {
      let message = response.added + " (" + response.count + " in queue)"
      console.log(message);
      //Update history
    }).fail(function(response) {
      let errorMessage = response.error;
      console.error(errorMessage);
      $flashMessage.text(errorMessage);
    });
  });
});
