package main

import (
	// "log"

	// "github.com/gofiber/fiber/v2"

	"fmt"
)

type Student struct {
	ID       int     `json:"id"`
	Username string  `json:"username"`
	Grade    float64 `json:"grade"`
	IsActive bool    `json:"is_active"`
}

func (s Student) GetInfo() string {
	return fmt.Sprintf("ID: %d\nUsername: %s\nGrade: %f\nIsActive: %v", s.ID, s.Username, s.Grade, s.IsActive)
}

func (s *Student) UpdateGrade(grade float64) {
	s.Grade = grade
}

func (s *Student) Activate()   { s.IsActive = true }
func (s *Student) Deactivate() { s.IsActive = false }

func main() {
	// app := fiber.New()
	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.SendString("Hello World")
	// })
	// port := ":3000"
	// log.Fatal(app.Listen(port))

	var string_var string = "Ini String"
	var int_var int = 12
	var float_var float64 = 2.34242242
	var bool_var bool = false
	var slice_var = []int{1, 2, 3}

	fmt.Println(string_var, int_var, float_var, bool_var, slice_var)

	mahasiswa := make(map[string]int)

	//Menambah ke map
	mahasiswa["Joko"] = 89
	mahasiswa["Wawan"] = 90
	mahasiswa["Andi"] = 88

	//Membaca dan cek keberadaan
	bacaNilai := mahasiswa["Joko"]
	fmt.Println(bacaNilai)
	if nilaiJoko, exist := mahasiswa["Joko"]; exist {
		fmt.Println("Nilai Joko adalah:", nilaiJoko)
	} else {
		fmt.Println("Joko belum memiliki nilai")
	}

	//Menghaus isi map
	delete(mahasiswa, "Joko")

	//Menelusuri isi map
	for nama, nilai := range mahasiswa {
		fmt.Printf("Nama: %s || Nilai: %d \n", nama, nilai)
	}

	a, b := 100, 150
	fmt.Printf("1. Nilai awal \n Nilai A: %d || Nilai B: %d \n", a, b)
	swapByValue(a, b)
	fmt.Printf("a. swap dengan pass by value \n Nilai A: %d || Nilai B: %d \n", a, b)
	swap(&a, &b)
	fmt.Printf("b. swap dengan pass by pointer \n Nilai A: %d || Nilai B: %d \n", a, b)

	makanan := make([]string, 0, 10)
	makanan = append(makanan, "Rendang", "Nasi_Goreng")
	fmt.Println("2. slice awal", makanan)
	updateSliceByValue(makanan, "Bakso")
	fmt.Println("a. update slice dengan pass by value \n", makanan)
	updateSlice(&makanan, "Bakso")
	fmt.Println("b. update slice dengan pass by pointer \n", makanan)

	u := Student{
		ID:       1,
		Username: "joko_123",
		Grade:    99.212,
		IsActive: true,
	}
	fmt.Println(u.GetInfo())

	newGrade := 80.12
	u.UpdateGrade(newGrade)

	fmt.Println("__Setelah di Update__\n",u.GetInfo())

	u.Deactivate()

	fmt.Println("__Setelah di Deactivate__\n",u.GetInfo())

	u.Activate()

	fmt.Println("__Setelah di Activate__\n",u.GetInfo())
	
}

func swap(a, b *int) {
	*a, *b = *b, *a
}

func swapByValue(a, b int) {
	a, b = b, a
}

func updateSlice(s *[]string, newItem string) {
	*s = append(*s, newItem)
}

func updateSliceByValue(s []string, newItem string) {
	s = append(s, newItem)
}
