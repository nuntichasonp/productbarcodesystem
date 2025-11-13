package main // จะเป็น package หลักของโปรเจค และจะรันไฟล์นี้เป็นหลัก โดย go จะเริ่มต้นการทำงานจากฟังก์ชัน main()

import ( //import package ที่จำเป็น
	"log" // ใช้สำหรับ logging ข้อความต่างๆ ไปยัง console เวลา run server
	"github.com/gin-gonic/gin" // ใช้ framework gin สำหรับสร้าง web server และจัดการ routing
	"github.com/gin-contrib/cors" // ใช้สำหรับจัดการ CORS (Cross-Origin Resource Sharing)
	"time" // ใช้สำหรับจัดการเวลา
	swaggerFiles "github.com/swaggo/files" // swaggerFiles ใช้สำหรับให้บริการไฟล์ Swagger UI
    ginSwagger "github.com/swaggo/gin-swagger" // ginSwagger ใช้สำหรับผนวก Swagger UI เข้ากับ Gin framework
	_ "go-rest-api/docs" // นำเข้า package docs ที่สร้างโดย swaggo เพื่อให้สามารถใช้งาน Swagger UI ได้
)

func main() {
	InitDB() // คือฟังก์ชันที่เราสร้างไว้ในไฟล์ database.go เพื่อเชื่อมต่อฐานข้อมูล

	r := gin.Default() //gin.Default() สร้าง router พร้อม middleware มาตรฐาน (logging, recovery) 
	
	// ใช้ตัวแปร r เป็นตัวแทนของเว็บเซิร์ฟเวอร์ทั้งหมด

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:4200"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"*"},
        AllowCredentials: true,
        MaxAge: 24 * time.Hour,
    }))

	// Routing
	r.GET("/products", GetProducts) // call function GetProducts จากไฟล์ handlers.go
	r.GET("/products/:id", GetProduct) // call function GetProduct จากไฟล์ handlers.go โดย id ตือ path parameter เช่น /products/1
	r.GET("/products/search", SearchProductByCode) // call function SearchProductByCode จากไฟล์ handlers.go โดยรับ product_code เป็น query parameter เช่น /products/search?product_code=XXXX-XXXX-XXXX-XXXX
	r.POST("/products", CreateProduct) // call function CreateProduct จากไฟล์ handlers.go
	r.PUT("/products/:id", UpdateProduct) // call function UpdateProduct จากไฟล์ handlers.go โดย id ตือ path parameter เช่น /products/1
	r.DELETE("/products/:id", DeleteProduct) // call function DeleteProduct จากไฟล์ handlers.go โดย id ตือ path parameter เช่น /products/1

	log.Println("Server started at :8080") // log ข้อความบอกว่า server เริ่มทำงานแล้วที่พอร์ต 8080
	r.Run(":8080") // เริ่มต้น web server ที่พอร์ต 8080
}
