package obs

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db         *sql.DB
	obstype    string
	parameters []Parameter
}

type Parameter struct {
	Parameter  string
	AccumHours float32
	Units      string
}

type Data struct {
	ValidDate time.Time
	SID       int
	Lat       float32
	Lon       float32
	Elev      int
	Data      map[string]float32
}

func open(filename, obstype string) (*Database, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	return &Database{
		db:      db,
		obstype: obstype,
	}, nil
}

func Create(filename, obstype string, parameters []Parameter) (*Database, error) {
	db, err := open(filename, obstype)
	if err != nil {
		return nil, err
	}

	if err := db.addObstype(obstype, parameters); err != nil {
		db.Close()
		return nil, err
	}

	db.parameters = parameters

	return db, nil
}

func (db *Database) Close() error {
	return nil
}

func (db *Database) Add(data ...Data) error {
	var baseStatement strings.Builder
	fmt.Fprintf(&baseStatement, "INSERT INTO %s (validdate, SID, lat, lon, elev", strings.ToUpper(db.obstype))
	for _, p := range db.parameters {
		fmt.Fprintf(&baseStatement, ", %s", p.Parameter)
	}
	fmt.Fprint(&baseStatement, ") VALUES (?, ?, ?, ?, ?")
	for range db.parameters {
		fmt.Fprint(&baseStatement, ", ?")
	}
	fmt.Fprint(&baseStatement, ")")

	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(baseStatement.String())
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, d := range data {
		args := []interface{}{
			d.ValidDate.Unix(),
			d.SID,
			d.Lat,
			d.Lon,
			d.Elev,
		}
		for _, p := range db.parameters {
			value, ok := d.Data[p.Parameter]
			if !ok {
				tx.Rollback()
				return fmt.Errorf("missing value for parameter %s", p.Parameter)
			}
			args = append(args, value)
		}

		if _, err := stmt.Exec(args...); err != nil {
			tx.Rollback()
			return err
		}
	}

	stmt.Close()

	return tx.Commit()
}

var tableNamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)

func isValidNameInSql(word string) bool {
	return tableNamePattern.MatchString(word)
}

func (db *Database) addObstype(obstype string, parameters []Parameter) error {
	obstype = strings.ToUpper(obstype)
	if !isValidNameInSql(obstype) {
		return errors.New("invalid obstype")
	}
	for _, p := range parameters {
		if !isValidNameInSql(p.Parameter) {
			return fmt.Errorf("invalid parameter: %s", p.Parameter)
		}
	}

	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(fmt.Sprintf("CREATE TABLE %s_params(parameter TEXT, accum_hours REAL, units TEXT)", obstype)); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(fmt.Sprintf("CREATE UNIQUE INDEX %s_params_parameter ON %s_params(parameter)", obstype, obstype)); err != nil {
		return err
	}
	insertParameters := fmt.Sprintf("INSERT INTO %s_params VALUES (?, ?, ?)", obstype)
	for _, p := range parameters {
		if _, err := tx.Exec(insertParameters, p.Parameter, p.AccumHours, p.Units); err != nil {
			return err
		}
	}

	var mkTable strings.Builder
	fmt.Fprintf(&mkTable, "CREATE TABLE %s (validdate REAL, SID INTEGER, lat REAL, lon REAL, elev INTEGER", obstype)
	for _, p := range parameters {
		fmt.Fprintf(&mkTable, ", %s REAL", p.Parameter)
	}
	fmt.Fprint(&mkTable, ")")
	if _, err := tx.Exec(mkTable.String()); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(fmt.Sprintf("CREATE UNIQUE INDEX %s_SID_validdate ON %s(SID, validdate)", obstype, obstype)); err != nil {
		return err
	}

	return tx.Commit()
}
