$(document).ready(function() {
  var $skipButton = $('#btnSkip');
  var $dataTable = $('#dt').DataTable();
  var $flashMessage = $('#flashMessage');
  var $history = $('#history');

  $skipButton.click(function() {
    $.post("/skip", function(response) {
      console.log(response.message);
      $flashMessage.text(response.message);
    });
  });

  $('#dt tbody').on('click', 'tr', function() {
    let songId = $('td:first', $(this)).text();
    let data = {songId: songId};

    $.post("/queue", data, function(response) {
      let message = response.added + " (" + response.count + " in queue)"
      console.log(message);
      $flashMessage.text(message);
      //Update history
    }).fail(function(response) {
      let errorMessage = response.error;
      console.error(errorMessage);
      $flashMessage.text(errorMessage);
    });
  });
});
