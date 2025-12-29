package main

import (
	"cool-games/config"
	"fmt"
	"os"
	"time"

	gameDelivery "cool-games/internal/game/delivery"
	gameRepo "cool-games/internal/game/repository"
	gameUcase "cool-games/internal/game/usecase"

	authDelivery "cool-games/internal/auth/delivery"
	authRepo "cool-games/internal/auth/repository"
	authUcase "cool-games/internal/auth/usecase"

	orderDelivery "cool-games/internal/order/delivery"
	orderRepo "cool-games/internal/order/repository"
	orderUcase "cool-games/internal/order/usecase"

	genreDelivery "cool-games/internal/genre/delivery"
    genreRepo "cool-games/internal/genre/repository"
    genreUcase "cool-games/internal/genre/usecase"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	db := config.ConnectDB()
    if err := godotenv.Load(); err != nil {
        fmt.Println("Warning: .env file not found, using system environment")
    }

    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        jwtSecret = "default_secret_for_dev_only"
    }
    
    fmt.Printf("LOG: JWT Secret loaded. Length: %d characters\n", len(jwtSecret))

	r := gin.Default()

	uRepo := authRepo.NewPsqlUserRepository(db)
	cRepo := authRepo.NewPsqlCustomerRepository(db)
	
	aUcase := authUcase.NewAuthUsecase(uRepo, cRepo, jwtSecret, 5*time.Second)
	custUcase := authUcase.NewCustomerUsecase(cRepo, 5*time.Second)
	
	authDelivery.NewAuthHandler(r, aUcase)
	authDelivery.NewCustomerHandler(r, custUcase, jwtSecret)

	gRepo := gameRepo.NewPsqlGameRepository(db)
	gUcase := gameUcase.NewGameUsecase(gRepo, 5*time.Second)
	gameDelivery.NewGameHandler(r, gUcase, jwtSecret)

	oRepo := orderRepo.NewPsqlOrderRepository(db)
    lRepo := orderRepo.NewPsqlLibraryRepository(db)
	oUcase := orderUcase.NewOrderUsecase(gRepo, cRepo, oRepo, lRepo, 10*time.Second) 
    orderDelivery.NewOrderHandler(r, oUcase, jwtSecret)

	genreRepo := genreRepo.NewPsqlGenreRepository(db)
	genreUcase := genreUcase.NewGenreUsecase(genreRepo, 5*time.Second)
	genreDelivery.NewGenreHandler(r, genreUcase, jwtSecret)

	r.Run(":8080")
}