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
	switch a := asset.(type) {
	case *models.Chart:
		_, err := p.pool.Exec(ctx,
			"INSERT INTO charts (id, title, description, x_axis_title, y_axis_title) VALUES ($1,$2,$3,$4,$5)",
			a.ID, a.Title, a.Description, a.XAxisTitle, a.YAxisTitle,
		)
		if err != nil {
			log.Println(err)
			return
		}
		for _, d := range a.Data {
			_, _ = p.pool.Exec(ctx,
				"INSERT INTO chart_data (chart_id, datapoint_code, value) VALUES ($1,$2,$3)",
				a.ID, d.DatapointCode, d.Value,
			)
		}
	case *models.Insight:
		_, err := p.pool.Exec(ctx,
			"INSERT INTO insights (id, description) VALUES ($1,$2)",
			a.ID, a.Description,
		)
		if err != nil {
			log.Println(err)
			return
		}
	case *models.Audience:
		_, err := p.pool.Exec(ctx,
			"INSERT INTO audiences (id, gender, country, age_group, social_hours, purchases, description) VALUES ($1,$2,$3,$4,$5,$6,$7)",
			a.ID, a.Gender, a.Country, a.AgeGroup, a.SocialHours, a.Purchases, a.Description,
		)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (p *PostgresStore) Get(userID string) []models.Asset {
	ctx := context.Background()
	rows, err := p.pool.Query(ctx, "SELECT asset_id, asset_type FROM favourites WHERE user_id=$1", userID)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	var assets []models.Asset
	for rows.Next() {
		var assetID, assetType string
		if err := rows.Scan(&assetID, &assetType); err != nil {
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
	_, err := p.pool.Exec(ctx,
		`DELETE FROM chart_data WHERE chart_id=$1;
		 DELETE FROM charts WHERE id=$1;
		 DELETE FROM insights WHERE id=$1;
		 DELETE FROM audiences WHERE id=$1;
		 DELETE FROM favourites WHERE asset_id=$1 AND user_id=$2`,
		assetID, userID,
	)
	return err == nil
}

func (p *PostgresStore) EditDescription(userID, assetID, newDesc string) bool {
	ctx := context.Background()

	// Try updating all three tables
	_, err := p.pool.Exec(ctx,
		`UPDATE charts SET description=$1 WHERE id=$2;
		 UPDATE insights SET description=$1 WHERE id=$2;
		 UPDATE audiences SET description=$1 WHERE id=$2`,
		newDesc, assetID,
	)
	return err == nil
}

// ----------------- Favourite Methods -----------------

func (p *PostgresStore) AddFavourite(userID, assetID, assetType string) bool {
	_, err := p.pool.Exec(context.Background(),
		"INSERT INTO favourites (user_id, asset_id, asset_type) VALUES ($1,$2,$3) ON CONFLICT DO NOTHING",
		userID, assetID, assetType,
	)
	return err == nil
}

func (p *PostgresStore) RemoveFavourite(userID, assetID string) bool {
	_, err := p.pool.Exec(context.Background(),
		"DELETE FROM favourites WHERE user_id=$1 AND asset_id=$2",
		userID, assetID,
	)
	return err == nil
}

func (p *PostgresStore) GetFavourites(userID string) []models.Favourite {
	ctx := context.Background()
	rows, err := p.pool.Query(ctx,
		"SELECT asset_id, asset_type FROM favourites WHERE user_id=$1", userID)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	var favs []models.Favourite
	for rows.Next() {
		var assetID, assetType string
		if err := rows.Scan(&assetID, &assetType); err != nil {
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
