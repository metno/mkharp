package obs

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/metno/mkharp/internal/harp"
)

type Database struct {
	db      *sql.DB
	obstype string
}

type Data struct {
	SID          int
	Lat          float32
	Lon          float32
	Elev         int
	Observations []Observation
}

func (d *Data) Parameters() ([]harp.Parameter, error) {
	encounterdParameters := make(map[string]bool)
	for _, obs := range d.Observations {
		for parameter := range obs.Data {
			encounterdParameters[parameter] = true
		}
	}
	var ret []harp.Parameter
	for p := range encounterdParameters {
		harpParameter := harp.GetParameter(p)
		if harpParameter == nil {
			return nil, fmt.Errorf("%s: no such parameter", p)
		}
		ret = append(ret, *harpParameter)
	}

	sort.Slice(
		ret,
		func(i, j int) bool {
			return ret[i].Parameter < ret[j].Parameter
		},
	)

	return ret, nil
}

type Observation struct {
	ValidDate time.Time
	Data      map[string]float32
}

func Open(filename, obstype string) (*Database, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	return &Database{
		db:      db,
		obstype: obstype,
	}, nil
}

func Create(filename, obstype string, parameters []harp.Parameter) (*Database, error) {
	db, err := Open(filename, obstype)
	if err != nil {
		return nil, err
	}

	if err := db.addObstype(obstype, parameters); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func (db *Database) Close() error {
	return nil
}

func (db *Database) Add(data Data) error {
	var baseStatement strings.Builder

	allParameters, err := data.Parameters()
	if err != nil {
		return err
	}

	fmt.Fprintf(&baseStatement, "INSERT INTO %s (validdate, SID, lat, lon, elev", strings.ToUpper(db.obstype))
	for _, p := range allParameters {
		fmt.Fprintf(&baseStatement, ", %s", p.Parameter)
	}
	fmt.Fprint(&baseStatement, ") VALUES (?, ?, ?, ?, ?")
	for range allParameters {
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

	for _, obs := range data.Observations {
		args := []interface{}{
			obs.ValidDate.Unix(),
			data.SID,
			data.Lat,
			data.Lon,
			data.Elev,
		}
		for _, p := range allParameters {
			var value interface{}
			value, ok := obs.Data[p.Parameter]
			if !ok {
				value = ""
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

func (db *Database) addObstype(obstype string, parameters []harp.Parameter) error {
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
