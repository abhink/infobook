goog.provide('pr.js');
goog.provide('pr.js.start');
goog.provide('pr.js.send');

goog.require('pr.js.login');
goog.require('pr.js.create');

goog.require('pr.js.profile');
goog.require('pr.js.update');

goog.require('goog.dom');
goog.require('goog.events');
goog.require('goog.events.EventType');
goog.require('goog.net.XhrIo');

pr.js.xsrf = "";

pr.js.start = function() {
    var hash = window.location.hash;
    
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

pr.js.send = function(url, opt_callback, opt_method, opt_content,
    opt_headers, opt_timeoutInterval) {

    var params = opt_content || {};
    if (pr.js.xsrf) {
        params['token'] = pr.js.xsrf;
    }

    var callback = function(event) {
        var xhr = event.target;
        var response = xhr.getResponseJson();
        console.log('Received: ', response);
        if (xhr.getStatus() == 401) {
            console.log('Access denied: ', response['error']);
            // window.location.href = "/";
            return;
        }
        if (xhr.getStatus() >= 400) {
            console.log('Error: ', response['error']);
            return;
        }
        pr.js.xsrf = response['token'];

        var data = response['data'];
        opt_callback(data);
    }
    
    goog.net.XhrIo.send(url, callback, opt_method,
        encodeQueryData(params), opt_headers, opt_timeoutInterval);
}

var encodeQueryData = function(data) {
   let ret = [];
   for (let d in data)
     ret.push(encodeURIComponent(d) + '=' + encodeURIComponent(data[d]));
   return ret.join('&');
}


goog.events.listen(window, goog.events.EventType.LOAD, pr.js.start);
