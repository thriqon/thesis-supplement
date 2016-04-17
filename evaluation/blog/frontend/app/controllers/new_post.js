import Ember from 'ember';

const {Controller} = Ember;

export default Controller.extend({
  actions: {
    save(title, content) {
      this.store.createRecord('post', {
        title, content
      }).save().then(() => this.transitionToRoute('posts'));
    }
  }
});
