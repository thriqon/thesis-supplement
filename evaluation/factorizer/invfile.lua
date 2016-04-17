
inv.task('build')
  .using("frolvlad/alpine-gcc:latest")
    .run('/bin/sh', '-c', 'mkdir -p dist && gcc -o dist/factorizer factorizer.c')

inv.task('package')
  .wrap('dist')
    .inImage('alpine:3.3')
    .at("/")
    .withConfig({cmd = {"/factorizer"}})
    .as(VAR.TAG)
    
