package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"example.com/m/internal/api/auth"
	"example.com/m/internal/storage/models"
	"example.com/m/pkg/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var DB *sqlx.DB

func InitDb(ctx context.Context) {
	dbUri := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s host=%s port=%s", viper.GetString("DP_USERNAME"), viper.GetString("DB_NAME"), viper.GetString("DB_PASS"), viper.GetString("DB_IP"), viper.GetString("DB_PORT"))
	fmt.Println(dbUri)
	db, err := sqlx.Connect("postgres", dbUri)
	if err != nil {
		log.Fatalln(err)
	}

	DB = db

	if err := DB.Ping(); err != nil {
		logger.Fatal(ctx, "failed to connect to database", zap.Error(err))
	} else {
		logger.Info(ctx, "connected to database")
	}
	runMigrations(ctx)
}

func runMigrations(ctx context.Context) {
	migrationsPath := "migrations"

	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		logger.Fatal(ctx, "failed to read migrations folder", zap.Error(err))
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".sql") {
			continue
		}

		path := filepath.Join(migrationsPath, f.Name())
		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			logger.Fatal(ctx, "failed to read migration file", zap.Error(err), zap.String("file", f.Name()))
		}

		query := string(sqlBytes)
		_, err = DB.Exec(query)
		if err != nil {
			logger.Fatal(ctx, "failed to execute migration", zap.Error(err), zap.String("file", f.Name()))
		}

		logger.Info(ctx, fmt.Sprintf("migration %s applied successfully", f.Name()))
	}
}

func CreateUser(user *models.User) error {
	user.CreatedAt = time.Now()

	query := `
		INSERT INTO users (password_hash, username, first_name, created_at, user_role, user_status)
		VALUES (:password_hash, :username, :first_name, :created_at, :user_role, :user_status)
		RETURNING id`

	stmt, err := DB.PrepareNamed(query)
	if err != nil {
		return err
	}

	// stmt.Get will execute and scan the returned id into user.ID
	return stmt.Get(&user.ID, user)
}

// GetUser retrieves a user from the database by ID
func GetUser(id int) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, password_hash, username, first_name, created_at, user_role, user_status
		FROM users 
		WHERE id = $1`

	err := DB.Get(&user, query, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, password_hash, username, first_name, created_at, user_role, user_status
		FROM users 
		WHERE username = $1`
	err := DB.Get(&user, query, username)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func StoreToken(tokenString string) error {
	claims, err := auth.ValidateToken(tokenString, "", false)
	if err != nil {
		return fmt.Errorf("failed to validate token: %w", err)
	}

	hashedToken := auth.HashToken(tokenString)

	token := models.JwtToken{
		UserID:    claims.UserID,
		TokenHash: hashedToken,
		CreatedAt: claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
	}

	query := `
		INSERT INTO jwt_tokens (user_id, token_hash, created_at, expires_at)
		VALUES (:user_id, :token_hash, :created_at, :expires_at)
	`

	_, err = DB.NamedExec(query, token)
	if err != nil {
		return fmt.Errorf("failed to insert token: %w", err)
	}

	return nil
}

func GetJwtToken(tokenString string) (*models.JwtToken, error) {
	hashedToken := auth.HashToken(tokenString)

	var token models.JwtToken
	query := `
	 SELECT user_id, token_hash, active, created_at, expires_at
	 FROM jwt_tokens
	 WHERE token_hash = $1
	`

	err := DB.Get(&token, query, hashedToken)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
