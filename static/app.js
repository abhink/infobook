goog.provide('pr.start');

goog.require('goog.dom');
goog.require('goog.events');
goog.require('goog.events.EventType');


pr.start = function() {
  var newDiv = goog.dom.createDom('h1', {'style': 'background-color:#EEE'},
    'Hello world!');
  goog.dom.appendChild(document.body, newDiv);
};

goog.events.listen(window, goog.events.EventType.LOAD, pr.start);
