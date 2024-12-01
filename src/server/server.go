package server

import (
	"context"
	"gtmx/src/database/repository"
	"gtmx/src/handler"
	customMiddleware "gtmx/src/middleware"
	"gtmx/src/server/routes"
	"gtmx/src/service"
	"gtmx/src/service/auth"
	"os"
	"os/signal"
	"time"

	"github.com/antonlindstrom/pgstore"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type Server struct {
	app         *echo.Echo
	db          *gorm.DB
	store       *pgstore.PGStore
	ServerReady chan bool
}

func (s *Server) Start(address string) error {
	s.setRoutes()
	go s.app.Start(address)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	if s.ServerReady != nil {
		s.ServerReady <- true
	}

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.app.Shutdown(ctx)
}

func (s *Server) setRoutes() {
	s.app.Use(middleware.Logger())
	s.app.Use(session.Middleware(s.store))
	s.app.Use(middleware.Recover())

	errorHandler := handler.ErrorHandler{}
	s.app.HTTPErrorHandler = errorHandler.HandleError

	s.app.Static("/static", "static")

	s.app.GET("/health", func(c echo.Context) error {
		c.Response().Status = 200
		return nil
	})

	userRepository := repository.NewUserRepository(s.db)
	sectionRepository := repository.NewSectionRepository(s.db)
	categoryRepository := repository.NewCategoryRepository(s.db)
	menuRepository := repository.NewMenuRepository(s.db)
	productRepository := repository.NewProductRepository(s.db)
	variantRepository := repository.NewVariantRepository(s.db)

	variantService := service.NewVariantService(&variantRepository)
	productService := service.NewProductService(&productRepository)
	categoryService := service.NewCategoryService(&categoryRepository)
	sectionService := service.NewSectionService(&sectionRepository)
	menuService := service.NewMenuService(&menuRepository)
	userService := service.NewUserService(&userRepository)

	authService := auth.NewAuthService(&userService)

	authHandler := handler.NewAuthHandler(&authService)
	s.app.GET("/signup", authHandler.HandleShowSignUp).Name = routes.SHOW_SIGN_UP
	s.app.POST("/signup", authHandler.HandleSignUp).Name = routes.SIGN_UP
	s.app.GET("/login", authHandler.HandleShowLogin).Name = routes.SHOW_LOGIN
	s.app.POST("/login", authHandler.HandleLogin).Name = routes.LOGIN
	s.app.GET("/logout", authHandler.HandleLogout).Name = routes.LOGOUT

	authenticatedRoutes := s.app.Group("/admin")
	authenticatedRoutes.Use(customMiddleware.Authenticated)
	s.app.Use(customMiddleware.ResponseHeaders)

	menuHandler := handler.NewMenuHandler(&menuService)
	authenticatedRoutes.GET("/menus", menuHandler.HandleIndex).Name = routes.INDEX_MENU
	authenticatedRoutes.GET("/menus/new", menuHandler.HandleNew).Name = routes.NEW_MENU
	authenticatedRoutes.POST("/menus", menuHandler.HandleCreate).Name = routes.CREATE_MENU
	authenticatedRoutes.GET("/menus/:id", menuHandler.HandleShow).Name = routes.SHOW_MENU

	productHandler := handler.NewProductHandler(&productService, &categoryService)
	authenticatedRoutes.GET("/products", productHandler.HandleIndex).Name = routes.INDEX_PRODUCT
	authenticatedRoutes.POST("/products", productHandler.HandleCreate).Name = routes.CREATE_PRODUCT
	authenticatedRoutes.GET("/products/:id", productHandler.HandleShow).Name = routes.SHOW_PRODUCT
	authenticatedRoutes.GET("/products/new", productHandler.HandleNew).Name = routes.NEW_PRODUCT
	authenticatedRoutes.GET("/products/:id/edit", productHandler.HandleEdit).Name = routes.EDIT_PRODUCT
	authenticatedRoutes.POST("/products/:id", productHandler.HandleUpdate).Name = routes.UPDATE_PRODUCT

	sectionHandler := handler.NewSectionHandler(&sectionService)
	authenticatedRoutes.GET("/sections", sectionHandler.HandleIndex).Name = routes.INDEX_SECTION
	authenticatedRoutes.POST("/sections", sectionHandler.HandleCreate).Name = routes.CREATE_SECTION
	authenticatedRoutes.GET("/sections/:id", sectionHandler.HandleShow).Name = routes.SHOW_SECTION
	authenticatedRoutes.GET("/sections/new", sectionHandler.HandleNew).Name = routes.NEW_SECTION
	authenticatedRoutes.GET("/sections/:id/edit", sectionHandler.HandleEdit).Name = routes.EDIT_SECTION
	authenticatedRoutes.POST("/sections/:id", sectionHandler.HandleUpdate).Name = routes.UPDATE_SECTION

	categoryHandler := handler.NewCategoryHandler(&sectionService, &categoryService)
	authenticatedRoutes.GET("/categories", categoryHandler.HandleIndex).Name = routes.INDEX_CATEGORY
	authenticatedRoutes.POST("/categories", categoryHandler.HandleCreate).Name = routes.CREATE_CATEGORY
	authenticatedRoutes.GET("/categories/:id", categoryHandler.HandleShow).Name = routes.SHOW_CATEGORY
	authenticatedRoutes.GET("/categories/new", categoryHandler.HandleNew).Name = routes.NEW_CATEGORY
	authenticatedRoutes.GET("/categories/:id/edit", categoryHandler.HandleEdit).Name = routes.EDIT_CATEGORY
	authenticatedRoutes.POST("/categories/:id", categoryHandler.HandleUpdate).Name = routes.UPDATE_CATEGORY

	variantHandler := handler.NewVariantHandler(&variantService)
	authenticatedRoutes.GET("/variants", variantHandler.HandleIndex).Name = routes.INDEX_VARIANT
	authenticatedRoutes.POST("/variants", variantHandler.HandleCreate).Name = routes.CREATE_VARIANT
	authenticatedRoutes.GET("/variants/:id", variantHandler.HandleShow).Name = routes.SHOW_VARIANT
	authenticatedRoutes.GET("/variants/new", variantHandler.HandleNew).Name = routes.NEW_VARIANT
	authenticatedRoutes.GET("/variants/:id/edit", variantHandler.HandleEdit).Name = routes.EDIT_VARIANT
	authenticatedRoutes.POST("/variants/:id", variantHandler.HandleUpdate).Name = routes.UPDATE_VARIANT

	routes.SetRoutesMap(s.app.Routes())
}

func New(db *gorm.DB, store *pgstore.PGStore, serveReady chan bool) Server {
	return Server{
		app:         echo.New(),
		store:       store,
		db:          db,
		ServerReady: serveReady,
	}
}
