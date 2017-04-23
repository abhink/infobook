goog.provide('pr.js.login');
goog.provide('pr.js.create');
goog.provide('pr.js.encodeQueryData');

goog.require('goog.dom')
goog.require('goog.net.XhrIo');


pr.js.login = function() {
    var email = goog.dom.getElement('mail').value || '';
    var pass = goog.dom.getElement('pass'). value || '';

    var callback = function(response) {
        goog.dom.getElement('oauthlogin').innerHTML = "";
        
        var logout = goog.dom.createDom('a', null, 'Logout')
        logout.href = '/logout?email=' + response['email'];
        goog.dom.appendChild(goog.dom.getElement('logout'), logout);
        
        new pr.js.profile(response['email']);
    }

    var header = {
        "Authorization": "Basic " + btoa(email + ":" + pass)
    }
    var p = {'email': email}
    pr.js.send('/authorise', callback, 'POST', p, header);
}

pr.js.create = function() {
    var email = goog.dom.getElement('newmail').value || '';
    var pass = goog.dom.getElement('newpass'). value || '';
    var repass = goog.dom.getElement('repass'). value || '';

    if (repass !== pass) {
        console.log("password mismatch");
        return
    }

    var callback = function(response) {
        new pr.js.update(response);
    }

    var p = {
        'email': email,
        'pass': pass
    }
    pr.js.send('/create/', callback, 'POST', p);
}
