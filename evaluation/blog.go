package main

func init() {
	registerSuite("Blog", "Dockerfile JS", blogDockerfileJs)
	registerSuite("Blog", "Dockerfile Go", blogDockerfileGo)
	registerSuite("Blog", "Squash and Load", blogSquashAndLoad)
	registerSuite("Blog", "Layer Donning", blogLayerDonning)
}

func blogDockerfileGo() result {
	ensureNotPresent("alpine", "test4/blog", "test4/blog2")
	defer ensureNotPresent("alpine", "test4/blog", "test4/blog2")

	writeFile("blog/backend/flags.go", flagsGo)
	first := measure(commandRunner("docker", "build", "-q", "-t", "test4/blog", "-f", "blog/Dockerfile", "blog"), nil)

	writeFile("blog/backend/flags.go", flagsGo2)
	second := measure(commandRunner("docker", "build", "-q", "-t", "test4/blog", "-f", "blog/Dockerfile", "blog"), nil)

	return result{first, second, imageSize("test4/blog")}
}

func blogDockerfileJs() result {
	ensureNotPresent("alpine", "test4/blog", "test4/blog2")
	defer ensureNotPresent("alpine", "test4/blog", "test4/blog2")

	writeFile("blog/frontend/app/router.js", routerJs)
	first := measure(commandRunner("docker", "build", "-q", "-t", "test4/blog", "-f", "blog/Dockerfile", "blog"), nil)

	writeFile("blog/frontend/app/router.js", routerJs2)
	second := measure(commandRunner("docker", "build", "-q", "-t", "test4/blog", "-f", "blog/Dockerfile", "blog"), nil)

	return result{first, second, imageSize("test4/blog")}
}

func blogSquashAndLoad() result {
	ensureNotPresent("alpine", "test4/blog", "test4/blog2", "test4/blog_t")
	defer ensureNotPresent("alpine", "test4/blog", "test4/blog2")

	defer ensureNotPresent("test4/blog_t")
	squashAndLoad := func(tag string) {
		mustRun("docker", "build", "-f", "blog/Dockerfile.squash", "-t", "test4/blog_t", "blog")
		ids := mustOutputString("docker", "create", "test4/blog_t", "/bin/sh")
		defer mustRun("docker", "rm", ids)
		mustRun("/bin/sh", "-c", "docker export "+ids+" | "+
			"docker import -c 'ENTRYPOINT /blog/blog' - "+tag)
	}

	writeFile("blog/backend/flags.go", flagsGo)
	first := measure(func() error {
		squashAndLoad("test4/blog")
		return nil
	}, nil)

	writeFile("blog/backend/flags.go", flagsGo2)
	second := measure(func() error {
		squashAndLoad("test4/blog2")
		return nil
	}, nil)

	return result{first, second, imageSize("test4/blog")}
}

func blogLayerDonning() result {
	ensureNotPresent("golang:1.6", "danlynn/ember-cli:2.3.0", "test4/blog", "test4/blog2")
	defer ensureNotPresent("golang:1.6", "danlynn/ember-cli:2.3.0", "test4/blog", "test4/blog2")

	cleanInstance := func() {
		mustRun("involucro", "-w", "blog", "-f", "blog/invfile.lua", "clean")
	}

	writeFile("blog/backend/flags.go", flagsGo)
	first := measure(commandRunner("involucro", "-w", "blog", "-f", "blog/invfile.lua", "--set", "TAG=test4/blog", "build", "package"), cleanInstance)

	writeFile("blog/backend/flags.go", flagsGo2)
	second := measure(commandRunner("involucro", "-w", "blog", "-f", "blog/invfile.lua", "--set", "TAG=test4/blog2", "build", "package"), cleanInstance)

	return result{first, second, imageSize("test4/blog")}
}

var flagsGo = `
package main

import "flag"

var (
	httpSpec = flag.String("http", ":8040", "HTTP address and port")
	dbFile   = flag.String("file", "db.json", "Database file")
)
`

var flagsGo2 = `
package main

import "flag"

var (
	httpSpec = flag.String("http", ":8080", "HTTP address and port")
	dbFile   = flag.String("file", "db.json", "Database file")
)
`

var routerJs = `
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
`

var routerJs2 = `
import Ember from 'ember';
import config from './config/environment';

const Router = Ember.Router.extend({
  location: config.locationType
});

Router.map(function() {
  this.route('posts', {path: '/'});
  this.route('new_post', {path: '/new'});
});

export default Router;
`
