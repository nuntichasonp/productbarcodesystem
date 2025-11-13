package main // package เดียวกับ main.go เพื่อให้สามารถเรียกใช้ต functions handler ได้

import ( // import package ที่จำเป็น
	"net/http" // ใช้สำหรับสถานะรหัส HTTP
	"github.com/gin-gonic/gin"// ใช้ framework gin สำหรับสร้าง web server และจัดการ routing
)
// GET /products

// GetProducts godoc
// @Summary Get all products
// @Description Get list of all products
// @Tags products
// @Produce json
// @Success 200 {array} Product
// @Router /products [get]
func GetProducts(c *gin.Context) {// ฟังก์ชันสำหรับดึงข้อมูลสินค้าทั้งหมด
	rows, _ := db.Query("SELECT id, product_code, barcode FROM products ORDER BY id DESC")// ดึงข้อมูลสินค้าทั้งหมดจากฐานข้อมูล เรียงลำดับจาก id มากไปน้อย
	defer rows.Close()// ปิดการเชื่อมต่อเมื่อฟังก์ชันนี้ทำงานเสร็จ

	var products []Product// สร้าง slice สำหรับเก็บข้อมูลสินค้า
	for rows.Next() {// วนลูปอ่านข้อมูลแต่ละแถว
		var p Product// สร้างตัวแปร p สำหรับเก็บข้อมูลสินค้าแต่ละแถว
		rows.Scan(&p.ID, &p.ProductCode, &p.Barcode)// อ่านข้อมูลจากแถวปัจจุบันมาเก็บในตัวแปร p
		products = append(products, p)// เพิ่มข้อมูลสินค้า p ลงใน slice products
	}

	c.JSON(http.StatusOK, products)// ส่งข้อมูลสินค้าในรูปแบบ JSON พร้อมสถานะรหัส 200 OK
}



// GET /products/:id

// GetProduct godoc
// @Summary Get a product by ID
// @Description Get a product by its ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} Product
// @Failure 404 {object} map[string]string
// @Router /products/{id} [get]
func GetProduct(c *gin.Context) {// ฟังก์ชันสำหรับดึงข้อมูลสินค้าตาม id
	id := c.Param("id")// ดึงค่า id จาก path parameter
	var p Product// สร้างตัวแปร p สำหรับเก็บข้อมูลสินค้า
	err := db.QueryRow("SELECT id, product_code, barcode FROM products WHERE id = ?", id).Scan(&p.ID, &p.ProductCode, &p.Barcode) // ดึงข้อมูลสินค้าตาม id จากฐานข้อมูล
	if err != nil {// ตรวจสอบ error ในการดึงข้อมูล
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})// ถ้าไม่พบสินค้าตาม id ให้ส่งสถานะรหัส 404 Not Found พร้อมข้อความแสดงข้อผิดพลาด
		return// ออกจากฟังก์ชัน
	}
	c.JSON(http.StatusOK, p)// ส่งข้อมูลสินค้าในรูปแบบ JSON พร้อมสถานะรหัส 200 OK
}

// GET /products/search?product_code=XXXX-XXXX-XXXX-XXXX

// SearchProductByCode godoc
// @Summary Search a product by code
// @Description Search a product by product_code
// @Tags products
// @Produce json
// @Param product_code query string true "Product Code"
// @Success 200 {object} Product
// @Failure 404 {object} map[string]string
// @Router /products/search [get]
func SearchProductByCode(c *gin.Context) { // ฟังก์ชันสำหรับค้นหาสินค้าตาม product_code
	code := c.Query("product_code") // ดึงค่า product_code จาก query parameter
	var p Product // สร้างตัวแปร p สำหรับเก็บข้อมูลสินค้า
	err := db.QueryRow("SELECT id, product_code, barcode FROM products WHERE product_code = ?", code).Scan(&p.ID, &p.ProductCode, &p.Barcode) // ดึงข้อมูลสินค้าตาม product_code จากฐานข้อมูล
	if err != nil { // ตรวจสอบ error ในการดึงข้อมูล
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"}) // ถ้าไม่พบสินค้าตาม product_code ให้ส่งสถานะรหัส 404 Not Found พร้อมข้อความแสดงข้อผิดพลาด
		return // ออกจากฟังก์ชัน
	}
	c.JSON(http.StatusOK, p) // ส่งข้อมูลสินค้าในรูปแบบ JSON พร้อมสถานะรหัส 200 OK
}

// POST /products

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with JSON payload
// @Tags products
// @Accept json
// @Produce json
// @Param product body Product true "Product to create"
// @Success 201 {object} Product
// @Failure 400 {object} map[string]string
// @Router /products [post]
func CreateProduct(c *gin.Context) {
	var p Product // สร้างตัวแปร p สำหรับเก็บข้อมูลสินค้า
	if err := c.ShouldBindJSON(&p); err != nil { // ผูกข้อมูล JSON ที่ส่งมาในคำขอเข้ากับตัวแปร p
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"}) // ส่งสถานะรหัส 400 Bad Request พร้อมข้อความแสดงข้อผิดพลาด
		return // ออกจากฟังก์ชันถ้า JSON ที่ส่งมาไม่ถูกต้อง
	}

	if !ValidateProductcode(p.ProductCode) { // ตรวจสอบรูปแบบของ product_code
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product_code format"}) // ส่งสถานะรหัส 400 Bad Request พร้อมข้อความแสดงข้อผิดพลาด
		return // ออกจากฟังก์ชันถ้า product_code ไม่ถูกต้อง
	}

	// Insert product
	result, err := db.Exec(
		"INSERT INTO products (product_code, barcode, create_at, update_at) VALUES (?, ?, ?, ?)",
		p.ProductCode, p.ProductCode, p.CreatedAt, p.UpdateAt,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Duplicate productcode"})
		return
	}

	// ดึง ID ที่เพิ่ง insert 
	lastID, _ := result.LastInsertId()

	// ดึงข้อมูลจาก DB กลับมา
	err = db.QueryRow(
		"SELECT id, product_code, barcode, create_at, update_at FROM products WHERE id = ?",
		lastID,
	).Scan(&p.ID, &p.ProductCode, &p.Barcode)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot read created product"})
		return
	}

	c.JSON(http.StatusCreated, p) // ส่งข้อมูลสินค้าในรูปแบบ JSON พร้อมสถานะรหัส 201 Created
}

// PUT /products/:id

// UpdateProduct godoc
// @Summary Update a product
// @Description Update a product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body Product true "Updated product data"
// @Success 200 {object} Product
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id} [put]
func UpdateProduct(c *gin.Context) {
	id := c.Param("id") // ดึงค่า id จาก path parameter
	var p Product // สร้างตัวแปร p สำหรับเก็บข้อมูลสินค้า
	if err := c.ShouldBindJSON(&p); err != nil { // ผูกข้อมูล JSON ที่ส่งมาในคำขอเข้ากับตัวแปร p
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"}) // ส่งสถานะรหัส 400 Bad Request พร้อมข้อความแสดงข้อผิดพลาด
		return // ออกจากฟังก์ชันถ้า JSON ที่ส่งมาไม่ถูกต้อง
	}

	if !ValidateProductcode(p.ProductCode) { // ตรวจสอบรูปแบบของ product_code
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ProductCode format"}) // ส่งสถานะรหัส 400 Bad Request พร้อมข้อความแสดงข้อผิดพลาด
		return // ออกจากฟังก์ชันถ้า product_code ไม่ถูกต้อง
	}

	_, err := db.Exec("UPDATE products SET product_code = ?, update_at = CURRENT_TIMESTAMP WHERE id = ?", p.ProductCode, id) // อัปเดตข้อมูลสินค้าตาม id ในฐานข้อมูล
	if err != nil { // ตรวจสอบ error ในการอัปเดตข้อมูล
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"}) // ถ้าเกิดข้อผิดพลาดในการอัปเดต ให้ส่งสถานะรหัส 500 Internal Server Error พร้อมข้อความแสดงข้อผิดพลาด
		return // ออกจากฟังก์ชัน
	}

	c.JSON(http.StatusOK, p) // ส่งข้อมูลสินค้าในรูปแบบ JSON พร้อมสถานะรหัส 200 OK
}

// DELETE /products/:id

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by its ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 204 {string} string "No Content"
// @Failure 500 {object} map[string]string
// @Router /products/{id} [delete]
func DeleteProduct(c *gin.Context) { 
	id := c.Param("id") // ดึงค่า id จาก path parameter
	_, err := db.Exec("DELETE FROM products WHERE id = ?", id) // ลบสินค้าตาม id จากฐานข้อมูล
	if err != nil { // ตรวจสอบ error ในการลบข้อมูล
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"}) // ถ้าเกิดข้อผิดพลาดในการลบ ให้ส่งสถานะรหัส 500 Internal Server Error พร้อมข้อความแสดงข้อผิดพลาด
		return // ออกจากฟังก์ชัน
	}
	c.Status(http.StatusNoContent) // ส่งสถานะรหัส 204 No Content เพื่อบอกว่าการลบสำเร็จแล้ว
}
