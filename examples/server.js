var express = require('express');
var faye = require('faye');
var http = require('http');

var app = express();
var server = http.createServer(app);
var bayeux = new faye.NodeAdapter({mount: '/faye', timeout: 45});

bayeux.attach(server);

FayeLogger = {
  incoming: function(message, callback){
    console.log("incomming");
    console.log(message)
    callback(message);
  },
  outgoing: function(message, callback){
    console.log("outgoing");
    console.log(message)
    callback(message)
  }
}

bayeux.addExtension(FayeLogger);

server.listen(5000, function(){
  console.log('Listening on port 5000');
});
