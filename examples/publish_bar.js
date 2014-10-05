var faye = require('faye');
var client = new faye.Client('http://localhost:5000/faye');

for(var i = 0; i < 5; i++) {
  client.publish("/bar", {hello: "bar"+i}).then(function(){})
}
