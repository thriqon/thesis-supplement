import Ember from 'ember';
import config from './config/environment';

const Router = Ember.Router.extend({
  location: config.locationType
});

Router.map(function() {
  this.route('new_post', {path: '/new'});
  this.route('posts', {path: '/'});
});

export default Router;
