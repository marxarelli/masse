module: "masse.example"
language: {
	version: "v0.13.0"
}
source: {
	kind: "self"
}
deps: {
	"github.com/marxarelli/masse-go@v2": {
		v: "v2.0.0"
	}
	"github.com/marxarelli/masse@v1": {
		v:       "v1.9.0"
		default: true
	}
}
