package main

import (
	"log"
	"os"
	"server/internal/config"
	"server/internal/handler"
	"server/internal/repository"
	"server/internal/service"
	"server/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main(){
	err := godotenv.Load()
	if err != nil{
		log.Fatal("Error loading .env")
	}
	config.ConnectDB()

	repoInv:=repository.NewInventoryRepository(config.DB)
	serviceInv:=service.NewInventoryService(repoInv)
	handlerInv:=handler.NewInventoryHandler(serviceInv)

	repoProd:=repository.NewProductRepository(config.DB)
	serviceProd:=service.NewProductService(repoProd)
	handlerProd:=handler.NewProductHandler(serviceProd)


	r:=gin.Default()

	r.GET("/ping",func(c *gin.Context){
		c.JSON(200,gin.H{
			"message":"pong",
		})
	})

	routes.SetupRoutes(r,handlerInv,handlerProd)
	port := os.Getenv("PORT")
	r.Run(":"+port)
}