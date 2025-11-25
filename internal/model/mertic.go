package model

import (
	"log"
	"database/sql"
	"time"
)

type Metric struct {
	ID           int    `json:"id"`
	ProjectID    string `json:"projectId"`
	Route        string `json:"route"`
	Method       string `json:"method"`
	StatusCode   int    `json:"statusCode"`
	ResponseTime int64  `json:"responseTime"`
	Timestamp    int64  `json:"timestamp"` // Unix timestamp in milliseconds
	CreatedAt    time.Time `json:"createdAt"`
}

type metricRepository struct {
	db *sql.DB
}
 func NewMetricRepository(db *sql.DB) *metricRepository {
	return &metricRepository{db:db}
}
//create metric table
func (r *metricRepository) CreateMetricTable() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS metrics (
			id SERIAL PRIMARY KEY,
			project_id VARCHAR(255) NOT NULL,
			route VARCHAR(500) NOT NULL,
			method VARCHAR(10) NOT NULL,
			status_code INT NOT NULL,
			response_time BIGINT NOT NULL,
			timestamp BIGINT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_project_id ON metrics(project_id)`,
		`CREATE INDEX IF NOT EXISTS idx_timestamp ON metrics(timestamp)`,
	}

	log.Println("Creating metrics table if not exists...")
	for _, query := range queries {
		_, err := r.db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}
func (r *metricRepository) Save(metric *Metric) error {
	query := `
		INSERT INTO metrics (project_id, route, method, status_code, response_time, timestamp)
		VALUES ($1, $2,$3, $4, $5, $6)	
		returning id ,created_at`
		err:=r.db.QueryRow(query,
		metric.ProjectID,
		metric.Route,
		metric.Method,
		metric.StatusCode,
		metric.ResponseTime,
		metric.Timestamp,
	).Scan(&metric.ID,&metric.CreatedAt	)
	log.Println("Saved metric with ID:", metric.ID)
		return err
}
func( r *metricRepository ) FindAll()([]Metric,error){
	query:=`SELECT id,project_id,route,method,status_code,response_time,timestamp,created_at FROM metrics ORDER BY created_at DESC`
	rows,err:=r.db.Query(query)
	if err!=nil{
		return nil,err
	}
	defer rows.Close()
	var metrics []Metric
	for rows.Next(){
		var m Metric
		err:=rows.Scan(
			&m.ID,
			&m.ProjectID,&m.Route,
			&m.Method,
			&m.StatusCode,&m.ResponseTime,
			&m.Timestamp,
			&m.CreatedAt,
		)
		if err!=nil{
			return nil,err
		}
		metrics=append(metrics, m)
		log.Println("Fetched metric ID:", m.ID);
	}
	return metrics,nil
} 

/// get metrics for a project in last one minute
func (r*metricRepository) DATA_A_min_ago(projectID string)([]Metric,error){
	query:=`SELECT id,project_id,route,method,status_code,response_time,timestamp,created_at FROM metrics WHERE project_id=$1 AND timestamp >= $2 ORDER BY created_at DESC`
	minAgo:=time.Now().Add(-1*time.Minute).UnixMilli()
     rows,err:=r.db.Query(query,projectID,minAgo);
	 if err!=nil{
		log.Println("Error querying metrics:", err)
		return nil,err
	 }
	 defer rows.Close()
	 var metrics []Metric
	 for rows.Next(){
		var m Metric
		err:=rows.Scan(
			&m.ID,
			&m.ProjectID,&m.Route,
			&m.Method,
			&m.StatusCode,&m.ResponseTime,
			&m.Timestamp,
			&m.CreatedAt,
		)
		if err!=nil{
			log.Println("Error scanning metric:", err)
			return nil,err	
		}
		metrics=append(metrics, m)
		log.Println("Fetched metric ID:", m.ID);
	}
	 return metrics,nil
}




	