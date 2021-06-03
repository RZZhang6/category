package main

import (
	"github.com/RZZhang6/category/common"
	"github.com/RZZhang6/category/domain/repository"
	service2 "github.com/RZZhang6/category/domain/service"
	"github.com/RZZhang6/category/handler"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"

	category "github.com/RZZhang6/category/proto/category"
)

func main() {
	// 配置中心
	consulConfig, err := common.GetConsulConfig("192.168.153.135", 8500,
		"/micro/config")
	if err != nil {
		log.Error(err)
	}
	// 注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			 "192.168.153.135:8500",
		}
	})


	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.category"),
		micro.Version("latest"),
		// 这里设置地址和需要暴露的端口
		micro.Address("127.0.0.1:8082"),
		//添加consul 作为注册中心
		micro.Registry(consulRegistry),
	)

	//获取Mysql配置，路径中不带前缀
	mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")

	// 连接数据库
	db,err := gorm.Open("mysql",mysqlInfo.User+":"+mysqlInfo.Pwd+"@/"+
		mysqlInfo.Database+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Error(err)
	}
	defer db.Close()

	// 禁止复表
	db.SingularTable(true)


	// Initialise service
	service.Init()

	categoryDataService := service2.NewCategoryDataService(repository.NewCategoryRepository(db))

	// Register Handler
	err = category.RegisterCategoryHandler(service.Server(), &handler.Category{CategoryDataService: categoryDataService})
	if err != nil {
		log.Error(err)
	}

	// Register Struct as Subscriber
	//micro.RegisterSubscriber("go.micro.service.category", service.Server(), new(subscriber.Category))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
