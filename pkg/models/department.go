package models

import (
	"context"
	"time"
)

type Department struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func (m *Model) GetDepartmentById(ctx context.Context, departmentID string) (Department, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var department Department
	row := m.Conn.QueryRow(ctx, "select id, title FROM departments WHERE id = $1", departmentID)
	err := row.Scan(&department.ID, &department.Title)
	if err != nil {
		return Department{}, err
	}

	return department, nil
}

func (m *Model) GetDepartments(ctx context.Context) ([]Department, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	rows, err := m.Conn.Query(ctx, "select id, title FROM departments")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departments []Department
	for rows.Next() {
		var department Department
		err = rows.Scan(&department.ID, &department.Title)
		if err != nil {
			return nil, err
		}

		departments = append(departments, department)
	}

	return departments, nil
}
