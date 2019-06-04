package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var createTableStmts = []string{
	`CREATE TABLE IF NOT EXISTS "public"."enriched_profiles" (
	"id" serial,
	"CustomerID" text NOT NULL,
	"Partner" text,
	"Dependents" text,
	"Tenure" smallint,
	"PhoneService" text,
	"MultipleLines" text,
	"InternetService" text,
	"OnlineSecurity" text,
	"OnlineBackup" text,
	"DeviceProtection" text,
	"TechSupport" text,
	"StreamingTV" text,
	"StreamingMovies" text,
	"Contract" text,
	"PaperlessBilling" text,
	"PaymentMethod" text,
	"MonthlyCharges" decimal,
	"TotalCharges" decimal,
	"ChurnScore" int,
	"EnrichedAt" timestamp,
	PRIMARY KEY ("id")
)`,
}

var createIndexStmts = []string{
	`CREATE INDEX "customer_id_idx" ON "public"."enriched_profile"("CustomerID")`,
}

// createTables creates the tables
func createTables(conn *sql.DB) error {
	for _, stmt := range createTableStmts {
		_, err := conn.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

// createIndexes creates any indexes
func createIndexes(conn *sql.DB) error {
	for _, stmt := range createIndexStmts {
		_, err := conn.Exec(stmt)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

// PrepareDB creates our tables if they don't exist , and optionally attempts to setup indexes
func PrepareDB(conn *sql.DB, setupIndexes bool) error {
	if err := createTables(db); err != nil {
		return err
	}
	if setupIndexes {
		if err := createIndexes(db); err != nil {
			return err
		}
	}
	return nil
}
