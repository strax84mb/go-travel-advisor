// Code generated by SQLBoiler 4.2.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Airport is an object representing the database table.
type Airport struct {
	ID        int64  `boil:"id" json:"id" toml:"id" yaml:"id"`
	AirportID int64  `boil:"airport_id" json:"airport_id" toml:"airport_id" yaml:"airport_id"`
	Name      string `boil:"name" json:"name" toml:"name" yaml:"name"`
	CityID    int64  `boil:"city_id" json:"city_id" toml:"city_id" yaml:"city_id"`

	R *airportR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L airportL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var AirportColumns = struct {
	ID        string
	AirportID string
	Name      string
	CityID    string
}{
	ID:        "id",
	AirportID: "airport_id",
	Name:      "name",
	CityID:    "city_id",
}

// Generated where

type whereHelperint64 struct{ field string }

func (w whereHelperint64) EQ(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperint64) NEQ(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperint64) LT(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperint64) LTE(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperint64) GT(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperint64) GTE(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperint64) IN(slice []int64) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperint64) NIN(slice []int64) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

type whereHelperstring struct{ field string }

func (w whereHelperstring) EQ(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperstring) NEQ(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperstring) LT(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperstring) LTE(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperstring) GT(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperstring) GTE(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperstring) IN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperstring) NIN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

var AirportWhere = struct {
	ID        whereHelperint64
	AirportID whereHelperint64
	Name      whereHelperstring
	CityID    whereHelperint64
}{
	ID:        whereHelperint64{field: "\"airports\".\"id\""},
	AirportID: whereHelperint64{field: "\"airports\".\"airport_id\""},
	Name:      whereHelperstring{field: "\"airports\".\"name\""},
	CityID:    whereHelperint64{field: "\"airports\".\"city_id\""},
}

// AirportRels is where relationship names are stored.
var AirportRels = struct {
	City              string
	DestinationRoutes string
	SourceRoutes      string
}{
	City:              "City",
	DestinationRoutes: "DestinationRoutes",
	SourceRoutes:      "SourceRoutes",
}

// airportR is where relationships are stored.
type airportR struct {
	City              *City      `boil:"City" json:"City" toml:"City" yaml:"City"`
	DestinationRoutes RouteSlice `boil:"DestinationRoutes" json:"DestinationRoutes" toml:"DestinationRoutes" yaml:"DestinationRoutes"`
	SourceRoutes      RouteSlice `boil:"SourceRoutes" json:"SourceRoutes" toml:"SourceRoutes" yaml:"SourceRoutes"`
}

// NewStruct creates a new relationship struct
func (*airportR) NewStruct() *airportR {
	return &airportR{}
}

// airportL is where Load methods for each relationship are stored.
type airportL struct{}

var (
	airportAllColumns            = []string{"id", "airport_id", "name", "city_id"}
	airportColumnsWithoutDefault = []string{"airport_id", "name", "city_id"}
	airportColumnsWithDefault    = []string{"id"}
	airportPrimaryKeyColumns     = []string{"id"}
)

type (
	// AirportSlice is an alias for a slice of pointers to Airport.
	// This should generally be used opposed to []Airport.
	AirportSlice []*Airport
	// AirportHook is the signature for custom Airport hook methods
	AirportHook func(context.Context, boil.ContextExecutor, *Airport) error

	airportQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	airportType                 = reflect.TypeOf(&Airport{})
	airportMapping              = queries.MakeStructMapping(airportType)
	airportPrimaryKeyMapping, _ = queries.BindMapping(airportType, airportMapping, airportPrimaryKeyColumns)
	airportInsertCacheMut       sync.RWMutex
	airportInsertCache          = make(map[string]insertCache)
	airportUpdateCacheMut       sync.RWMutex
	airportUpdateCache          = make(map[string]updateCache)
	airportUpsertCacheMut       sync.RWMutex
	airportUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var airportBeforeInsertHooks []AirportHook
var airportBeforeUpdateHooks []AirportHook
var airportBeforeDeleteHooks []AirportHook
var airportBeforeUpsertHooks []AirportHook

var airportAfterInsertHooks []AirportHook
var airportAfterSelectHooks []AirportHook
var airportAfterUpdateHooks []AirportHook
var airportAfterDeleteHooks []AirportHook
var airportAfterUpsertHooks []AirportHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Airport) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range airportBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Airport) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range airportBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Airport) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range airportBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Airport) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range airportBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Airport) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range airportAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Airport) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range airportAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Airport) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range airportAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Airport) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range airportAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Airport) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range airportAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddAirportHook registers your hook function for all future operations.
func AddAirportHook(hookPoint boil.HookPoint, airportHook AirportHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		airportBeforeInsertHooks = append(airportBeforeInsertHooks, airportHook)
	case boil.BeforeUpdateHook:
		airportBeforeUpdateHooks = append(airportBeforeUpdateHooks, airportHook)
	case boil.BeforeDeleteHook:
		airportBeforeDeleteHooks = append(airportBeforeDeleteHooks, airportHook)
	case boil.BeforeUpsertHook:
		airportBeforeUpsertHooks = append(airportBeforeUpsertHooks, airportHook)
	case boil.AfterInsertHook:
		airportAfterInsertHooks = append(airportAfterInsertHooks, airportHook)
	case boil.AfterSelectHook:
		airportAfterSelectHooks = append(airportAfterSelectHooks, airportHook)
	case boil.AfterUpdateHook:
		airportAfterUpdateHooks = append(airportAfterUpdateHooks, airportHook)
	case boil.AfterDeleteHook:
		airportAfterDeleteHooks = append(airportAfterDeleteHooks, airportHook)
	case boil.AfterUpsertHook:
		airportAfterUpsertHooks = append(airportAfterUpsertHooks, airportHook)
	}
}

// One returns a single airport record from the query.
func (q airportQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Airport, error) {
	o := &Airport{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for airports")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Airport records from the query.
func (q airportQuery) All(ctx context.Context, exec boil.ContextExecutor) (AirportSlice, error) {
	var o []*Airport

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Airport slice")
	}

	if len(airportAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Airport records in the query.
func (q airportQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count airports rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q airportQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if airports exists")
	}

	return count > 0, nil
}

// City pointed to by the foreign key.
func (o *Airport) City(mods ...qm.QueryMod) cityQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.CityID),
	}

	queryMods = append(queryMods, mods...)

	query := Cities(queryMods...)
	queries.SetFrom(query.Query, "\"cities\"")

	return query
}

// DestinationRoutes retrieves all the route's Routes with an executor via destination_id column.
func (o *Airport) DestinationRoutes(mods ...qm.QueryMod) routeQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"routes\".\"destination_id\"=?", o.ID),
	)

	query := Routes(queryMods...)
	queries.SetFrom(query.Query, "\"routes\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"routes\".*"})
	}

	return query
}

// SourceRoutes retrieves all the route's Routes with an executor via source_id column.
func (o *Airport) SourceRoutes(mods ...qm.QueryMod) routeQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"routes\".\"source_id\"=?", o.ID),
	)

	query := Routes(queryMods...)
	queries.SetFrom(query.Query, "\"routes\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"routes\".*"})
	}

	return query
}

// LoadCity allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (airportL) LoadCity(ctx context.Context, e boil.ContextExecutor, singular bool, maybeAirport interface{}, mods queries.Applicator) error {
	var slice []*Airport
	var object *Airport

	if singular {
		object = maybeAirport.(*Airport)
	} else {
		slice = *maybeAirport.(*[]*Airport)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &airportR{}
		}
		args = append(args, object.CityID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &airportR{}
			}

			for _, a := range args {
				if a == obj.CityID {
					continue Outer
				}
			}

			args = append(args, obj.CityID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`cities`),
		qm.WhereIn(`cities.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load City")
	}

	var resultSlice []*City
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice City")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for cities")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for cities")
	}

	if len(airportAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.City = foreign
		if foreign.R == nil {
			foreign.R = &cityR{}
		}
		foreign.R.Airports = append(foreign.R.Airports, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.CityID == foreign.ID {
				local.R.City = foreign
				if foreign.R == nil {
					foreign.R = &cityR{}
				}
				foreign.R.Airports = append(foreign.R.Airports, local)
				break
			}
		}
	}

	return nil
}

// LoadDestinationRoutes allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (airportL) LoadDestinationRoutes(ctx context.Context, e boil.ContextExecutor, singular bool, maybeAirport interface{}, mods queries.Applicator) error {
	var slice []*Airport
	var object *Airport

	if singular {
		object = maybeAirport.(*Airport)
	} else {
		slice = *maybeAirport.(*[]*Airport)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &airportR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &airportR{}
			}

			for _, a := range args {
				if a == obj.ID {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`routes`),
		qm.WhereIn(`routes.destination_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load routes")
	}

	var resultSlice []*Route
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice routes")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on routes")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for routes")
	}

	if len(routeAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.DestinationRoutes = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &routeR{}
			}
			foreign.R.Destination = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.DestinationID {
				local.R.DestinationRoutes = append(local.R.DestinationRoutes, foreign)
				if foreign.R == nil {
					foreign.R = &routeR{}
				}
				foreign.R.Destination = local
				break
			}
		}
	}

	return nil
}

// LoadSourceRoutes allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (airportL) LoadSourceRoutes(ctx context.Context, e boil.ContextExecutor, singular bool, maybeAirport interface{}, mods queries.Applicator) error {
	var slice []*Airport
	var object *Airport

	if singular {
		object = maybeAirport.(*Airport)
	} else {
		slice = *maybeAirport.(*[]*Airport)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &airportR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &airportR{}
			}

			for _, a := range args {
				if a == obj.ID {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`routes`),
		qm.WhereIn(`routes.source_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load routes")
	}

	var resultSlice []*Route
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice routes")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on routes")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for routes")
	}

	if len(routeAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.SourceRoutes = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &routeR{}
			}
			foreign.R.Source = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.SourceID {
				local.R.SourceRoutes = append(local.R.SourceRoutes, foreign)
				if foreign.R == nil {
					foreign.R = &routeR{}
				}
				foreign.R.Source = local
				break
			}
		}
	}

	return nil
}

// SetCity of the airport to the related item.
// Sets o.R.City to related.
// Adds o to related.R.Airports.
func (o *Airport) SetCity(ctx context.Context, exec boil.ContextExecutor, insert bool, related *City) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"airports\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 0, []string{"city_id"}),
		strmangle.WhereClause("\"", "\"", 0, airportPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.CityID = related.ID
	if o.R == nil {
		o.R = &airportR{
			City: related,
		}
	} else {
		o.R.City = related
	}

	if related.R == nil {
		related.R = &cityR{
			Airports: AirportSlice{o},
		}
	} else {
		related.R.Airports = append(related.R.Airports, o)
	}

	return nil
}

// AddDestinationRoutes adds the given related objects to the existing relationships
// of the airport, optionally inserting them as new records.
// Appends related to o.R.DestinationRoutes.
// Sets related.R.Destination appropriately.
func (o *Airport) AddDestinationRoutes(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*Route) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.DestinationID = o.ID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"routes\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 0, []string{"destination_id"}),
				strmangle.WhereClause("\"", "\"", 0, routePrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.DestinationID = o.ID
		}
	}

	if o.R == nil {
		o.R = &airportR{
			DestinationRoutes: related,
		}
	} else {
		o.R.DestinationRoutes = append(o.R.DestinationRoutes, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &routeR{
				Destination: o,
			}
		} else {
			rel.R.Destination = o
		}
	}
	return nil
}

// AddSourceRoutes adds the given related objects to the existing relationships
// of the airport, optionally inserting them as new records.
// Appends related to o.R.SourceRoutes.
// Sets related.R.Source appropriately.
func (o *Airport) AddSourceRoutes(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*Route) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.SourceID = o.ID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"routes\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 0, []string{"source_id"}),
				strmangle.WhereClause("\"", "\"", 0, routePrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.SourceID = o.ID
		}
	}

	if o.R == nil {
		o.R = &airportR{
			SourceRoutes: related,
		}
	} else {
		o.R.SourceRoutes = append(o.R.SourceRoutes, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &routeR{
				Source: o,
			}
		} else {
			rel.R.Source = o
		}
	}
	return nil
}

// Airports retrieves all the records using an executor.
func Airports(mods ...qm.QueryMod) airportQuery {
	mods = append(mods, qm.From("\"airports\""))
	return airportQuery{NewQuery(mods...)}
}

// FindAirport retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindAirport(ctx context.Context, exec boil.ContextExecutor, iD int64, selectCols ...string) (*Airport, error) {
	airportObj := &Airport{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"airports\" where \"id\"=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, airportObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from airports")
	}

	return airportObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Airport) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no airports provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(airportColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	airportInsertCacheMut.RLock()
	cache, cached := airportInsertCache[key]
	airportInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			airportAllColumns,
			airportColumnsWithDefault,
			airportColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(airportType, airportMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(airportType, airportMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"airports\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"airports\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT \"%s\" FROM \"airports\" WHERE %s", strings.Join(returnColumns, "\",\""), strmangle.WhereClause("\"", "\"", 0, airportPrimaryKeyColumns))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	result, err := exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into airports")
	}

	var lastID int64
	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	o.ID = int64(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == airportMapping["id"] {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.ID,
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, identifierCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "models: unable to populate default values for airports")
	}

CacheNoHooks:
	if !cached {
		airportInsertCacheMut.Lock()
		airportInsertCache[key] = cache
		airportInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Airport.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Airport) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	airportUpdateCacheMut.RLock()
	cache, cached := airportUpdateCache[key]
	airportUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			airportAllColumns,
			airportPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update airports, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"airports\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 0, wl),
			strmangle.WhereClause("\"", "\"", 0, airportPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(airportType, airportMapping, append(wl, airportPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update airports row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for airports")
	}

	if !cached {
		airportUpdateCacheMut.Lock()
		airportUpdateCache[key] = cache
		airportUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q airportQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for airports")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for airports")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o AirportSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), airportPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"airports\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, airportPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in airport slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all airport")
	}
	return rowsAff, nil
}

// Delete deletes a single Airport record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Airport) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Airport provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), airportPrimaryKeyMapping)
	sql := "DELETE FROM \"airports\" WHERE \"id\"=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from airports")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for airports")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q airportQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no airportQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from airports")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for airports")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o AirportSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(airportBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), airportPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"airports\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, airportPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from airport slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for airports")
	}

	if len(airportAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Airport) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindAirport(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *AirportSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := AirportSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), airportPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"airports\".* FROM \"airports\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, airportPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in AirportSlice")
	}

	*o = slice

	return nil
}

// AirportExists checks if the Airport row exists.
func AirportExists(ctx context.Context, exec boil.ContextExecutor, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"airports\" where \"id\"=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if airports exists")
	}

	return exists, nil
}
