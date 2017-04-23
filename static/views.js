// This file was automatically generated from views.soy.
// Please don't edit this file by hand.

/**
 * @fileoverview Templates in namespace views.infobook.
 */

if (typeof views == 'undefined') { var views = {}; }
if (typeof views.infobook == 'undefined') { views.infobook = {}; }


views.infobook.login = function(opt_data, opt_ignored) {
  return '<div id="loginview"><h3>Existing users, log in below:</h3><div><label for="mail">E-mail:</label><input type="email" id="mail" name="user_mail"></div><div><label for="pass">Password:</label><input type="password" id="pass" name="user_password"></div><div id="login" class="button"><button type="submit">Login</button></div><br><br><h3>New users, please enter below details to create a new account:</h3><div><label for="newmail">E-mail:</label><input type="email" id="newmail" name="new_mail"></div><div><label for="newpass">Password:</label><input type="password" id="newpass" name="new_password"></div><div><label for="repass">Re-enter Password:</label><input type="password" id="repass" name="reentered_password"></div><div id="create" class="button"><button type="submit">Create Account</button></div></div>';
};
if (goog.DEBUG) {
  views.infobook.login.soyTemplateName = 'views.infobook.login';
}


views.infobook.profile = function(opt_data, opt_ignored) {
  return '<div id="profileview"><br>Email: ' + soy.$$escapeHtml(opt_data.email) + '<br>Name: ' + soy.$$escapeHtml(opt_data.name) + '<br>Address: ' + soy.$$escapeHtml(opt_data.address) + '<br>Phone No.: ' + soy.$$escapeHtml(opt_data.phone) + '</div><div id="update" class="button"><button type="submit">Update Account</button></div>';
};
if (goog.DEBUG) {
  views.infobook.profile.soyTemplateName = 'views.infobook.profile';
}


views.infobook.update = function(opt_data, opt_ignored) {
  return '<div id="updateview"><br><div><label for="updatemail">E-mail:</label><input type="email" id="updatemail" value="' + soy.$$escapeHtml(opt_data.email) + '"></div><br><div><label for="updatename">Name:</label><input id="updatename" value="' + soy.$$escapeHtml(opt_data.name) + '"></div><br><div><label for="updateaddress">Address:</label><input id="updateaddress" value="' + soy.$$escapeHtml(opt_data.address) + '"></div><br><div><label for="updatephone">Phone:</label><input type="email" id="updatephone" value="' + soy.$$escapeHtml(opt_data.phone) + '"></div></div><div id="update" class="button"><button type="submit">Update Account</button></div>';
};
if (goog.DEBUG) {
  views.infobook.update.soyTemplateName = 'views.infobook.update';
}
