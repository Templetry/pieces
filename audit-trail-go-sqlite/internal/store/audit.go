package store

import (
	"database/sql"
	"time"
)

// init registers the audit table through the migrations socket, so no
// existing file changes (ADR-0014).
func init() {
	Register(`
CREATE TABLE IF NOT EXISTS audit_entry (
	id          INTEGER PRIMARY KEY AUTOINCREMENT,
	at          TEXT NOT NULL DEFAULT (datetime('now')),
	actor       TEXT NOT NULL DEFAULT '',
	action      TEXT NOT NULL,
	object_name TEXT NOT NULL,
	object_id   TEXT NOT NULL DEFAULT '',
	detail      TEXT NOT NULL DEFAULT '',
	ip          TEXT NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS audit_entry_object ON audit_entry (object_name, object_id);`)
}

// AuditEntry is one immutable record of a change.
type AuditEntry struct {
	ID         int64  `json:"id"`
	At         string `json:"at"`
	Actor      string `json:"actor"`
	Action     string `json:"action"`
	ObjectName string `json:"objectName"`
	ObjectID   string `json:"objectId"`
	Detail     string `json:"detail"`
	IP         string `json:"ip"`
}

// Audit is the append-only repository. It exposes no update or delete:
// a trail the application can rewrite is not a trail.
type Audit struct{ DB *sql.DB }

// Record appends one entry.
func (a Audit) Record(e AuditEntry) (AuditEntry, error) {
	if e.At == "" {
		e.At = time.Now().UTC().Format(time.RFC3339)
	}
	res, err := a.DB.Exec(
		`INSERT INTO audit_entry (at, actor, action, object_name, object_id, detail, ip)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		e.At, e.Actor, e.Action, e.ObjectName, e.ObjectID, e.Detail, e.IP)
	if err != nil {
		return AuditEntry{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return AuditEntry{}, err
	}
	e.ID = id
	return e, nil
}

// List returns the most recent entries, optionally filtered by object type.
func (a Audit) List(objectName string, limit, offset int) ([]AuditEntry, error) {
	query := `SELECT id, at, actor, action, object_name, object_id, detail, ip
	          FROM audit_entry`
	args := []any{}
	if objectName != "" {
		query += ` WHERE object_name = ?`
		args = append(args, objectName)
	}
	query += ` ORDER BY id DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := a.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []AuditEntry{}
	for rows.Next() {
		var e AuditEntry
		if err := rows.Scan(&e.ID, &e.At, &e.Actor, &e.Action,
			&e.ObjectName, &e.ObjectID, &e.Detail, &e.IP); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
