goog.provide('pr.js.profile');

goog.require('goog.dom')
goog.require('goog.net.XhrIo');


pr.js.profile = function(userId, opt_obj) {
    this.userId_ =  userId;

    this.view_ = views.infobook['profile'];

    this.updateHash_();
    
    if (!opt_obj) {
        this.fetchProfile_(userId);
    } else {
        this.obj_ = opt_obj;
        pr.js.switchView(
            this.view_, goog.bind(this.attachListeners_, this), opt_obj);
    }
}

pr.js.profile.prototype.fetchProfile_ = function(userId) {
    var param = 'id=' +  userId;
    goog.net.XhrIo.send(
        '/profile/', goog.bind(this.profileCallback_, this), 'POST', param, null);
}

pr.js.profile.prototype.profileCallback_ = function(event) {
    var xhr = event.target;
    this.obj_ = xhr.getResponseJson();
    console.log('Received: ', this.obj_);
    if (xhr.getStatus() != 200) {
        console.log('Error: ', obj);
        return;
    }
    pr.js.switchView(
        this.view_, goog.bind(this.attachListeners_, this), this.obj_);
}

pr.js.profile.prototype.updateHash_ = function() {
    document.location.hash = 'profiles/' + this.userId_;
}

pr.js.profile.prototype.attachListeners_ = function() {
    var f = function() {
        new pr.js.update(this.obj_);
    }
    goog.events.listen(
        goog.dom.getElement('update'),
        goog.events.EventType.CLICK, goog.bind(f, this, this.obj_));
}
