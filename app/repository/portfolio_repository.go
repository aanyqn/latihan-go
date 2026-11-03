package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"latihan-fiber/app/model"
)

var (
	ErrPortfolioNotFound = errors.New("portofolio tidak ditemukan")
	ErrProjectNotFound   = errors.New("proyek tidak ditemukan")
	ErrProjectLimit      = errors.New("jumlah proyek sudah mencapai batas")
)

// MaxProjectsPerPortfolio membatasi pertumbuhan array projects. Satu dokumen
// MongoDB maksimal 16 MB; batas ini juga menjaga ukuran dokumen tetap wajar.
const MaxProjectsPerPortfolio = 20

// PortfolioRepository adalah satu-satunya pintu akses ke collection
// student_portfolios. Service tidak pernah menyentuh *mongo.Collection.
type PortfolioRepository interface {
	EnsureIndexes(ctx context.Context) error
	FindByStudentID(ctx context.Context, studentID int) (model.Portfolio, error)
	Replace(ctx context.Context, p model.Portfolio) (model.Portfolio, error)
	AddProject(ctx context.Context, studentID int, project model.Project) (model.Portfolio, error)
	RemoveProject(ctx context.Context, studentID int, projectID string) (model.Portfolio, error)
	TopTechnologies(ctx context.Context, limit int) ([]model.TechnologyStat, error)
}

// portfolioDocument adalah bentuk portofolio SEBAGAIMANA TERSIMPAN di
// MongoDB. Tag bson hanya hidup di repository; package model tetap bebas dari
// tipe milik driver.
type portfolioDocument struct {
	ID           bson.ObjectID         `bson:"_id,omitempty"`
	StudentID    int                   `bson:"student_id"`
	Headline     string                `bson:"headline"`
	Bio          string                `bson:"bio"`
	Projects     []projectDocument     `bson:"projects"`
	Certificates []certificateDocument `bson:"certificates"`
	CreatedAt    time.Time             `bson:"created_at"`
	UpdatedAt    time.Time             `bson:"updated_at"`
}

type projectDocument struct {
	ID           bson.ObjectID `bson:"_id"`
	Title        string        `bson:"title"`
	Technologies []string      `bson:"technologies"`
	Link         string        `bson:"link"`
	Year         int           `bson:"year"`
}

type certificateDocument struct {
	Name     string    `bson:"name"`
	Issuer   string    `bson:"issuer"`
	IssuedAt time.Time `bson:"issued_at"`
}

func (d portfolioDocument) toModel() model.Portfolio {
	projects := make([]model.Project, 0, len(d.Projects))
	for _, p := range d.Projects {
		technologies := p.Technologies
		if technologies == nil {
			technologies = []string{}
		}
		projects = append(projects, model.Project{
			ID:           p.ID.Hex(),
			Title:        p.Title,
			Technologies: technologies,
			Link:         p.Link,
			Year:         p.Year,
		})
	}

	certificates := make([]model.Certificate, 0, len(d.Certificates))
	for _, c := range d.Certificates {
		certificates = append(certificates, model.Certificate{
			Name:     c.Name,
			Issuer:   c.Issuer,
			IssuedAt: c.IssuedAt,
		})
	}

	return model.Portfolio{
		ID:           d.ID.Hex(),
		StudentID:    d.StudentID,
		Headline:     d.Headline,
		Bio:          d.Bio,
		Projects:     projects,
		Certificates: certificates,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}

type portfolioMongoRepository struct {
	coll *mongo.Collection
}

func NewPortfolioRepository(db *mongo.Database) PortfolioRepository {
	return &portfolioMongoRepository{coll: db.Collection("student_portfolios")}
}

// EnsureIndexes membuat aturan "satu mahasiswa satu portofolio" di level
// database, bukan hanya di kode. Index kedua adalah multikey index pada array
// technologies untuk pencarian/filter teknologi di masa depan.
func (r *portfolioMongoRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "student_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("student_portfolio_unique"),
		},
		{
			Keys:    bson.D{{Key: "projects.technologies", Value: 1}},
			Options: options.Index().SetName("portfolio_technologies"),
		},
	})
	if err != nil {
		return fmt.Errorf("membuat index student_portfolios: %w", err)
	}
	return nil
}

func (r *portfolioMongoRepository) FindByStudentID(
	ctx context.Context, studentID int,
) (model.Portfolio, error) {
	var doc portfolioDocument
	err := r.coll.FindOne(ctx, bson.M{"student_id": studentID}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Portfolio{}, ErrPortfolioNotFound
		}
		return model.Portfolio{}, fmt.Errorf("mengambil portofolio: %w", err)
	}
	return doc.toModel(), nil
}

// Replace membuat atau mengganti portofolio dalam SATU operasi upsert.
//
// $setOnInsert dipakai agar created_at hanya diisi saat dokumen pertama kali
// dibuat; saat mengganti portofolio yang sudah ada, created_at lama tetap
// dipertahankan. Filter student_id menjadi kunci upsert, dan unique index
// memastikan dua request serentak tidak dapat membuat dua dokumen.
func (r *portfolioMongoRepository) Replace(
	ctx context.Context, p model.Portfolio,
) (model.Portfolio, error) {
	t := now()
	set := bson.M{
		"headline":     p.Headline,
		"bio":          p.Bio,
		"projects":     projectDocumentsFromModel(p.Projects),
		"certificates": certificateDocumentsFromModel(p.Certificates),
		"updated_at":   t,
	}
	update := bson.M{
		"$set": set,
		"$setOnInsert": bson.M{
			"_id":        bson.NewObjectID(),
			"created_at": t,
		},
	}

	var doc portfolioDocument
	err := r.coll.FindOneAndUpdate(ctx,
		bson.M{"student_id": p.StudentID},
		update,
		options.FindOneAndUpdate().
			SetUpsert(true).
			SetReturnDocument(options.After),
	).Decode(&doc)

	// Dua PUT serentak dapat sama-sama melihat "belum ada dokumen" lalu
	// mencoba insert. Yang kalah mendapat duplicate key; ulangi sebagai update
	// biasa karena dokumennya sekarang sudah ada.
	if mongo.IsDuplicateKeyError(err) {
		err = r.coll.FindOneAndUpdate(ctx,
			bson.M{"student_id": p.StudentID},
			bson.M{"$set": set},
			options.FindOneAndUpdate().SetReturnDocument(options.After),
		).Decode(&doc)
	}
	if err != nil {
		return model.Portfolio{}, fmt.Errorf("menyimpan portofolio: %w", err)
	}
	return doc.toModel(), nil
}

// AddProject menambahkan satu proyek dengan $push.
//
// Syarat "projects.19 tidak ada" berarti array belum berisi 20 elemen
// (indeks 0 sampai 19). Pemeriksaan batas dan penambahan terjadi dalam satu
// operasi atomik, sehingga 10 request serentak tidak dapat melewati batas.
func (r *portfolioMongoRepository) AddProject(
	ctx context.Context, studentID int, project model.Project,
) (model.Portfolio, error) {
	doc := projectDocument{
		ID:           bson.NewObjectID(),
		Title:        project.Title,
		Technologies: project.Technologies,
		Link:         project.Link,
		Year:         project.Year,
	}
	if doc.Technologies == nil {
		doc.Technologies = []string{}
	}

	limitField := "projects." + strconv.Itoa(MaxProjectsPerPortfolio-1)
	var updated portfolioDocument
	err := r.coll.FindOneAndUpdate(ctx,
		bson.M{
			"student_id": studentID,
			limitField:   bson.M{"$exists": false},
		},
		bson.M{
			"$push": bson.M{"projects": doc},
			"$set":  bson.M{"updated_at": now()},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updated)
	if err == nil {
		return updated.toModel(), nil
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return model.Portfolio{}, fmt.Errorf("menambah proyek: %w", err)
	}

	// Tidak ada dokumen yang cocok. Bedakan "portofolionya belum ada" dari
	// "portofolionya ada tetapi sudah berisi 20 proyek".
	count, err := r.coll.CountDocuments(ctx, bson.M{"student_id": studentID})
	if err != nil {
		return model.Portfolio{}, fmt.Errorf("memeriksa portofolio: %w", err)
	}
	if count == 0 {
		return model.Portfolio{}, ErrPortfolioNotFound
	}
	return model.Portfolio{}, ErrProjectLimit
}

// RemoveProject menghapus satu proyek dengan $pull. Filter memastikan
// portofolio milik mahasiswa yang dimaksud dan proyeknya benar-benar ada.
func (r *portfolioMongoRepository) RemoveProject(
	ctx context.Context, studentID int, projectID string,
) (model.Portfolio, error) {
	oid, err := parseID(projectID)
	if err != nil {
		return model.Portfolio{}, err
	}

	var updated portfolioDocument
	err = r.coll.FindOneAndUpdate(ctx,
		bson.M{
			"student_id":   studentID,
			"projects._id": oid,
		},
		bson.M{
			"$pull": bson.M{"projects": bson.M{"_id": oid}},
			"$set":  bson.M{"updated_at": now()},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updated)
	if err == nil {
		return updated.toModel(), nil
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return model.Portfolio{}, fmt.Errorf("menghapus proyek: %w", err)
	}

	// Bedakan portofolio yang tidak ada dari proyek yang tidak ada.
	count, err := r.coll.CountDocuments(ctx, bson.M{"student_id": studentID})
	if err != nil {
		return model.Portfolio{}, fmt.Errorf("memeriksa portofolio: %w", err)
	}
	if count == 0 {
		return model.Portfolio{}, ErrPortfolioNotFound
	}
	return model.Portfolio{}, ErrProjectNotFound
}

// TopTechnologies menghitung 10 teknologi paling banyak dipakai di seluruh
// portofolio.
//
// Dua $unwind dibutuhkan karena technologies bersarang di dalam array
// projects: $unwind pertama memecah projects, $unwind kedua memecah array
// technologies di setiap proyek. Setelah itu $group menghitung kemunculan
// setiap teknologi.
func (r *portfolioMongoRepository) TopTechnologies(
	ctx context.Context, limit int,
) ([]model.TechnologyStat, error) {
	if limit <= 0 {
		limit = 10
	}
	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$projects"}},
		{{Key: "$unwind", Value: "$projects.technologies"}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$projects.technologies",
			"total": bson.M{"$sum": 1},
		}}},
		{{Key: "$sort", Value: bson.D{
			{Key: "total", Value: -1},
			{Key: "_id", Value: 1},
		}}},
		{{Key: "$limit", Value: int64(limit)}},
	}

	cursor, err := r.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("menghitung teknologi: %w", err)
	}
	defer cursor.Close(ctx)

	var rows []struct {
		Technology string `bson:"_id"`
		Total      int    `bson:"total"`
	}
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, fmt.Errorf("membaca statistik teknologi: %w", err)
	}

	result := make([]model.TechnologyStat, 0, len(rows))
	for _, row := range rows {
		result = append(result, model.TechnologyStat{
			Technology: row.Technology,
			Total:      row.Total,
		})
	}
	return result, nil
}

func projectDocumentsFromModel(projects []model.Project) []projectDocument {
	result := make([]projectDocument, 0, len(projects))
	for _, p := range projects {
		oid := bson.NewObjectID()
		if p.ID != "" {
			if parsed, err := bson.ObjectIDFromHex(p.ID); err == nil {
				oid = parsed
			}
		}
		technologies := p.Technologies
		if technologies == nil {
			technologies = []string{}
		}
		result = append(result, projectDocument{
			ID:           oid,
			Title:        p.Title,
			Technologies: technologies,
			Link:         p.Link,
			Year:         p.Year,
		})
	}
	return result
}

func certificateDocumentsFromModel(certificates []model.Certificate) []certificateDocument {
	result := make([]certificateDocument, 0, len(certificates))
	for _, c := range certificates {
		result = append(result, certificateDocument{
			Name:     c.Name,
			Issuer:   c.Issuer,
			IssuedAt: c.IssuedAt,
		})
	}
	return result
}
