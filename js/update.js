goog.provide('pr.js.update');

goog.require('goog.dom')
goog.require('goog.net.XhrIo');

pr.js.update = function(data) {
    this.obj_ = data;

    this.view_ = views.infobook.update;
    
    this.updateHash_();

    pr.js.switchView(
        this.view_, goog.bind(this.attachListeners_, this), this.obj_);
}

pr.js.update.prototype.updateProfile_ = function() {
    var email = goog.dom.getElement('updatemail').value || '';
    var name = goog.dom.getElement('updatename'). value || '';
    var add = goog.dom.getElement('updateaddress'). value || '';
    var phone = goog.dom.getElement('updatephone'). value || '';
    
    var param = {
        'email': email,
        'name': name,
        'address': add,
        'phone': phone
    }
    pr.js.send(
        '/update/', goog.bind(this.updateCallback_, this), 'POST', param);
}

pr.js.update.prototype.updateCallback_ = function(response) {
    this.obj_ = response;
    new pr.js.profile(this.obj_['email'], this.obj_);
}

pr.js.update.prototype.updateHash_ = function() {
    document.location.hash = 'update/' + this.userId_;
}

pr.js.update.prototype.attachListeners_ = function() {
    goog.events.listen(
        goog.dom.getElement('update'),
        goog.events.EventType.CLICK, goog.bind(this.updateProfile_, this));
}

