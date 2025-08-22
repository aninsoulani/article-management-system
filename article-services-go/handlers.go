package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Helper untuk error response
func respondError(w http.ResponseWriter, code int, msg string, details ...any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := map[string]any{
		"error":   true,
		"message": msg,
	}
	if len(details) > 0 {
		resp["details"] = details
	}

	json.NewEncoder(w).Encode(resp)
}

// Helper untuk sukses response
func respondJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// Fungsi untuk parsing error validasi biar readable
func parseValidationError(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		var msgs []string
		for _, e := range errs {
			field := strings.ToLower(e.Field())
			switch e.Tag() {
			case "required":
				msgs = append(msgs, fmt.Sprintf("Field '%s' wajib diisi", field))
			case "oneof":
				msgs = append(msgs, fmt.Sprintf("Field '%s' harus salah satu dari: %s", field, e.Param()))
			case "min":
				msgs = append(msgs, fmt.Sprintf("Field '%s' minimal %s karakter", field, e.Param()))
			default:
				msgs = append(msgs, fmt.Sprintf("Field '%s' tidak valid", field))
			}
		}
		return strings.Join(msgs, ", ")
	}
	return err.Error()
}


// ================= HANDLERS ===================

// create article (POST /article/)
func CreatePost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Tidak bisa membaca body request")
		return
	}

	body = bytes.TrimSpace(body)
	if len(body) == 0 {
		respondError(w, http.StatusBadRequest, "Body request kosong")
		return
	}

	if body[0] == '[' {
		// bulk insert
		var posts []Post
		if err := json.Unmarshal(body, &posts); err != nil {
			respondError(w, http.StatusBadRequest, "Format JSON tidak valid")
			return
		}

		for _, p := range posts {
			if err := validate.Struct(p); err != nil {
				respondError(w, http.StatusBadRequest, parseValidationError(err))
				return
			}
		}

		if err := DB.Create(&posts).Error; err != nil {
			respondError(w, http.StatusInternalServerError, "Gagal menyimpan artikel")
			return
		}

		respondJSON(w, http.StatusCreated, map[string]any{
			"message": fmt.Sprintf("%d artikel berhasil dibuat", len(posts)),
		})
		return
	}

	// single insert
	var post Post
	if err := json.Unmarshal(body, &post); err != nil {
		respondError(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	if err := validate.Struct(post); err != nil {
		respondError(w, http.StatusBadRequest, parseValidationError(err))
		return
	}

	if err := DB.Create(&post).Error; err != nil {
		respondError(w, http.StatusInternalServerError, "Gagal menyimpan artikel")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]string{"message": "Artikel berhasil dibuat"})
}

// get semua article dengan paging (GET /article/{limit}/{offset})
func GetPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	limit, _ := strconv.Atoi(vars["limit"])
	offset, _ := strconv.Atoi(vars["offset"])

	var posts []Post
	var total int64

	if err := DB.Model(&Post{}).Count(&total).Error; err != nil {
		respondError(w, http.StatusInternalServerError, "Gagal menghitung artikel")
		return
	}

	if err := DB.Limit(limit).Offset(offset).Find(&posts).Error; err != nil {
		respondError(w, http.StatusInternalServerError, "Gagal mengambil data artikel")
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"data":  posts,
		"total": total,
	})
}

// get semua article by status dengan paging (GET /article/status/{status}/{limit}/{offset})
func GetPostsByStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status := vars["status"]
	limit, _ := strconv.Atoi(vars["limit"])
	offset, _ := strconv.Atoi(vars["offset"])

	var posts []Post
	var total int64

	if err := DB.Model(&Post{}).Where("status = ?", status).Count(&total).Error; err != nil {
		respondError(w, http.StatusInternalServerError, "Gagal menghitung artikel berdasarkan status")
		return
	}

	if err := DB.Where("status = ?", status).Limit(limit).Offset(offset).Find(&posts).Error; err != nil {
		respondError(w, http.StatusInternalServerError, "Gagal mengambil artikel berdasarkan status")
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"data":  posts,
		"total": total,
	})
}


// get 1 article (GET /article/{id})
func GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var post Post
	if err := DB.First(&post, id).Error; err == gorm.ErrRecordNotFound {
		respondError(w, http.StatusNotFound, "Artikel tidak ditemukan")
		return
	} else if err != nil {
		respondError(w, http.StatusInternalServerError, "Gagal mengambil artikel")
		return
	}

	respondJSON(w, http.StatusOK, post)
}

// update article (PUT /article/{id})
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var post Post
	if err := DB.First(&post, id).Error; err == gorm.ErrRecordNotFound {
		respondError(w, http.StatusNotFound, "Artikel tidak ditemukan")
		return
	}

	var updateData Post
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		respondError(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	if err := validate.Struct(updateData); err != nil {
		respondError(w, http.StatusBadRequest, parseValidationError(err))
		return
	}

	post.Title = updateData.Title
	post.Content = updateData.Content
	post.Category = updateData.Category
	post.Status = updateData.Status

	if err := DB.Save(&post).Error; err != nil {
		respondError(w, http.StatusInternalServerError, "Gagal memperbarui artikel")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Artikel berhasil diperbarui"})
}

// update only status (PATCH /article/{id}/status)
func UpdatePostStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var post Post
	if err := DB.First(&post, id).Error; err == gorm.ErrRecordNotFound {
		respondError(w, http.StatusNotFound, "Artikel tidak ditemukan")
		return
	}

	var body struct {
		Status string `json:"status" validate:"required,oneof=publish draft trash"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}
	if err := validate.Struct(body); err != nil {
		respondError(w, http.StatusBadRequest, parseValidationError(err))
		return
	}

	post.Status = body.Status
	if err := DB.Save(&post).Error; err != nil {
		respondError(w, http.StatusInternalServerError, "Gagal memperbarui status artikel")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Status artikel berhasil diperbarui"})
}

// delete article permanently (DELETE /article/{id})
func DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var post Post
	if err := DB.First(&post, id).Error; err == gorm.ErrRecordNotFound {
		respondError(w, http.StatusNotFound, "Artikel tidak ditemukan")
		return
	}

	if err := DB.Delete(&post).Error; err != nil {
		respondError(w, http.StatusInternalServerError, "Gagal menghapus artikel")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Artikel berhasil dihapus"})
}
