$(function() {
  var $status = $('.status');
  $.getJSON('/alarm', render);
  $('button').click(function() {
    $status.fadeOut();
    $.ajax({
      type: 'PUT',
      url: '/alarm',
      data: JSON.stringify(parse()),
      contentType: "application/json; charset=utf-8",
      dataType: "json",
      success: function(alarm) {
        render(alarm);
        $status.text('Saved').fadeIn();
      },
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
    $time.val(zeroPad(alarm.hour)+':'+zeroPad(alarm.minute))
    $enabled.prop('checked', alarm.enabled)
  }
  function zeroPad(num) {
    num = ''+num;
    return (num.length == 1)
      ? '0' + num
      : num;
  }
})

