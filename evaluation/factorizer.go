package main

import "time"

func init() {
	registerSuite("Factorizer", "Dockerfile a", factorizerDockerfileA)
	registerSuite("Factorizer", "Dockerfile b", factorizerDockerfileB)
	registerSuite("Factorizer", "Squash and Load", factorizerSquashAndLoad)
	registerSuite("Factorizer", "Layer Donning", factorizerLayerDonning)
}

func factorizerDockerfileVariant(variant string) result {
	ensureNotPresent("alpine", "test3/fac", "test3/fac2")
	defer ensureNotPresent("alpine", "test3/fac", "test3/fac2")

	buildVariant := func(variant string) func() error {
		return commandRunner("docker", "build", "-q", "-t", "test3/fac", "-f", "factorizer/Dockerfile."+variant, "factorizer")
	}

	writeFile("factorizer/factorizer.c", factorizerC)
	first := measure(buildVariant(variant), nil)

	writeFile("factorizer/factorizer.c", factorizerC2)
	second := measure(buildVariant(variant), nil)

	return result{first, second, imageSize("test3/fac")}
}

func factorizerDockerfileA() result {
	return factorizerDockerfileVariant("a")
}

func factorizerDockerfileB() result {
	return factorizerDockerfileVariant("b")
}

func factorizerSquashAndLoad() result {
	ensureNotPresent("alpine", "test3/fac", "test3/fac2")
	defer ensureNotPresent("alpine", "test3/fac", "test3/fac2")

	squashAndLoad := func(tag string) {
		defer ensureNotPresent("test3/fac_t")
		mustRun("docker", "build", "-f", "factorizer/Dockerfile.squash", "-t", "test3/fac_t", "factorizer")
		ids := mustOutputString("docker", "create", "test3/fac_t", "/bin/sh")
		defer mustRun("docker", "rm", ids)
		mustRun("/bin/sh", "-c", "docker export "+ids+" | "+
			"docker import -c 'CMD /factorizer' - "+tag)
	}

	writeFile("factorizer/factorizer.c", factorizerC)
	first := measure(func() error {
		squashAndLoad("test3/fac")
		return nil
	}, nil)

	writeFile("factorizer/factorizer.c", factorizerC2)
	second := measure(func() error {
		squashAndLoad("test3/fac2")
		return nil
	}, nil)

	return result{first, second, imageSize("test3/fac")}
}

func factorizerLayerDonning() result {
	ensureNotPresent("frolvlad/alpine-gcc", "alpine:3.3", "test3/fac", "test3/fac2")
	defer ensureNotPresent("frolvlad/alpine-gcc", "alpine:3.3", "test3/fac", "test3/fac2")

	writeFile("factorizer/factorizer.c", factorizerC)
	first := measure(commandRunner("involucro", "-w", "factorizer", "-f", "factorizer/invfile.lua", "--set", "TAG=test3/fac", "build", "package"), nil)

	time.Sleep(20 * time.Second)

	writeFile("factorizer/factorizer.c", factorizerC2)
	second := measure(commandRunner("involucro", "-w", "factorizer", "-f", "factorizer/invfile.lua", "--set", "TAG=test3/fac2", "build", "package"), nil)

	return result{first, second, imageSize("test3/fac")}
}

var factorizerC = `
#include <stdio.h>

int main() {
  int target = 0;

  if (scanf("%d", &target) < 1 || target <= 1) {
    fprintf(stderr, "Invalid number, should greater 1\n");
    return 1;
  }

  printf("%d:", target);

  int divider = 2;
  while(target > 1) {

    while (target % divider == 0) {
      printf(" %d", divider);
      target /= divider;
    }

    divider++;
  }
  printf("\n");
  return 0;
}
`

var factorizerC2 = `
#include <stdio.h>

#define INITIAL_DIVIDER 2

int main() {
  int target = 0;

  if (scanf("%d", &target) < 1 || target <= 1) {
    fprintf(stderr, "Invalid number, should greater 1\n");
    return 1;
  }

  printf("%d:", target);

  int divider = INITIAL_DIVIDER;
  while(target > 1) {

    while (target % divider == 0) {
      printf(" %d", divider);
      target /= divider;
    }

    divider++;
  }
  printf("\n");
  return 0;
}
`
