
inv.task('package')
  .wrap('native').at('/').inImage('busybox:latest')
    .withConfig({Cmd = {'/bin/sh', '/helloworld.sh'}})
    .as(VAR.TAG)

