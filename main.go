package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"lenslocked.com/controllers"
	"lenslocked.com/middleware"
	"lenslocked.com/models"
	"lenslocked.com/rand"
)

func main() {
	boolPtr := flag.Bool("prod", false, "Provide this flag in production. This ensures that a .config file is provided before the application starts.")
	r := mux.NewRouter()
	static := controllers.NewStatic()

	cfg := LoadConfig(*boolPtr)
	dbCfg := DefaultPosgresConfig()
	services, err := models.NewServices(
		models.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInfo()),
		models.WithLogMode(!cfg.IsProd()),
		models.WithUser(cfg.Pepper, cfg.HMACKey),
		models.WithGallery(),
		models.WithImage(),
	)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	usersController := controllers.NewUser(services.User)
	galleriesC := controllers.NewGallery(services.Gallery, services.Image, r)

	//CSRF
	csrfSeed, err := rand.Bytes(32)
	must(err)
	// TODO: update this to be a config variable
	isProd := false
	csrfMw := csrf.Protect(csrfSeed, csrf.Secure(isProd))

	userMw := middleware.User{
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{
		User: userMw,
	}

	r.HandleFunc("/", static.HomeView.ServeHTTP).Methods("GET")
	r.HandleFunc("/contact", static.ContactView.ServeHTTP).Methods("GET")
	r.HandleFunc("/sign-up", usersController.New).Methods("GET")
	r.HandleFunc("/sign-up", usersController.Create).Methods("POST")
	r.HandleFunc("/login", usersController.LoginView.ServeHTTP).Methods("GET")
	r.HandleFunc("/login", usersController.SignIn).Methods("POST")

	// Image routes
	imageHandler := http.FileServer(http.Dir("./images/"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imageHandler))

	// Assets
	assetHandler := http.FileServer(http.Dir("./assets/"))
	assetHandler = http.StripPrefix("/assets/", assetHandler)
	r.PathPrefix("/assets/").Handler(assetHandler)

	// Gallery routes
	r.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesC.Index)).Methods("GET")
	r.HandleFunc("/galleries/new", requireUserMw.ApplyFn(galleriesC.New)).Methods("GET")
	r.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesC.Create)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMw.ApplyFn(galleriesC.Edit)).Methods("GET").Name("edit_gallery")
	r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMw.ApplyFn(galleriesC.Update)).Methods("POST").Name("update_gallery")
	r.HandleFunc("/galleries/{id:[0-9]+}/images", requireUserMw.ApplyFn(galleriesC.ImageUpload)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete", requireUserMw.ApplyFn(galleriesC.ImageDelete)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMw.ApplyFn(galleriesC.Delete)).Methods("POST").Name("delete_gallery")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").Name("show_gallery")

	r.NotFoundHandler = http.HandlerFunc(static.NotFoundView.ServeHTTP)

	p := fmt.Sprintf(":%d", cfg.Port)

	fmt.Println("Running on port 3000")
	log.Fatal(http.ListenAndServe(p, csrfMw(userMw.Apply(r))))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
