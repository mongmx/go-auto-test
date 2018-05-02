get-dep:
	go get github.com/stretchr/testify
	go get github.com/smartystreets/goconvey

autotest:
	goconvey
