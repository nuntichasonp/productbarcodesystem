package main // package เดียวกับ main.go เพื่อให้สามารถเรียกใช้ตโครงสร้าง Model Product ได้

import "regexp" // ใช้สำหรับตรวจสอบรูปแบบของ Productcode

// โครงสร้างข้อมูลสินค้า
type Product struct {
	ID      int    `json:"id"` // รหัสสินค้า
	ProductCode string `json:"product_code"` // รหัสผลิตภัณฑ์
	Barcode string `json:"barcode"` // บาร์โค้ดของสินค้า
	CreatedAt string `json:"created_at"` // วันที่สร้างสินค้า
	UpdateAt string `json:"update_at"` // วันที่แก้ไขสินค้า
}

// ตรวจสอบรูปแบบ Productcode
func ValidateProductcode(code string) bool {
	re := regexp.MustCompile(`^[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{4}$`) // รูปแบบ Productcode ที่ถูกต้อง เช่น ABCD-1234-EFGH-5678
	//MustCompile สร้าง regex หนึ่งครั้งแล้วเก็บไว้ใช้ซ้ำได้เลย
	//^[A-Z0-9]{4} — ตัวอักษรใหญ่หรือตัวเลข 4 ตัว
	// - — ตามด้วยขีดกลาง
	// ทำซ้ำ 4 ชุด
	return re.MatchString(code)// ตรวจสอบว่า code ตรงกับรูปแบบที่กำหนดหรือไม่
}
