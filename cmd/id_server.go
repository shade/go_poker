package main
import (
	"net/http"
	"go_poker/internal/identity"
)

func main {
	gen := identity.NewIdentityGenerator("id_secret");

	routes := path.Prefix("/api/v1/identity")
	routes = gen.CreateRoutes(routes)

	http.ListenAndServe("id_port", routes)
}