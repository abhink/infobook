goog.provide('pr.js');
goog.provide('pr.js.start');
goog.provide('pr.js.send');
goog.provide('pr.js.addLogou');

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
    var oauth = goog.dom.getElement('oauthemail');
    var token = goog.dom.getElement('oauthtoken');

    if (oauth && token) {
        var email = oauth.innerText;
        pr.js.addLogout(email);
        new pr.js.profile(email);
        return;
    }
    var hash = window.location.hash.split('/');
    
    switch(hash[0]) {
    case "":
        pr.js.switchView(views.infobook.login);
        
        goog.events.listen(
            goog.dom.getElement('login'),
            goog.events.EventType.CLICK, pr.js.login);
        goog.events.listen(
            goog.dom.getElement('create'),
            goog.events.EventType.CLICK, pr.js.create);
        break;
    case "#profiles":
    case "#update":
        var email = hash[1];
        if (!email) {
            console.log("username missing")
            return
        }
        pr.js.addLogout(email);
        new pr.js.profile(email);
        break;
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
            showError(response['error']);
            return;
        }
        if (xhr.getStatus() >= 400) {
            showError(response['error']);
            return;
        }
        showError("", true);
        pr.js.xsrf = response['token'];

        var data = response['data'];
        opt_callback(data);
    }
    
    goog.net.XhrIo.send(url, callback, opt_method,
        encodeQueryData(params), opt_headers, opt_timeoutInterval);
}

pr.js.addLogout = function(email) {
    goog.dom.getElement('oauthlogin').innerHTML = "";
    
    var logout = goog.dom.createDom('a', null, 'Logout')
    logout.href = '/logout?email=' + email;
    goog.dom.appendChild(goog.dom.getElement('logout'), logout);
}

var encodeQueryData = function(data) {
   let ret = [];
   for (let d in data)
     ret.push(encodeURIComponent(d) + '=' + encodeURIComponent(data[d]));
   return ret.join('&');
}

var showError = function(error, clearErr) {
    var errElem = goog.dom.getElement('error');
    if (clearErr) {
        errElem.innerText = "";
        return;
    }
    errElem.innerText = error;
}


goog.events.listen(window, goog.events.EventType.LOAD, pr.js.start);
