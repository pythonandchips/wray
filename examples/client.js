var faye = require('faye');
var client = new faye.Client('http://localhost:5000/faye');
client.disable('websocket')

var subscription = client.subscribe("/foo", function(message){
  console.log(message)
})

subscription.then(function(a, b) {
  console.log("----------Promise response----------")
  console.log(a)
  console.log(b)
})
