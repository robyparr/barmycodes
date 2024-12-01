function eh(action, data) {
  var data = data || {};

  console.log(window.EVENT_HORIZON_PROJECT_KEY, action, data)
  if (window.EVENT_HORIZON_PROJECT_KEY !== '') {
    var xhr = new XMLHttpRequest();
    xhr.open('POST', 'https://event-horizon.robyparr.com/api/events', true);
    xhr.setRequestHeader('Content-type', 'application/json; charset=UTF-8');
    xhr.setRequestHeader('Authorization', 'Bearer ' + window.EVENT_HORIZON_PROJECT_KEY);

    xhr.send(JSON.stringify({
      event: {
        action: action,
        resource: data.resource,
        resource_type: data.resourceType,
        count: data.count || 1
      }
    }));
  }
}
