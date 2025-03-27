module: "masse.example"
language: {
	version: "v0.13.0"
}
source: {
	kind: "self"
}
deps: {
	"github.com/marxarelli/masse@v1": {
		v: "v1.5.0"
		default: true
	}
	"github.com/marxarelli/masse-go@v1": {
		v: "v1.0.3"
	}
}
