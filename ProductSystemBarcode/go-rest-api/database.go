package main // package เดียวกับ main.go เพื่อให้สามารถเรียกใช้ตัวแปร db ได้

import ( // import package ที่จำเป็น
	"database/sql" // ใช้สำหรับเชื่อมต่อและจัดการฐานข้อมูล SQL
	"log" // ใช้สำหรับ logging ข้อความต่างๆ ไปยัง console เวลา run server
	_ "modernc.org/sqlite" // ใช้ SQLite driver สำหรับเชื่อมต่อฐานข้อมูล SQLite
)

var db *sql.DB // ตัวแปร db สำหรับเก็บการเชื่อมต่อฐานข้อมูล

// InitDB ฟังก์ชันสำหรับเชื่อมต่อฐานข้อมูล SQLite และสร้างตาราง products หากยังไม่มี
func InitDB() { // ฟังก์ชันนี้จะถูกเรียกใช้ใน main.go
	var err error // ตัวแปรสำหรับเก็บ error
	db, err = sql.Open("sqlite", "./products.db") // เชื่อมต่อฐานข้อมูล SQLite ที่ไฟล์ products.db
	if err != nil { // ตรวจสอบ error ในการเชื่อมต่อ
		log.Fatal(err) // ถ้ามี error ให้ log ข้อความและหยุดโปรแกรม
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		product_code TEXT NOT NULL UNIQUE,
		barcode TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		update_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);` // คำสั่ง SQL สำหรับสร้างตาราง products หากยังไม่มี
	_, err = db.Exec(createTable) // รันคำสั่ง SQL เพื่อสร้างตาราง
	if err != nil { // ตรวจสอบ error ในการรันคำสั่ง SQL
		log.Fatal(err) // ถ้ามี error ให้ log ข้อความและหยุดโปรแกรม
	}
}
