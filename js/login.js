goog.provide('pr.js.login');
goog.provide('pr.js.create');
goog.provide('pr.js.encodeQueryData');

goog.require('goog.dom')
goog.require('goog.net.XhrIo');


pr.js.login = function() {
    var email = goog.dom.getElement('mail').value || '';
    var pass = goog.dom.getElement('pass'). value || '';

    var callback = function(event) {
        var xhr = event.target;
        var obj = xhr.getResponseText();
        console.log('Received: ', obj);
        if (xhr.getStatus() == 401) {
            return;
        }
        new pr.js.profile(email);
    }

    var header = {
        "Authorization": "Basic " + btoa(email + ":" + pass)
    }
    var p = {
        'email': email,
        'pass': pass
    }
    goog.net.XhrIo.send(
        '/authorise', callback, 'POST', pr.js.encodeQueryData(p), header);
}

pr.js.create = function() {
    var email = goog.dom.getElement('newmail').value || '';
    var pass = goog.dom.getElement('newpass'). value || '';
    var repass = goog.dom.getElement('repass'). value || '';

    if (repass !== pass) {
        console.log("password mismatch");
        return
    }

    var callback = function(event) {
        var xhr = event.target;
        var obj = xhr.getResponseJson();
        console.log('Received: ', obj);
        if (xhr.getStatus() != 200) {
            return;
        }
        new pr.js.update(obj);
    }

    var p = {
        'email': email,
        'pass': pass
    }
    goog.net.XhrIo.send(
        '/create/', callback, 'POST', pr.js.encodeQueryData(p));
}

pr.js.encodeQueryData = function(data) {
   let ret = [];
   for (let d in data)
     ret.push(encodeURIComponent(d) + '=' + encodeURIComponent(data[d]));
   return ret.join('&');
}
