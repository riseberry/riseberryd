$(function() {
  $.getJSON('/alarm', render);
  $('button').click(function() {
    $.ajax({
      type: 'PUT',
      url: '/alarm',
      data: JSON.stringify(parse()),
      contentType: "application/json; charset=utf-8",
      dataType: "json",
      success: render,
    });
  });

  var $time = $('input[type=time]');
  var $enabled = $('input[type=checkbox]');
  function parse() {
    var time = $time.val().split(':').map(Number);
    return {
      hour: time[0],
      minute: time[1],
      enabled: $enabled.prop('checked'),
    }
  }
  function render(alarm) {
    $time.val(alarm.hour+':'+alarm.minute)
    $enabled.prop('checked', alarm.enabled)
  }
})

