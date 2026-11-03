package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"latihan-fiber/app/model"
	"latihan-fiber/config"
	"latihan-fiber/database"
	"latihan-fiber/helper"
)

type testUser struct {
	ID       int
	Username string
	Role     string
}

type result struct {
	Username   string
	StatusCode int
	Body       string
	Err        error
}

func main() {
	mode := flag.String("mode", "crowd", "crowd")
	baseURL := flag.String("url", "", "Base URL API")
	flag.Parse()

	if *mode != "crowd" {
		fmt.Println("Mode yang didukung saat ini: crowd")
		return
	}

	config.LoadEnv()

	if *baseURL == "" {
		*baseURL = "http://localhost:" +
			config.GetEnv("APP_PORT", "3000")
	}

	secret := config.GetEnv("JWT_SECRET", "")
	if len(secret) < 32 {
		fmt.Println("ERROR: JWT_SECRET kosong atau kurang dari 32 karakter")
		return
	}

	issuer := config.GetEnv("JWT_ISSUER", "praktikum-backend")
	ttl := time.Duration(
		config.GetEnvInt("JWT_ACCESS_TTL_MINUTES", 15),
	) * time.Minute

	jwtManager := helper.NewJWTManager(secret, issuer, ttl)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	pool, err := database.NewPool(ctx)
	if err != nil {
		fmt.Println("Gagal terhubung ke database:", err)
		return
	}
	defer pool.Close()

	var course model.Course
	err = pool.QueryRow(ctx, `
		SELECT id, code, name, credits, quota,
		       enrolled_count, version, created_at
		FROM courses
		WHERE code = $1
	`, "IF399").Scan(
		&course.ID,
		&course.Code,
		&course.Name,
		&course.Credits,
		&course.Quota,
		&course.EnrolledCount,
		&course.Version,
		&course.CreatedAt,
	)
	if err != nil {
		fmt.Println("Gagal mencari IF399:", err)
		return
	}

	users := make([]testUser, 0, 60)

	for i := 1; i <= 60; i++ {
		username := fmt.Sprintf("mhs%03d", i)

		var u testUser
		err := pool.QueryRow(ctx, `
			SELECT id, username, role
			FROM users
			WHERE username = $1 AND is_active = TRUE
		`, username).Scan(&u.ID, &u.Username, &u.Role)

		if err != nil {
			fmt.Printf("Gagal mengambil akun %s: %v\n", username, err)
			return
		}

		users = append(users, u)
	}

	var before int
	err = pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM enrollments
		WHERE course_id = $1
	`, course.ID).Scan(&before)
	if err != nil {
		fmt.Println("Gagal menghitung enrollment awal:", err)
		return
	}

	fmt.Println("Mode: crowd")
	fmt.Println("Base URL:", *baseURL)
	fmt.Printf("Mata kuliah: %s (ID %d)\n", course.Code, course.ID)
	fmt.Printf("Kuota: %d\n", course.Quota)
	fmt.Printf("Enrollment sebelum tes: %d\n", before)
	fmt.Println("Mengirim 60 request bersamaan...")

	start := make(chan struct{})
	results := make(chan result, len(users))

	var wg sync.WaitGroup

	client := &http.Client{Timeout: 20 * time.Second}
	enrollURL := fmt.Sprintf(
		"%s/api/v1/courses/%d/enroll",
		*baseURL, course.ID,
	)

	for _, user := range users {
		wg.Add(1)

		go func(u testUser) {
			defer wg.Done()
			<-start

			accessToken, err := jwtManager.GenerateAccess(model.User{
				ID:       u.ID,
				Username: u.Username,
				Role:     u.Role,
			})
			if err != nil {
				results <- result{Username: u.Username, Err: err}
				return
			}

			fmt.Printf("Username: %s\nToken: %s\n\n", u.Username, accessToken)

			req, err := http.NewRequest(
				http.MethodPost, enrollURL, bytes.NewReader([]byte("{}")),
			)
			if err != nil {
				results <- result{Username: u.Username, Err: err}
				return
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+accessToken)

			resp, err := client.Do(req)
			if err != nil {
				results <- result{Username: u.Username, Err: err}
				return
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)

			results <- result{
				Username:   u.Username,
				StatusCode: resp.StatusCode,
				Body:       string(body),
			}
		}(user)
	}

	close(start)
	wg.Wait()
	close(results)

	counts := make(map[int]int)
	errorsCount := 0

	for r := range results {
		if r.Err != nil {
			errorsCount++
			fmt.Printf("[ERROR] %s: %v\n", r.Username, r.Err)
			continue
		}

		counts[r.StatusCode]++
		fmt.Printf("[%d] %s: %s\n",
			r.StatusCode, r.Username, r.Body)
	}

	var after int
	err = pool.QueryRow(context.Background(), `
		SELECT COUNT(*)
		FROM enrollments
		WHERE course_id = $1
	`, course.ID).Scan(&after)
	if err != nil {
		fmt.Println("Gagal menghitung enrollment akhir:", err)
		return
	}

	fmt.Println("\n========== HASIL CROWD ==========")
	fmt.Println("Total request:", len(users))
	fmt.Println("HTTP 201 (berhasil):", counts[201])
	fmt.Println("HTTP 409 (conflict):", counts[409])
	fmt.Println("HTTP lainnya:")
	for code, count := range counts {
		if code != 201 && code != 409 {
			fmt.Printf("  HTTP %d: %d\n", code, count)
		}
	}
	fmt.Println("Request error:", errorsCount)
	fmt.Println("Enrollment sebelum:", before)
	fmt.Println("Enrollment sesudah:", after)
	fmt.Println("Kuota:", course.Quota)

	if after > course.Quota {
		fmt.Println("KESIMPULAN: KUOTA TERLAMPAUI — indikasi race condition")
	} else {
		fmt.Println("KESIMPULAN: Kuota belum terlampaui pada eksekusi ini")
	}

	_ = json.Valid
}
