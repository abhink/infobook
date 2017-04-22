goog.provide('pr.js');
goog.provide('pr.js.start');

goog.require('pr.js.login');
goog.require('pr.js.create');

goog.require('pr.js.profile');
goog.require('pr.js.update');

goog.require('goog.dom');
goog.require('goog.events');
goog.require('goog.events.EventType');
goog.require('goog.net.XhrIo');


pr.js.start = function() {
    var newDiv = goog.dom.createDom('h1', {'style': 'background-color:#EEE'},
        'Hello world!');
    goog.dom.appendChild(document.body, newDiv);

    switch(window.location.pathname) {
    case "/":
        pr.js.switchView(views.infobook.login);
        
        goog.events.listen(
            goog.dom.getElement('login'),
            goog.events.EventType.CLICK, pr.js.login);
        goog.events.listen(
            goog.dom.getElement('create'),
            goog.events.EventType.CLICK, pr.js.create);
    }
    
};

pr.js.switchView = function(view, opt_attachListenerFunc_, opt_param) {
    var viewElem = goog.dom.getElement('view');
    viewElem.innerHTML = view(opt_param);
    if (opt_attachListenerFunc_) {
        opt_attachListenerFunc_();
    }
}

goog.events.listen(window, goog.events.EventType.LOAD, pr.js.start);
