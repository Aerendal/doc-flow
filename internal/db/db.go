package db

import (
        "database/sql"
        "fmt"
        "os"
        "path/filepath"
        "strings"
        "time"

        "docflow/internal/model"

        _ "modernc.org/sqlite"
)

type Store struct {
        db *sql.DB
}

func New(dbPath string) (*Store, error) {
        dir := filepath.Dir(dbPath)
        if dir != "." && dir != "" {
                if err := os.MkdirAll(dir, 0755); err != nil {
                        return nil, fmt.Errorf("cannot create directory %s: %w", dir, err)
                }
        }

        db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(wal)&_pragma=foreign_keys(on)")
        if err != nil {
                return nil, fmt.Errorf("cannot open database: %w", err)
        }

        if err := db.Ping(); err != nil {
                return nil, fmt.Errorf("cannot connect to database: %w", err)
        }

        s := &Store{db: db}
        if err := s.migrate(); err != nil {
                return nil, fmt.Errorf("migration failed: %w", err)
        }

        return s, nil
}

func (s *Store) Close() error {
        return s.db.Close()
}

func (s *Store) migrate() error {
        schema := `
        CREATE TABLE IF NOT EXISTS project (
                id INTEGER PRIMARY KEY,
                name TEXT NOT NULL DEFAULT 'docflow',
                created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );

        CREATE TABLE IF NOT EXISTS phase (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                day INTEGER NOT NULL UNIQUE,
                name TEXT NOT NULL,
                description TEXT NOT NULL DEFAULT '',
                status TEXT NOT NULL DEFAULT 'todo',
                created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );

        CREATE TABLE IF NOT EXISTS phase_dependency (
                phase_id INTEGER NOT NULL,
                depends_on_id INTEGER NOT NULL,
                PRIMARY KEY (phase_id, depends_on_id),
                FOREIGN KEY (phase_id) REFERENCES phase(id) ON DELETE CASCADE,
                FOREIGN KEY (depends_on_id) REFERENCES phase(id) ON DELETE CASCADE
        );

        CREATE TABLE IF NOT EXISTS task (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                phase_id INTEGER NOT NULL,
                name TEXT NOT NULL,
                done INTEGER NOT NULL DEFAULT 0,
                created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                FOREIGN KEY (phase_id) REFERENCES phase(id) ON DELETE CASCADE
        );
        `
        _, err := s.db.Exec(schema)
        return err
}

func (s *Store) InitProject(name string) error {
        _, err := s.db.Exec(`INSERT OR IGNORE INTO project (id, name) VALUES (1, ?)`, name)
        return err
}

func (s *Store) GetProjectName() (string, error) {
        var name string
        err := s.db.QueryRow(`SELECT name FROM project WHERE id = 1`).Scan(&name)
        if err == sql.ErrNoRows {
                return "", fmt.Errorf("project not initialized – run 'docflow init' first")
        }
        return name, err
}

func (s *Store) AddPhase(p model.Phase) (int64, error) {
        res, err := s.db.Exec(
                `INSERT INTO phase (day, name, description, status) VALUES (?, ?, ?, ?)`,
                p.Day, p.Name, p.Description, p.Status,
        )
        if err != nil {
                return 0, err
        }

        id, err := res.LastInsertId()
        if err != nil {
                return 0, err
        }

        return id, nil
}

func (s *Store) LinkDependencies(phaseDay int, depDays []int) (missing []int, err error) {
        var phaseID int
        err = s.db.QueryRow(`SELECT id FROM phase WHERE day = ?`, phaseDay).Scan(&phaseID)
        if err != nil {
                return nil, fmt.Errorf("phase day=%d not found: %w", phaseDay, err)
        }

        for _, depDay := range depDays {
                var depID int
                err := s.db.QueryRow(`SELECT id FROM phase WHERE day = ?`, depDay).Scan(&depID)
                if err != nil {
                        missing = append(missing, depDay)
                        continue
                }
                _, err = s.db.Exec(`INSERT OR IGNORE INTO phase_dependency (phase_id, depends_on_id) VALUES (?, ?)`, phaseID, depID)
                if err != nil {
                        return missing, err
                }
        }

        return missing, nil
}

func (s *Store) ListPhases(statusFilter string) ([]model.Phase, error) {
        query := `SELECT id, day, name, description, status, created_at, updated_at FROM phase`
        var args []interface{}

        if statusFilter != "" {
                query += ` WHERE status = ?`
                args = append(args, statusFilter)
        }
        query += ` ORDER BY day ASC`

        rows, err := s.db.Query(query, args...)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var phases []model.Phase
        for rows.Next() {
                var p model.Phase
                var status string
                err := rows.Scan(&p.ID, &p.Day, &p.Name, &p.Description, &status, &p.CreatedAt, &p.UpdatedAt)
                if err != nil {
                        return nil, err
                }
                p.Status = model.PhaseStatus(status)

                deps, err := s.getPhaseDependencies(p.ID)
                if err != nil {
                        return nil, err
                }
                p.DependsOn = deps
                phases = append(phases, p)
        }
        return phases, nil
}

func (s *Store) getPhaseDependencies(phaseID int) ([]int, error) {
        rows, err := s.db.Query(
                `SELECT p.day FROM phase_dependency pd JOIN phase p ON p.id = pd.depends_on_id WHERE pd.phase_id = ?`,
                phaseID,
        )
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var deps []int
        for rows.Next() {
                var day int
                if err := rows.Scan(&day); err != nil {
                        return nil, err
                }
                deps = append(deps, day)
        }
        return deps, nil
}

func (s *Store) UpdatePhaseStatus(day int, status model.PhaseStatus) error {
        valid := map[model.PhaseStatus]bool{
                model.StatusTodo:       true,
                model.StatusInProgress: true,
                model.StatusDone:       true,
                model.StatusBlocked:    true,
        }
        if !valid[status] {
                return fmt.Errorf("invalid status: %s (valid: %s)", status,
                        strings.Join([]string{string(model.StatusTodo), string(model.StatusInProgress), string(model.StatusDone), string(model.StatusBlocked)}, ", "))
        }

        res, err := s.db.Exec(
                `UPDATE phase SET status = ?, updated_at = ? WHERE day = ?`,
                status, time.Now(), day,
        )
        if err != nil {
                return err
        }
        n, _ := res.RowsAffected()
        if n == 0 {
                return fmt.Errorf("phase with day %d not found", day)
        }
        return nil
}

func (s *Store) GetPhaseByDay(day int) (*model.Phase, error) {
        var p model.Phase
        var status string
        err := s.db.QueryRow(
                `SELECT id, day, name, description, status, created_at, updated_at FROM phase WHERE day = ?`,
                day,
        ).Scan(&p.ID, &p.Day, &p.Name, &p.Description, &status, &p.CreatedAt, &p.UpdatedAt)
        if err == sql.ErrNoRows {
                return nil, fmt.Errorf("phase with day %d not found", day)
        }
        if err != nil {
                return nil, err
        }
        p.Status = model.PhaseStatus(status)

        deps, err := s.getPhaseDependencies(p.ID)
        if err != nil {
                return nil, err
        }
        p.DependsOn = deps
        return &p, nil
}

func (s *Store) AddTask(phaseDay int, name string) error {
        var phaseID int
        err := s.db.QueryRow(`SELECT id FROM phase WHERE day = ?`, phaseDay).Scan(&phaseID)
        if err == sql.ErrNoRows {
                return fmt.Errorf("phase with day %d not found", phaseDay)
        }
        if err != nil {
                return err
        }

        _, err = s.db.Exec(`INSERT INTO task (phase_id, name) VALUES (?, ?)`, phaseID, name)
        return err
}

func (s *Store) CompleteTask(taskID int) error {
        res, err := s.db.Exec(`UPDATE task SET done = 1, updated_at = ? WHERE id = ?`, time.Now(), taskID)
        if err != nil {
                return err
        }
        n, _ := res.RowsAffected()
        if n == 0 {
                return fmt.Errorf("task with id %d not found", taskID)
        }
        return nil
}

func (s *Store) ListTasks(phaseDay int) ([]model.Task, error) {
        var phaseID int
        err := s.db.QueryRow(`SELECT id FROM phase WHERE day = ?`, phaseDay).Scan(&phaseID)
        if err == sql.ErrNoRows {
                return nil, fmt.Errorf("phase with day %d not found", phaseDay)
        }
        if err != nil {
                return nil, err
        }

        rows, err := s.db.Query(
                `SELECT id, phase_id, name, done, created_at, updated_at FROM task WHERE phase_id = ? ORDER BY id`,
                phaseID,
        )
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var tasks []model.Task
        for rows.Next() {
                var t model.Task
                var done int
                err := rows.Scan(&t.ID, &t.PhaseID, &t.Name, &done, &t.CreatedAt, &t.UpdatedAt)
                if err != nil {
                        return nil, err
                }
                t.Done = done == 1
                tasks = append(tasks, t)
        }
        return tasks, nil
}

func (s *Store) GetStats() (*model.ProjectStats, error) {
        stats := &model.ProjectStats{}

        rows, err := s.db.Query(`SELECT status, COUNT(*) FROM phase GROUP BY status`)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        for rows.Next() {
                var status string
                var count int
                if err := rows.Scan(&status, &count); err != nil {
                        return nil, err
                }
                stats.TotalPhases += count
                switch model.PhaseStatus(status) {
                case model.StatusTodo:
                        stats.TodoPhases = count
                case model.StatusInProgress:
                        stats.InProgressPhases = count
                case model.StatusDone:
                        stats.DonePhases = count
                case model.StatusBlocked:
                        stats.BlockedPhases = count
                }
        }

        err = s.db.QueryRow(`SELECT COUNT(*) FROM task`).Scan(&stats.TotalTasks)
        if err != nil {
                return nil, err
        }
        err = s.db.QueryRow(`SELECT COUNT(*) FROM task WHERE done = 1`).Scan(&stats.DoneTasks)
        if err != nil {
                return nil, err
        }

        if stats.TotalPhases > 0 {
                stats.PercentComplete = float64(stats.DonePhases) / float64(stats.TotalPhases) * 100
        }

        return stats, nil
}

func (s *Store) ClearPhases() error {
        _, err := s.db.Exec(`DELETE FROM phase`)
        return err
}
