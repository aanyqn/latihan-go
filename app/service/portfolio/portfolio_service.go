package portfolio

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"latihan-fiber/app/model"
	"latihan-fiber/app/repository"
	"latihan-fiber/app/service/authz"
	"latihan-fiber/helper"
)

// PortfolioService memakai DUA database: data mahasiswa diperiksa di
// PostgreSQL, portofolionya disimpan di MongoDB. Hanya service yang mengenal
// keduanya; repository tetap terpisah per database.
type PortfolioService struct {
	portfolios repository.PortfolioRepository
	students   repository.StudentRepository
	perms      *helper.PermissionSet
}

func NewPortfolioService(
	portfolios repository.PortfolioRepository,
	students repository.StudentRepository,
	perms *helper.PermissionSet,
) *PortfolioService {
	return &PortfolioService{portfolios: portfolios, students: students, perms: perms}
}

// authorize memastikan mahasiswa ada di PostgreSQL sekaligus memeriksa hak
// akses. Pemilik data boleh mengakses; role lain membutuhkan permission yang
// disebutkan (student:read:any atau student:update:any).
func (s *PortfolioService) authorize(
	ctx context.Context, c *fiber.Ctx, anyPermission string,
) (model.Student, error) {
	current, ok := helper.CurrentUser(c)
	if !ok {
		return model.Student{}, helper.Unauthorized("belum terautentikasi")
	}
	studentID, valid := helper.ParamID(c)
	if !valid || studentID < 1 {
		return model.Student{}, helper.BadRequest("id mahasiswa harus berupa angka positif")
	}
	student, err := s.students.FindByID(ctx, studentID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return model.Student{}, helper.NotFound("mahasiswa tidak ditemukan")
		}
		return model.Student{}, helper.Internal(err)
	}
	if !authz.CanAccessStudent(current, student.OwnerID, s.perms, anyPermission) {
		return model.Student{}, helper.Forbidden(
			"Anda tidak berhak mengakses portofolio mahasiswa ini")
	}
	return student, nil
}

// ---------- GET /students/:id/portfolio ----------
func (s *PortfolioService) Get(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	student, err := s.authorize(ctx, c, "student:read:any")
	if err != nil {
		return err
	}
	portfolio, err := s.portfolios.FindByStudentID(ctx, student.ID)
	if err != nil {
		return translatePortfolioError(err)
	}
	return helper.Success(c, fiber.StatusOK, "portofolio ditemukan", portfolio)
}

// ---------- PUT /students/:id/portfolio ----------
// PUT bersifat mengganti seluruh isi portofolio. Bila dokumennya belum ada,
// repository membuatnya dengan upsert.
func (s *PortfolioService) Replace(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	student, err := s.authorize(ctx, c, "student:update:any")
	if err != nil {
		return err
	}

	var req model.ReplacePortfolioRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}
	normalizeReplaceRequest(&req)
	if errs := helper.ValidateStruct(req); errs != nil {
		return helper.Validation(errs)
	}

	projects := make([]model.Project, 0, len(req.Projects))
	for _, p := range req.Projects {
		projects = append(projects, model.Project{
			Title:        p.Title,
			Technologies: p.Technologies,
			Link:         p.Link,
			Year:         p.Year,
		})
	}
	certificates := make([]model.Certificate, 0, len(req.Certificates))
	for _, cert := range req.Certificates {
		issuedAt, err := parseCertificateDate(cert.IssuedAt)
		if err != nil {
			return helper.Validation(map[string]string{
				"certificates.issued_at": "tanggal sertifikat harus berformat YYYY-MM-DD atau RFC3339",
			})
		}
		certificates = append(certificates, model.Certificate{
			Name:     cert.Name,
			Issuer:   cert.Issuer,
			IssuedAt: issuedAt,
		})
	}

	saved, err := s.portfolios.Replace(ctx, model.Portfolio{
		StudentID:    student.ID,
		Headline:     req.Headline,
		Bio:          req.Bio,
		Projects:     projects,
		Certificates: certificates,
	})
	if err != nil {
		return translatePortfolioError(err)
	}
	return helper.Success(c, fiber.StatusOK, "portofolio berhasil disimpan", saved)
}

// ---------- POST /students/:id/portfolio/projects ----------
func (s *PortfolioService) AddProject(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	student, err := s.authorize(ctx, c, "student:update:any")
	if err != nil {
		return err
	}

	var req model.ProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}
	normalizeProjectRequest(&req)
	if errs := helper.ValidateStruct(req); errs != nil {
		return helper.Validation(errs)
	}

	updated, err := s.portfolios.AddProject(ctx, student.ID, model.Project{
		Title:        req.Title,
		Technologies: req.Technologies,
		Link:         req.Link,
		Year:         req.Year,
	})
	if err != nil {
		return translatePortfolioError(err)
	}
	return helper.Success(c, fiber.StatusCreated, "proyek berhasil ditambahkan", updated)
}

// ---------- DELETE /students/:id/portfolio/projects/:projectId ----------
func (s *PortfolioService) RemoveProject(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	student, err := s.authorize(ctx, c, "student:update:any")
	if err != nil {
		return err
	}
	projectID := strings.TrimSpace(c.Params("projectId"))
	if projectID == "" {
		return helper.BadRequest("id proyek wajib diisi")
	}

	updated, err := s.portfolios.RemoveProject(ctx, student.ID, projectID)
	if err != nil {
		return translatePortfolioError(err)
	}
	return helper.Success(c, fiber.StatusOK, "proyek berhasil dihapus", updated)
}

// ---------- GET /portfolios/stats/technologies ----------
func (s *PortfolioService) TopTechnologies(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	stats, err := s.portfolios.TopTechnologies(ctx, 10)
	if err != nil {
		return helper.Internal(err)
	}
	return helper.Success(c, fiber.StatusOK, "statistik teknologi berhasil diambil", stats)
}

// NormalizeTechnologies merapikan nama teknologi: huruf kecil, tanpa spasi di
// tepi, tanpa duplikat, urutan pertama dipertahankan. Tanpa normalisasi,
// "Go", "go", dan " go " akan dihitung sebagai tiga teknologi berbeda.
func NormalizeTechnologies(technologies []string) []string {
	result := make([]string, 0, len(technologies))
	seen := make(map[string]struct{}, len(technologies))
	for _, t := range technologies {
		t = strings.ToLower(strings.TrimSpace(t))
		if t == "" {
			continue
		}
		if _, dup := seen[t]; dup {
			continue
		}
		seen[t] = struct{}{}
		result = append(result, t)
	}
	return result
}

func normalizeProjectRequest(req *model.ProjectRequest) {
	req.Title = strings.TrimSpace(req.Title)
	req.Link = strings.TrimSpace(req.Link)
	req.Technologies = NormalizeTechnologies(req.Technologies)
}

func normalizeReplaceRequest(req *model.ReplacePortfolioRequest) {
	req.Headline = strings.TrimSpace(req.Headline)
	req.Bio = strings.TrimSpace(req.Bio)
	for i := range req.Projects {
		normalizeProjectRequest(&req.Projects[i])
	}
	for i := range req.Certificates {
		req.Certificates[i].Name = strings.TrimSpace(req.Certificates[i].Name)
		req.Certificates[i].Issuer = strings.TrimSpace(req.Certificates[i].Issuer)
		req.Certificates[i].IssuedAt = strings.TrimSpace(req.Certificates[i].IssuedAt)
	}
}

// parseCertificateDate menerima tanggal saja maupun RFC3339. Waktu selalu
// disimpan dalam UTC dan dipotong ke milidetik, sama seperti now() pada
// repository ulasan.
func parseCertificateDate(raw string) (time.Time, error) {
	layouts := []string{time.RFC3339, "2006-01-02", "2006-01-02 15:04:05"}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return t.UTC().Truncate(time.Millisecond), nil
		}
	}
	return time.Time{}, errors.New("format tanggal tidak valid")
}

func translatePortfolioError(err error) error {
	switch {
	case errors.Is(err, repository.ErrInvalidID):
		return helper.BadRequest("id proyek tidak valid")
	case errors.Is(err, repository.ErrPortfolioNotFound):
		return helper.NotFound("portofolio belum dibuat")
	case errors.Is(err, repository.ErrProjectNotFound):
		return helper.NotFound("proyek tidak ditemukan")
	case errors.Is(err, repository.ErrProjectLimit):
		return helper.NewError(fiber.StatusConflict, "PROJECT_LIMIT",
			"jumlah proyek sudah mencapai batas 20")
	default:
		return helper.Internal(err)
	}
}
