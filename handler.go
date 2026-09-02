package main

import (
	"errors"
	"latihan-fiber/app/model"
	"latihan-fiber/app/repository"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	repo repository.UserRepository
}

type StudentHandler struct {
	repo repository.StudentRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func NewStudentHandler(repo repository.StudentRepository) *StudentHandler {
	return &StudentHandler{repo: repo}
}

func terjemahkanError(c *fiber.Ctx, err error, pesanUmum string) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return fail(c, fiber.StatusNotFound, "User not found")
	case errors.Is(err, repository.ErrDuplicate):
		return fail(c, fiber.StatusConflict, "Username already used")
	default:
		return fail(c, fiber.StatusInternalServerError, pesanUmum)
	}
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	q := parseListQuery(c)

	users, total, err := h.repo.FindAll(ctx, q)
	if err != nil {
		return fail(c, fiber.StatusInternalServerError, "Error fetching users")
	}

	totalPages := 0
	if q.Limit > 0 {
		totalPages = (total + q.Limit - 1) / q.Limit
	}

	return okList(c, "Successful fetching Users data", users, &model.Meta{
		Page: q.Page, Limit: q.Limit, Total: total, TotalPages: totalPages,
	})
}

func (h *UserHandler) Get(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "ID must be valid")
	}

	user, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return terjemahkanError(c, err, "Fail to get user data")
	}

	return ok(c, "User found!", user)
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	var req model.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "Body must be in valid JSON format!")
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	errs := map[string]string{}
	if req.Username == "" {
		errs["username"] = "must be filled"
	}
	if !strings.Contains(req.Email, "@") {
		errs["email"] = "format isn't valid"
	}
	if len(req.Password) < 8 {
		errs["password"] = "minimum in 8 characters"
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	baru, err := h.repo.Create(ctx, model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		IsActive: true,
	})
	if err != nil {
		return terjemahkanError(c, err, "Failed to save user data")
	}
	return created(c, "User succesfully created", baru,
		"/api/v1/users/"+strconv.Itoa(baru.ID))
}

func (h *UserHandler) Replace(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "ID must be valid!")
	}
	var req model.ReplaceUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "Body must be valid in JSON format")
	}
	errs := map[string]string{}
	if strings.TrimSpace(req.Username) == "" {
		errs["username"] = "must be filled"
	}
	if !strings.Contains(req.Email, "@") {
		errs["email"] = "must be filled with email format"
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	hasil, err := h.repo.Update(ctx, model.User{
		ID: id, Username: req.Username, Email: req.Email, IsActive: req.IsActive,
	})
	if err != nil {
		return terjemahkanError(c, err, "Failed to update user")
	}

	return ok(c, "User data succesfully replaced", hasil)
}

func (h *UserHandler) Patch(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "ID must be valid!")
	}

	var req model.PatchUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "Body must be in valid JSON format")
	}
	if req.Username == nil && req.Email == nil && req.IsActive == nil {
		return fail(c, fiber.StatusBadRequest, "No field is changed")
	}

	saatIni, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return terjemahkanError(c, err, "Failed to fetch user")
	}
	if req.Username != nil {
		if strings.TrimSpace(*req.Username) == "" {
			return failValidation(c, map[string]string{"username": "can't be empty"})
		}
		saatIni.Username = *req.Username
	}
	if req.Email != nil {
		if !strings.Contains(*req.Email, "@") {
			return failValidation(c, map[string]string{"email": "invalid email format"})
		}
		saatIni.Email = *req.Email
	}
	if req.IsActive != nil {
		saatIni.IsActive = *req.IsActive
	}
	hasil, err := h.repo.Update(ctx, saatIni)
	if err != nil {
		return terjemahkanError(c, err, "Failed to update user")
	}
	return ok(c, "user successfuly patch updated", hasil)
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "Failed to delete user")
	}

	if err := h.repo.Delete(ctx, id); err != nil {
		return terjemahkanError(c, err, "Failed to delete user")
	}
	return noContent(c)
}

func (h *StudentHandler) ListStudents(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	q := parseListQuery(c)

	users, total, err := h.repo.FindAll(ctx, q)
	if err != nil {
		return fail(c, fiber.StatusInternalServerError, "Error fetching users")
	}

	totalPages := 0
	if q.Limit > 0 {
		totalPages = (total + q.Limit - 1) / q.Limit
	}

	return okList(c, "Successful fetching Students data", users, &model.Meta{
		Page: q.Page, Limit: q.Limit, Total: total, TotalPages: totalPages,
	})
}

func (h *StudentHandler) GetStudent(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "ID must be valid")
	}

	user, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return terjemahkanError(c, err, "Fail to get student data")
	}

	return ok(c, "Student found!", user)
}

func (h *StudentHandler) CreateStudent(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	var req model.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "Body must be in valid JSON format!")
	}

	req.Username = strings.TrimSpace(req.Username)
	req.NIM = strings.TrimSpace(req.NIM)

	errs := map[string]string{}
	if req.Username == "" {
		errs["username"] = "must be filled"
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}
	baru, err := h.repo.Create(ctx, model.Student{
		Username: req.Username,
		NIM:    req.NIM,
		Grade: *req.Grade,
		IsActive: true,
	})
	if err != nil {
		return terjemahkanError(c, err, "Failed to save user data")
	}
	return created(c, "User succesfully created", baru,
		"/api/v1/users/"+strconv.Itoa(baru.ID))
}

func (h *StudentHandler) ReplaceStudent(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "ID must be valid!")
	}
	var req model.ReplaceStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "Body must be valid in JSON format")
	}
	errs := map[string]string{}
	if strings.TrimSpace(req.Username) == "" {
		errs["username"] = "must be filled"
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	hasil, err := h.repo.Update(ctx, model.Student{
		ID: id, Username: req.Username, NIM: req.NIM, Grade: *req.Grade, IsActive: req.IsActive,
	})
	if err != nil {
		return terjemahkanError(c, err, "Failed to update student")
	}

	return ok(c, "Student data succesfully replaced", hasil)
}

func (h *StudentHandler) PatchStudent(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "ID must be valid!")
	}

	var req model.PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "Body must be in valid JSON format")
	}
	if req.Username == nil && req.NIM == nil && req.Grade == nil && req.IsActive == nil {
		return fail(c, fiber.StatusBadRequest, "No field is changed")
	}

	saatIni, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return terjemahkanError(c, err, "Failed to fetch student")
	}
	if req.Username != nil {
		if strings.TrimSpace(*req.Username) == "" {
			return failValidation(c, map[string]string{"username": "can't be empty"})
		}
		saatIni.Username = *req.Username
	}
	if req.IsActive != nil {
		saatIni.IsActive = *req.IsActive
	}
	hasil, err := h.repo.Update(ctx, saatIni)
	if err != nil {
		return terjemahkanError(c, err, "Failed to update student")
	}
	return ok(c, "student successfuly patch updated", hasil)
}

func (h *StudentHandler) DeleteStudent(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "Failed to delete student")
	}

	if err := h.repo.Delete(ctx, id); err != nil {
		return terjemahkanError(c, err, "Failed to delete student")
	}
	return noContent(c)
}

// var users []User
// var students []Student
// // var nextID = 1

// func findUserIndex(id int) int {
// 	for i := range users {
// 		if users[i].ID == id {
// 			return i
// 		}
// 	}
// 	return -1
// }

// func findStudentIndex(id int) int {
// 	for i := range students {
// 		if students[i].ID == id {
// 			return i
// 		}
// 	}
// 	return -1
// }

// func cocokPencarianUser(u User, kata string) bool {
// 	kata = strings.ToLower(kata)
// 	return strings.Contains(strings.ToLower(u.Username), kata) ||
// 		strings.Contains(strings.ToLower(u.Email), kata)
// }

// func cocokPencarianStudent(s Student, kata string) bool {
// 	kata = strings.ToLower(kata)
// 	return strings.Contains(strings.ToLower(s.Username), kata) ||
// 		strings.Contains(strings.ToLower(s.NIM), kata)
// }

func paramID(c *fiber.Ctx) (int, bool) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return 0, false
	}
	return id, true
}

// // Student
// func listStudents(c *fiber.Ctx) error {
// 	q := parseListQuery(c)
// 	hasil := []Student{}

// 	for _, s := range students {
// 		if q.isActive != nil && s.IsActive != *q.isActive {
// 			continue
// 		}
// 		if q.Search != "" && !cocokPencarianStudent(s, q.Search) {
// 			continue
// 		}
// 		if q.GradeMin != nil && q.GradeMax != nil && !(*q.GradeMin <= s.Grade && s.Grade <= *q.GradeMax) {
// 			continue
// 		}
// 		hasil = append(hasil, s)
// 	}

// 	sort.SliceStable(hasil, func(i, j int) bool {
// 		var lebihKecil bool
// 		switch q.Sort {
// 		case "username":
// 			lebihKecil = hasil[i].Username < hasil[j].Username
// 		default:
// 			lebihKecil = hasil[i].ID < hasil[j].ID
// 		}
// 		if q.Order == "desc" {
// 			return !lebihKecil
// 		}
// 		return lebihKecil
// 	})
// 	total := len(hasil)
// 	totalPages := (total + q.Limit - 1) / q.Limit
// 	mulai := (q.Page - 1) * q.Limit
// 	if mulai > total {
// 		mulai = total
// 	}
// 	akhir := mulai + q.Limit
// 	if akhir > total {
// 		akhir = total
// 	}
// 	return okList(c, "daftar students berhasil diambil", hasil[mulai:akhir], &Meta{
// 		Page: q.Page, Limit: q.Limit, Total: total, TotalPages: totalPages,
// 	})
// }

// func getStudent(c *fiber.Ctx) error {
// 	id, valid := paramID(c)
// 	if !valid {
// 		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
// 	}
// 	i := findStudentIndex(id)
// 	if i == -1 {
// 		return fail(c, fiber.StatusNotFound, "student tidak ditemukan")
// 	}
// 	return ok(c, "student ditemukan", students[i])
// }

// func createStudent(c *fiber.Ctx) error {
// 	var req CreateStudentRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return fail(c, fiber.StatusBadRequest, err.Error())
// 	}
// 	errs := map[string]string{}
// 	req.Username = strings.TrimSpace(req.Username)
// 	req.NIM = strings.TrimSpace(req.NIM)
// 	if req.Username == "" {
// 		errs["username"] = "wajib diisi"
// 	}
// 	if req.NIM == "" {
// 		errs["nim"] = "wajib diisi"
// 	}
// 	if len(errs) > 0 {
// 		return failValidation(c, errs)
// 	}
// 	for _, s := range students {
// 		if strings.EqualFold(s.Username, req.Username) || strings.EqualFold(s.NIM, req.NIM) {
// 			return fail(c, fiber.StatusConflict, "NIM/Username sudah ada")
// 		}
// 	}
// 	baru := Student{
// 		ID:       nextID,
// 		Username: req.Username,
// 		NIM:      req.NIM,
// 		Grade:    0,
// 		IsActive: true,
// 	}
// 	students = append(students, baru)
// 	nextID++
// 	return created(c, "student berhasil dibuat", baru,
// 		"/api/v1/students/"+strconv.Itoa(baru.ID))
// }

// func replaceStudent(c *fiber.Ctx) error {
// 	id, valid := paramID(c)
// 	if !valid {
// 		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
// 	}
// 	i := findStudentIndex(id)
// 	if i == -1 {
// 		return fail(c, fiber.StatusNotFound, "student tidak ditemukan")
// 	}
// 	var req ReplaceStudentRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
// 	}
// 	errs := map[string]string{}
// 	if strings.TrimSpace(req.Username) == "" {
// 		errs["username"] = "wajib diisi pada PUT"
// 	}
// 	if strings.TrimSpace(req.NIM) == "" {
// 		errs["nim"] = "wajib diisi pada PUT"
// 	}
// 	if req.Grade == nil {
// 		errs["grade"] = "wajib diisi pada PUT"
// 	}
// 	if len(errs) > 0 {
// 		return failValidation(c, errs)
// 	}
// 	students[i].Username = req.Username
// 	students[i].NIM = req.NIM
// 	students[i].Grade = *req.Grade
// 	students[i].IsActive = req.IsActive
// 	return ok(c, "student berhasil diganti seluruhnya", students[i])
// }

// func patchStudent(c *fiber.Ctx) error {
// 	id, valid := paramID(c)
// 	if !valid {
// 		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
// 	}
// 	i := findStudentIndex(id)
// 	if i == -1 {
// 		return fail(c, fiber.StatusNotFound, "student tidak ditemukan")
// 	}
// 	var req PatchStudentRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
// 	}

// 	if req.Username != nil {
// 		if strings.TrimSpace(*req.Username) == "" {
// 			return failValidation(c, map[string]string{"username": "tidak boleh kosong"})
// 		}
// 		students[i].Username = *req.Username
// 	}

// 	if req.NIM != nil {
// 		if strings.TrimSpace(*req.NIM) == "" {
// 			return failValidation(c, map[string]string{"nim": "tidak boleh kosong"})
// 		}
// 		students[i].NIM = *req.NIM
// 	}

// 	if req.Grade != nil {
// 		if *req.Grade < 0 {
// 			return failValidation(c, map[string]string{"grade": "tidak boleh kosong"})
// 		}
// 		students[i].Grade = *req.Grade
// 	}

// 	if req.IsActive != nil {
// 		students[i].IsActive = *req.IsActive
// 	}
// 	return ok(c, "student berhasil diperbarui sebagian", students[i])
// }

// func deleteStudent(c *fiber.Ctx) error {
// 	id, valid := paramID(c)
// 	if !valid {
// 		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
// 	}
// 	i := findStudentIndex(id)
// 	if i == -1 {
// 		return fail(c, fiber.StatusNotFound, "student tidak ditemukan")
// 	}
// 	students = append(students[:i], students[i+1:]...)
// 	return noContent(c)
// }

// // User
// func listUsers(c *fiber.Ctx) error {
// 	q := parseListQuery(c)

// 	hasil := []User{}

// 	for _, u := range users {
// 		if q.isActive != nil && u.IsActive != *q.isActive {
// 			continue
// 		}
// 		if q.Search != "" && !cocokPencarianUser(u, q.Search) {
// 			continue
// 		}
// 		hasil = append(hasil, u)
// 	}

// 	sort.SliceStable(hasil, func(i, j int) bool {
// 		var lebihKecil bool
// 		switch q.Sort {
// 		case "username":
// 			lebihKecil = hasil[i].Username < hasil[j].Username
// 		case "email":
// 			lebihKecil = hasil[i].Email < hasil[j].Email
// 		case "created_at":
// 			lebihKecil = hasil[i].CreatedAt.Before(hasil[j].CreatedAt)
// 		default:
// 			lebihKecil = hasil[i].ID < hasil[j].ID
// 		}
// 		if q.Order == "desc" {
// 			return !lebihKecil
// 		}
// 		return lebihKecil
// 	})
// 	total := len(hasil)
// 	totalPages := (total + q.Limit - 1) / q.Limit
// 	mulai := (q.Page - 1) * q.Limit
// 	if mulai > total {
// 		mulai = total
// 	}
// 	akhir := mulai + q.Limit
// 	if akhir > total {
// 		akhir = total
// 	}
// 	return okList(c, "daftar user berhasil diambil", hasil[mulai:akhir], &Meta{
// 		Page: q.Page, Limit: q.Limit, Total: total, TotalPages: totalPages,
// 	})
// }

// func getUser(c *fiber.Ctx) error {
// 	id, valid := paramID(c)
// 	if !valid {
// 		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
// 	}
// 	i := findUserIndex(id)
// 	if i == -1 {
// 		return fail(c, fiber.StatusNotFound, "user tidak ditemukan")
// 	}
// 	return ok(c, "user ditemukan", users[i])
// }

// func createUser(c *fiber.Ctx) error {
// 	var req CreateUserRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return fail(c, fiber.StatusBadRequest, err.Error())
// 	}
// 	errs := map[string]string{}
// 	req.Username = strings.TrimSpace(req.Username)
// 	req.Email = strings.TrimSpace(req.Email)
// 	if req.Username == "" {
// 		errs["username"] = "wajib diisi"
// 	}
// 	if !strings.Contains(req.Email, "@") {
// 		errs["email"] = "format email tidak valid"
// 	}
// 	if len(req.Password) < 8 {
// 		errs["password"] = "minimal 8 karakter"
// 	}
// 	for _, u := range users {
// 		if strings.EqualFold(u.Username, req.Username) {
// 			errs["username"] = "sudah dipakai"
// 		}
// 	}
// 	if len(errs) > 0 {
// 		return failValidation(c, errs)
// 	}
// 	baru := User{
// 		ID:        nextID,
// 		Username:  req.Username,
// 		Email:     req.Email,
// 		Password:  req.Password,
// 		IsActive:  true,
// 		CreatedAt: time.Now(),
// 	}
// 	users = append(users, baru)
// 	nextID++
// 	return created(c, "user berhasil dibuat", baru,
// 		"/api/v1/users/"+strconv.Itoa(baru.ID))
// }

// func replaceUser(c *fiber.Ctx) error {
// 	id, valid := paramID(c)
// 	if !valid {
// 		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
// 	}
// 	i := findUserIndex(id)
// 	if i == -1 {
// 		return fail(c, fiber.StatusNotFound, "user tidak ditemukan")
// 	}
// 	var req ReplaceUserRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
// 	}
// 	errs := map[string]string{}
// 	if strings.TrimSpace(req.Username) == "" {
// 		errs["username"] = "wajib diisi pada PUT"
// 	}
// 	if !strings.Contains(req.Email, "@") {
// 		errs["email"] = "wajib diisi dan berformat email pada PUT"
// 	}
// 	if len(errs) > 0 {
// 		return failValidation(c, errs)
// 	}
// 	users[i].Username = req.Username
// 	users[i].Email = req.Email
// 	users[i].IsActive = req.IsActive
// 	return ok(c, "user berhasil diganti seluruhnya", users[i])
// }

// func patchUser(c *fiber.Ctx) error {
// 	id, valid := paramID(c)
// 	if !valid {
// 		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
// 	}
// 	i := findUserIndex(id)
// 	if i == -1 {
// 		return fail(c, fiber.StatusNotFound, "user tidak ditemukan")
// 	}
// 	var req PatchUserRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
// 	}
// 	if req.Username == nil && req.Email == nil && req.IsActive == nil {
// 		return fail(c, fiber.StatusBadRequest, "tidak ada field yang diubah")
// 	}
// 	if req.Username != nil {
// 		if strings.TrimSpace(*req.Username) == "" {
// 			return failValidation(c, map[string]string{"username": "tidak boleh kosong"})
// 		}
// 		users[i].Username = *req.Username
// 	}
// 	if req.Email != nil {
// 		if !strings.Contains(*req.Email, "@") {
// 			return failValidation(c, map[string]string{"email": "format email tidak valid"})
// 		}
// 		users[i].Email = *req.Email
// 	}
// 	if req.IsActive != nil {
// 		users[i].IsActive = *req.IsActive
// 	}
// 	return ok(c, "user berhasil diperbarui sebagian", users[i])
// }

// func deleteUser(c *fiber.Ctx) error {
// 	id, valid := paramID(c)
// 	if !valid {
// 		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
// 	}
// 	i := findUserIndex(id)
// 	if i == -1 {
// 		return fail(c, fiber.StatusNotFound, "user tidak ditemukan")
// 	}
// 	users = append(users[:i], users[i+1:]...)
// 	return noContent(c)
// }
