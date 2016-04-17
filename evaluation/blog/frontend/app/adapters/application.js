
import DS from 'ember-data';
import ENV from 'blog/config/environment';

const {APP} = ENV;

const {RESTAdapter} = DS;

export default RESTAdapter.extend({
  host: APP.host,
  namespace: APP.prefix
});
