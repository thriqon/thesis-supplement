package main

func init() {
	registerSuite("Hello, World", "Dockerfile", helloWorldDockerfile)
	registerSuite("Hello, World", "Squash and Load", helloWorldSquashAndLoad)
	registerSuite("Hello, World", "Layer Donning", helloWorldLayerDonning)
}

var (
	helloWorld = "#!/bin/bash\necho \"Hello, World!\""
	helloMoon  = "#!/bin/bash\necho \"Hello, Moon!\""
)

func useHelloWorldFile(contents string) {
	writeFile("helloworld/native/helloworld.sh", contents)
}

func helloWorldDockerfile() result {
	ensureNotPresent("busybox", "test1/hw1", "test1/hw2")
	defer ensureNotPresent("test1/hw1", "test1/h2")

	useHelloWorldFile(helloWorld)
	first := measure(commandRunner("docker", "build", "-q", "-t", "test1/hw1", "helloworld"), nil)

	useHelloWorldFile(helloMoon)
	second := measure(commandRunner("docker", "build", "-q", "-t", "test1/hw2", "helloworld"), nil)

	return result{first, second, imageSize("test1/hw1")}
}

func helloWorldSquashAndLoad() result {
	ensureNotPresent("busybox", "test2/hw1", "test2/hw2")
	defer ensureNotPresent("test2/hw1", "test2/h2")

	squashAndLoad := func(tag string) func() {
		return func() {
			ids := mustOutputString("docker", "create", "busybox")
			mustRun("docker", "cp", "helloworld/native/helloworld.sh", ids+":/")
			mustRun("/bin/sh", "-c", "docker export "+ids+" | "+
				"docker import -c 'CMD /bin/sh /helloworld.sh' - "+tag)
		}
	}

	useHelloWorldFile(helloWorld)
	first := measure(squashAndLoad("test2/hw1"), nil)

	useHelloWorldFile(helloMoon)
	second := measure(squashAndLoad("test2/hw2"), nil)

	return result{first, second, imageSize("test2/hw1")}
}

func helloWorldLayerDonning() result {
	ensureNotPresent("busybox", "test3/hw1", "test3/hw2")
	defer ensureNotPresent("test3/hw1", "test3/h2")

	useHelloWorldFile(helloWorld)
	first := measure(commandRunner("involucro", "-w", "helloworld", "-f", "helloworld/invfile.lua", "--set", "TAG=test3/hw1", "package"), nil)
	useHelloWorldFile(helloMoon)
	second := measure(commandRunner("involucro", "-w", "helloworld", "-f", "helloworld/invfile.lua", "--set", "TAG=test3/hw2", "package"), nil)

	return result{first, second, imageSize("test3/hw1")}
}
