package storage

import (
	"context"
	"log"

	"assetsApp/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

// ----------------- Asset Methods -----------------
func (p *PostgresStore) Add(userID string, asset models.Asset) {
	ctx := context.Background()

	// Ensure user exists
	_, err := p.pool.Exec(ctx,
		"INSERT INTO users (id, name) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		userID, "User "+userID, //TODO FIX TABLE SCHEMA TO REMOVE NAME CONSTRAINT
	)
	if err != nil {
		log.Println("Failed to ensure user exists:", err)
		return
	}

	switch a := asset.(type) {
	case *models.Chart:
		tx, err := p.pool.Begin(ctx)
		if err != nil {
			log.Println("Failed to start transaction:", err)
			return
		}
		defer tx.Rollback(ctx)

		// 1. Insert parent first
		_, err = tx.Exec(ctx,
			"INSERT INTO assets (asset_id, title, description, asset_type, user_id) VALUES ($1, $2, $3, $4, $5)",
			a.ID, a.Title, a.Description, "chart", userID,
		)
		if err != nil {
			log.Println("Failed to insert into assets (chart):", err)
			return
		}

		// 2. Insert child now that parent exists
		_, err = tx.Exec(ctx,
			"INSERT INTO charts (id, title, description, x_axis_title, y_axis_title) VALUES ($1,$2,$3,$4,$5)",
			a.ID, a.Title, a.Description, a.XAxisTitle, a.YAxisTitle,
		)
		if err != nil {
			log.Println("Failed to insert chart:", err)
			return
		}

		for _, d := range a.Data {
			_, err = tx.Exec(ctx,
				"INSERT INTO chart_data (chart_id, datapoint_code, value) VALUES ($1,$2,$3)",
				a.ID, d.DatapointCode, d.Value,
			)
			if err != nil {
				log.Println("Failed to insert chart data:", err)
				return
			}
		}

		if err = tx.Commit(ctx); err != nil {
			log.Println("Failed to commit chart transaction:", err)
		}

	case *models.Insight:
		tx, err := p.pool.Begin(ctx)
		if err != nil {
			log.Println("Failed to start transaction:", err)
			return
		}
		defer tx.Rollback(ctx)

		// Insert into assets first (and include userID)
		_, err = tx.Exec(ctx,
			"INSERT INTO assets (asset_id, title, description, asset_type, user_id) VALUES ($1, $2, $3, $4, $5)",
			a.ID, "Insight", a.Description, "insight", userID,
		)
		if err != nil {
			log.Println("Failed to insert into assets:", err)
			return
		}

		_, err = tx.Exec(ctx,
			"INSERT INTO insights (id, description) VALUES ($1,$2)",
			a.ID, a.Description,
		)
		if err != nil {
			log.Println("Failed to insert insight:", err)
			return
		}

		if err = tx.Commit(ctx); err != nil {
			log.Println("Failed to commit insight transaction:", err)
		}

	case *models.Audience:
		tx, err := p.pool.Begin(ctx)
		if err != nil {
			log.Println("Failed to start transaction:", err)
			return
		}
		defer tx.Rollback(ctx)

		// Insert into assets first
		_, err = tx.Exec(ctx,
			"INSERT INTO assets (asset_id, title, description, asset_type, user_id) VALUES ($1, $2, $3, $4, $5)",
			a.ID, "Audience", a.Description, "audience", userID,
		)
		if err != nil {
			log.Println("Failed to insert into assets:", err)
			return
		}

		_, err = tx.Exec(ctx,
			"INSERT INTO audiences (id, gender, country, age_group, social_hours, purchases, description) VALUES ($1,$2,$3,$4,$5,$6,$7)",
			a.ID, a.Gender, a.Country, a.AgeGroup, a.SocialHours, a.Purchases, a.Description,
		)
		if err != nil {
			log.Println("Failed to insert audience:", err)
			return
		}

		if err = tx.Commit(ctx); err != nil {
			log.Println("Failed to commit audience transaction:", err)
		}
	}
}

func (p *PostgresStore) Get(userID string) []models.Asset {
	ctx := context.Background()
	rows, err := p.pool.Query(ctx, "SELECT asset_id, asset_type FROM assets WHERE user_id=$1", userID)
	if err != nil {
		log.Println("Failed to get assets:", err)
		return nil
	}
	defer rows.Close()

	var assets []models.Asset
	for rows.Next() {
		var assetID, assetType string
		if err := rows.Scan(&assetID, &assetType); err != nil {
			log.Println("Failed to scan asset row:", err)
			continue
		}

		var asset models.Asset
		switch assetType {
		case "chart":
			asset = &models.Chart{ID: assetID}
		case "insight":
			asset = &models.Insight{ID: assetID}
		case "audience":
			asset = &models.Audience{ID: assetID}
		}
		assets = append(assets, asset)
	}
	return assets
}

func (p *PostgresStore) Remove(userID, assetID string) bool {
	ctx := context.Background()
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		log.Println("Failed to start remove transaction:", err)
		return false
	}
	defer tx.Rollback(ctx)

	statements := []struct {
		query string
		args  []interface{}
	}{
		{"DELETE FROM chart_data WHERE chart_id=$1", []interface{}{assetID}},
		{"DELETE FROM charts WHERE id=$1", []interface{}{assetID}},
		{"DELETE FROM insights WHERE id=$1", []interface{}{assetID}},
		{"DELETE FROM audiences WHERE id=$1", []interface{}{assetID}},
		{"DELETE FROM assets WHERE asset_id=$1", []interface{}{assetID}},
		{"DELETE FROM favourites WHERE asset_id=$1 AND user_id=$2", []interface{}{assetID, userID}},
	}

	for _, stmt := range statements {
		if _, err := tx.Exec(ctx, stmt.query, stmt.args...); err != nil {
			log.Println("Failed to execute remove statement:", err)
			return false
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("Failed to commit remove transaction:", err)
		return false
	}

	return true
}

func (p *PostgresStore) EditDescription(userID, assetID, newDesc string) bool {
	ctx := context.Background()
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		log.Println("Failed to start edit transaction:", err)
		return false
	}
	defer tx.Rollback(ctx)

	statements := []string{
		"UPDATE charts SET description=$1 WHERE id=$2",
		"UPDATE insights SET description=$1 WHERE id=$2",
		"UPDATE audiences SET description=$1 WHERE id=$2",
		"UPDATE assets SET description=$1 WHERE asset_id=$2",
	}

	for _, stmt := range statements {
		if _, err := tx.Exec(ctx, stmt, newDesc, assetID); err != nil {
			log.Println("Failed to update description:", err)
			return false
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("Failed to commit edit transaction:", err)
		return false
	}

	return true
}

// ----------------- Favourite Methods -----------------

func (p *PostgresStore) AddFavourite(userID, assetID, assetType string) bool {
	ctx := context.Background()

	// Ensure user exists
	_, err := p.pool.Exec(ctx,
		"INSERT INTO users (id, name) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING",
		userID, "Unknown",
	)
	if err != nil {
		log.Println("Failed to ensure user exists:", err)
		return false
	}

	// Check that the asset belongs to this user
	var ownerID string
	err = p.pool.QueryRow(ctx, "SELECT user_id FROM assets WHERE asset_id=$1", assetID).Scan(&ownerID)
	if err != nil {
		log.Println("Failed to fetch asset owner or asset does not exist:", err)
		return false
	}
	if ownerID != userID {
		log.Println("Cannot favourite an asset not owned by the user")
		return false
	}

	// Insert into favourites using composite PK
	_, err = p.pool.Exec(ctx,
		"INSERT INTO favourites (user_id, asset_id, asset_type) VALUES ($1, $2, $3) ON CONFLICT (user_id, asset_id) DO NOTHING",
		userID, assetID, assetType,
	)
	if err != nil {
		log.Println("Failed to add favourite:", err)
		return false
	}

	return true
}

func (p *PostgresStore) RemoveFavourite(userID, assetID string) bool {
	_, err := p.pool.Exec(context.Background(),
		"DELETE FROM favourites WHERE user_id=$1 AND asset_id=$2",
		userID, assetID,
	)
	if err != nil {
		log.Println("Failed to remove favourite:", err)
		return false
	}
	return true
}

func (p *PostgresStore) GetFavourites(userID string) []models.Favourite {
	ctx := context.Background()
	rows, err := p.pool.Query(ctx,
		"SELECT asset_id, asset_type FROM favourites WHERE user_id=$1", userID)
	if err != nil {
		log.Println("Failed to get favourites:", err)
		return nil
	}
	defer rows.Close()

	var favs []models.Favourite
	for rows.Next() {
		var assetID, assetType string
		if err := rows.Scan(&assetID, &assetType); err != nil {
			log.Println("Failed to scan favourite row:", err)
			continue
		}

		var asset models.Asset
		switch assetType {
		case "chart":
			asset = &models.Chart{ID: assetID}
		case "insight":
			asset = &models.Insight{ID: assetID}
		case "audience":
			asset = &models.Audience{ID: assetID}
		}

		favs = append(favs, models.Favourite{
			UserID: userID,
			Asset:  asset,
		})
	}
	return favs
}
