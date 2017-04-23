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
    var param = {'email':  userId};
    pr.js.send(
        '/profile/', goog.bind(this.profileCallback_, this), 'POST', param);
}

pr.js.profile.prototype.profileCallback_ = function(response) {
    this.obj_ = response;
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
