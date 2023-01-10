package main

import (
	category3 "Selling/app/internal/controllers/http/category"
	client3 "Selling/app/internal/controllers/http/client"
	history3 "Selling/app/internal/controllers/http/history"
	product3 "Selling/app/internal/controllers/http/product"
	seller3 "Selling/app/internal/controllers/http/seller"
	subcategory3 "Selling/app/internal/controllers/http/subcategory"
	category2 "Selling/app/internal/service/category"
	client2 "Selling/app/internal/service/client"
	history2 "Selling/app/internal/service/history"
	product2 "Selling/app/internal/service/product"
	seller2 "Selling/app/internal/service/seller"
	subcategory2 "Selling/app/internal/service/subcategory"
	"Selling/app/internal/storage/category"
	"Selling/app/internal/storage/history"
	"Selling/app/internal/storage/product"
	"Selling/app/internal/storage/seller"
	"Selling/app/internal/storage/subcategory"
	"Selling/app/pkg/auth"
	"Selling/app/pkg/db"
	"Selling/app/pkg/digiseller"
	"Selling/app/pkg/server"
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"log"
)

func InitCfg() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
func main() {
	router := httprouter.New()
	err := InitCfg()
	if err != nil {
		log.Println(err)
		return
	}

	cfg, err := db.InitPGConfig()
	if err != nil {
		return
	}
	client, err := db.GetPGClient(context.Background(), cfg)
	if err != nil {
		log.Println(err)
		return
	}
	prodStore := product.NewProductStorage(client)
	sellerStore := seller.NewSellerStorage(client)
	subcatStore := subcategory.NewSubcategoryStorage(client)
	catStore := category.NewCategoryStorage(client)
	historyStore := history.NewHistoryStorage(client)

	digiClient := digiseller.NewDigiClient(catStore, subcatStore, sellerStore, prodStore)
	tm := auth.NewTokenManager(client)

	prodService := product2.NewProductService(prodStore)
	sellerService := seller2.NewSellerService(tm, sellerStore)
	catService := category2.NewCategoryService(catStore)
	subcatService := subcategory2.NewSubcategoryService(subcatStore)
	clientService := client2.NewClientService(digiClient, prodStore)
	historyService := history2.NewHistoryService(historyStore)
	ware := auth.NewMiddleWare(tm)

	historyHandler := history3.NewHistoryHandler(ware, historyService)
	historyHandler.Register(router)
	sellerHandler := seller3.NewSellerHandler(ware, sellerService)
	sellerHandler.Register(router)
	prodHandler := product3.NewProductHandler(ware, prodService)
	prodHandler.Register(router)
	catHandler := category3.NewCategoryHandler(ware, catService, subcatService, prodService)
	catHandler.Register(router)
	subcatHandler := subcategory3.NewSubcategoryHandler(ware, subcatService, prodService)
	subcatHandler.Register(router)
	clientHandler := client3.NewClientHandlers(clientService)
	clientHandler.Register(router)

	srv := server.NewServer()
	err = srv.Run(router)
	if err != nil {
		log.Fatal(err)
		return
	}
}
