package main

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var users []User
var students []Student
var nextID = 1

func findUserIndex(id int) int {
	for i := range users {
		if users[i].ID == id {
			return i
		}
	}
	return -1
}

func findStudentIndex(id int) int {
	for i := range students {
		if students[i].ID == id {
			return i
		}
	}
	return -1
}

func cocokPencarianUser(u User, kata string) bool {
	kata = strings.ToLower(kata)
	return strings.Contains(strings.ToLower(u.Username), kata) ||
		strings.Contains(strings.ToLower(u.Email), kata)
}

func cocokPencarianStudent(s Student, kata string) bool {
	kata = strings.ToLower(kata)
	return strings.Contains(strings.ToLower(s.Username), kata) ||
		strings.Contains(strings.ToLower(s.NIM), kata)
}

func paramID(c *fiber.Ctx) (int, bool) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return 0, false
	}
	return id, true
}

// Student
func listStudents(c *fiber.Ctx) error {
	q := parseListQuery(c)
	hasil := []Student{}

	for _, s := range students {
		if q.isActive != nil && s.IsActive != *q.isActive {
			continue
		}
		if q.Search != "" && !cocokPencarianStudent(s, q.Search) {
			continue
		}
		if q.GradeMin != nil && q.GradeMax != nil && !(*q.GradeMin <= s.Grade && s.Grade <= *q.GradeMax) {
			continue
		}
		hasil = append(hasil, s)
	}

	sort.SliceStable(hasil, func(i, j int) bool {
		var lebihKecil bool
		switch q.Sort {
		case "username":
			lebihKecil = hasil[i].Username < hasil[j].Username
		default:
			lebihKecil = hasil[i].ID < hasil[j].ID
		}
		if q.Order == "desc" {
			return !lebihKecil
		}
		return lebihKecil
	})
	total := len(hasil)
	totalPages := (total + q.Limit - 1) / q.Limit
	mulai := (q.Page - 1) * q.Limit
	if mulai > total {
		mulai = total
	}
	akhir := mulai + q.Limit
	if akhir > total {
		akhir = total
	}
	return okList(c, "daftar students berhasil diambil", hasil[mulai:akhir], &Meta{
		Page: q.Page, Limit: q.Limit, Total: total, TotalPages: totalPages,
	})
}

func getStudent(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}
	i := findStudentIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "student tidak ditemukan")
	}
	return ok(c, "student ditemukan", students[i])
}

func createStudent(c *fiber.Ctx) error {
	var req CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, err.Error())
	}
	errs := map[string]string{}
	req.Username = strings.TrimSpace(req.Username)
	req.NIM = strings.TrimSpace(req.NIM)
	if req.Username == "" {
		errs["username"] = "wajib diisi"
	}
	if req.NIM == ""  {
		errs["nim"] = "wajib diisi"
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}
	for _, s := range students {
		if strings.EqualFold(s.Username, req.Username) || strings.EqualFold(s.NIM, req.NIM) {
			return fail(c, fiber.StatusConflict, "NIM/Username sudah ada")
		}
	}
	baru := Student{
		ID:        nextID,
		Username:  req.Username,
		NIM: req.NIM,
		Grade: 0,
		IsActive:  true,
	}
	students = append(students, baru)
	nextID++
	return created(c, "student berhasil dibuat", baru,
		"/api/v1/students/"+strconv.Itoa(baru.ID))
}

func replaceStudent(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}
	i := findStudentIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "student tidak ditemukan")
	}
	var req ReplaceStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}
	errs := map[string]string{}
	if strings.TrimSpace(req.Username) == "" {
		errs["username"] = "wajib diisi pada PUT"
	}
	if strings.TrimSpace(req.NIM) == "" {
		errs["nim"] = "wajib diisi pada PUT"
	}
	if req.Grade == nil {
		errs["grade"] = "wajib diisi pada PUT"
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}
	students[i].Username = req.Username
	students[i].NIM = req.NIM
	students[i].Grade = *req.Grade
	students[i].IsActive = req.IsActive
	return ok(c, "student berhasil diganti seluruhnya", students[i])
}

func patchStudent(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}
	i := findStudentIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "student tidak ditemukan")
	}
	var req PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}
	
	if req.Username != nil {
		if strings.TrimSpace(*req.Username) == "" {
			return failValidation(c, map[string]string{"username": "tidak boleh kosong"})
		}
		students[i].Username = *req.Username
	}

	if req.NIM != nil {
		if strings.TrimSpace(*req.NIM) == "" {
			return failValidation(c, map[string]string{"nim": "tidak boleh kosong"})
		}
		students[i].NIM = *req.NIM
	}

	if req.Grade != nil {
		if *req.Grade < 0  {
			return failValidation(c, map[string]string{"grade": "tidak boleh kosong"})
		}
		students[i].Grade = *req.Grade
	}
	
	if req.IsActive != nil {
		students[i].IsActive = *req.IsActive
	}
	return ok(c, "student berhasil diperbarui sebagian", students[i])
}

func deleteStudent(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}
	i := findStudentIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "student tidak ditemukan")
	}
	students = append(students[:i], students[i+1:]...)
	return noContent(c)
}

// User
func listUsers(c *fiber.Ctx) error {
	q := parseListQuery(c)

	hasil := []User{}

	for _, u := range users {
		if q.isActive != nil && u.IsActive != *q.isActive {
			continue
		}
		if q.Search != "" && !cocokPencarianUser(u, q.Search) {
			continue
		}
		hasil = append(hasil, u)
	}

	sort.SliceStable(hasil, func(i, j int) bool {
		var lebihKecil bool
		switch q.Sort {
		case "username":
			lebihKecil = hasil[i].Username < hasil[j].Username
		case "email":
			lebihKecil = hasil[i].Email < hasil[j].Email
		case "created_at":
			lebihKecil = hasil[i].CreatedAt.Before(hasil[j].CreatedAt)
		default:
			lebihKecil = hasil[i].ID < hasil[j].ID
		}
		if q.Order == "desc" {
			return !lebihKecil
		}
		return lebihKecil
	})
	total := len(hasil)
	totalPages := (total + q.Limit - 1) / q.Limit
	mulai := (q.Page - 1) * q.Limit
	if mulai > total {
		mulai = total
	}
	akhir := mulai + q.Limit
	if akhir > total {
		akhir = total
	}
	return okList(c, "daftar user berhasil diambil", hasil[mulai:akhir], &Meta{
		Page: q.Page, Limit: q.Limit, Total: total, TotalPages: totalPages,
	})
}

func getUser(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}
	i := findUserIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "user tidak ditemukan")
	}
	return ok(c, "user ditemukan", users[i])
}

func createUser(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, err.Error())
	}
	errs := map[string]string{}
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	if req.Username == "" {
		errs["username"] = "wajib diisi"
	}
	if !strings.Contains(req.Email, "@") {
		errs["email"] = "format email tidak valid"
	}
	if len(req.Password) < 8 {
		errs["password"] = "minimal 8 karakter"
	}
	for _, u := range users {
		if strings.EqualFold(u.Username, req.Username) {
			errs["username"] = "sudah dipakai"
		}
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}
	baru := User{
		ID:        nextID,
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		IsActive:  true,
		CreatedAt: time.Now(),
	}
	users = append(users, baru)
	nextID++
	return created(c, "user berhasil dibuat", baru,
		"/api/v1/users/"+strconv.Itoa(baru.ID))
}

func replaceUser(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}
	i := findUserIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "user tidak ditemukan")
	}
	var req ReplaceUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}
	errs := map[string]string{}
	if strings.TrimSpace(req.Username) == "" {
		errs["username"] = "wajib diisi pada PUT"
	}
	if !strings.Contains(req.Email, "@") {
		errs["email"] = "wajib diisi dan berformat email pada PUT"
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}
	users[i].Username = req.Username
	users[i].Email = req.Email
	users[i].IsActive = req.IsActive
	return ok(c, "user berhasil diganti seluruhnya", users[i])
}

func patchUser(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}
	i := findUserIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "user tidak ditemukan")
	}
	var req PatchUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}
	if req.Username == nil && req.Email == nil && req.IsActive == nil {
		return fail(c, fiber.StatusBadRequest, "tidak ada field yang diubah")
	}
	if req.Username != nil {
		if strings.TrimSpace(*req.Username) == "" {
			return failValidation(c, map[string]string{"username": "tidak boleh kosong"})
		}
		users[i].Username = *req.Username
	}
	if req.Email != nil {
		if !strings.Contains(*req.Email, "@") {
			return failValidation(c, map[string]string{"email": "format email tidak valid"})
		}
		users[i].Email = *req.Email
	}
	if req.IsActive != nil {
		users[i].IsActive = *req.IsActive
	}
	return ok(c, "user berhasil diperbarui sebagian", users[i])
}

func deleteUser(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}
	i := findUserIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "user tidak ditemukan")
	}
	users = append(users[:i], users[i+1:]...)
	return noContent(c)
}
