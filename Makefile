binary: *.go */*.go go.*
	# the executable
	go build -o $@ -ldflags "-s -w"
	file $@

module.tar.gz: pca9685
	# the bundled module
	rm -f $@
	tar czf $@ $^