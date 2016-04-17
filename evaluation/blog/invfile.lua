
local package = "github.com/thriqon/blog"

inv.task('build:server')
  .using('golang:1.6')
    .withConfig({
      env = {"CGO_ENABLED=0"},
      workingdir = "/go/src/" .. package,
    })
    .withHostConfig({
      binds = {
        "./backend:/go/src/" .. package,
        "./dist:/dist"
      }
    })
    .run('go', 'build', '-o', '/dist/blog', './.')

inv.task('build:frontend')
  .using('thriqon/alpine-ember-cli:latest')
    .withConfig({
      entrypoint = {"/bin/sh", "-c"}
    })
    .withHostConfig({
      binds = {
        './frontend:/source',
        './dist:/dist'
      }
    })
    .run('npm install && bower install --allow-root')
    .run('ember build -prod --output-path=/dist/frontend')

inv.task('build')
  .using('thriqon/alpine-ember-cli:latest')
    .withConfig({entrypoint = {"/bin/sh", "-c"}})
    .run('mkdir -p dist')
  .runTask('build:server')
  .runTask('build:frontend')

inv.task('package')
  .wrap('dist')
    .at('/srv')
    .withConfig({
      cmd = {"/srv/blog"}
    })
    .as(VAR.TAG)

inv.task('clean')
  .using('alpine:latest')
  .withHostConfig({
    binds = {
      './frontend:/source',
    }
  })
  .run('rm', '-rf', 'node_modules/', 'bower_components/')
