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

// Route is an object representing the database table.
type Route struct {
	ID            int64   `boil:"id" json:"id" toml:"id" yaml:"id"`
	SourceID      int64   `boil:"source_id" json:"source_id" toml:"source_id" yaml:"source_id"`
	DestinationID int64   `boil:"destination_id" json:"destination_id" toml:"destination_id" yaml:"destination_id"`
	Price         float64 `boil:"price" json:"price" toml:"price" yaml:"price"`

	R *routeR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L routeL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var RouteColumns = struct {
	ID            string
	SourceID      string
	DestinationID string
	Price         string
}{
	ID:            "id",
	SourceID:      "source_id",
	DestinationID: "destination_id",
	Price:         "price",
}

// Generated where

type whereHelperfloat64 struct{ field string }

func (w whereHelperfloat64) EQ(x float64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperfloat64) NEQ(x float64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelperfloat64) LT(x float64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperfloat64) LTE(x float64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelperfloat64) GT(x float64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperfloat64) GTE(x float64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}
func (w whereHelperfloat64) IN(slice []float64) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperfloat64) NIN(slice []float64) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

var RouteWhere = struct {
	ID            whereHelperint64
	SourceID      whereHelperint64
	DestinationID whereHelperint64
	Price         whereHelperfloat64
}{
	ID:            whereHelperint64{field: "\"routes\".\"id\""},
	SourceID:      whereHelperint64{field: "\"routes\".\"source_id\""},
	DestinationID: whereHelperint64{field: "\"routes\".\"destination_id\""},
	Price:         whereHelperfloat64{field: "\"routes\".\"price\""},
}

// RouteRels is where relationship names are stored.
var RouteRels = struct {
	Destination string
	Source      string
}{
	Destination: "Destination",
	Source:      "Source",
}

// routeR is where relationships are stored.
type routeR struct {
	Destination *Airport `boil:"Destination" json:"Destination" toml:"Destination" yaml:"Destination"`
	Source      *Airport `boil:"Source" json:"Source" toml:"Source" yaml:"Source"`
}

// NewStruct creates a new relationship struct
func (*routeR) NewStruct() *routeR {
	return &routeR{}
}

// routeL is where Load methods for each relationship are stored.
type routeL struct{}

var (
	routeAllColumns            = []string{"id", "source_id", "destination_id", "price"}
	routeColumnsWithoutDefault = []string{"source_id", "destination_id", "price"}
	routeColumnsWithDefault    = []string{"id"}
	routePrimaryKeyColumns     = []string{"id"}
)

type (
	// RouteSlice is an alias for a slice of pointers to Route.
	// This should generally be used opposed to []Route.
	RouteSlice []*Route
	// RouteHook is the signature for custom Route hook methods
	RouteHook func(context.Context, boil.ContextExecutor, *Route) error

	routeQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	routeType                 = reflect.TypeOf(&Route{})
	routeMapping              = queries.MakeStructMapping(routeType)
	routePrimaryKeyMapping, _ = queries.BindMapping(routeType, routeMapping, routePrimaryKeyColumns)
	routeInsertCacheMut       sync.RWMutex
	routeInsertCache          = make(map[string]insertCache)
	routeUpdateCacheMut       sync.RWMutex
	routeUpdateCache          = make(map[string]updateCache)
	routeUpsertCacheMut       sync.RWMutex
	routeUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var routeBeforeInsertHooks []RouteHook
var routeBeforeUpdateHooks []RouteHook
var routeBeforeDeleteHooks []RouteHook
var routeBeforeUpsertHooks []RouteHook

var routeAfterInsertHooks []RouteHook
var routeAfterSelectHooks []RouteHook
var routeAfterUpdateHooks []RouteHook
var routeAfterDeleteHooks []RouteHook
var routeAfterUpsertHooks []RouteHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Route) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range routeBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Route) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range routeBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Route) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range routeBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Route) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range routeBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Route) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range routeAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Route) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range routeAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Route) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range routeAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Route) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range routeAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Route) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range routeAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddRouteHook registers your hook function for all future operations.
func AddRouteHook(hookPoint boil.HookPoint, routeHook RouteHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		routeBeforeInsertHooks = append(routeBeforeInsertHooks, routeHook)
	case boil.BeforeUpdateHook:
		routeBeforeUpdateHooks = append(routeBeforeUpdateHooks, routeHook)
	case boil.BeforeDeleteHook:
		routeBeforeDeleteHooks = append(routeBeforeDeleteHooks, routeHook)
	case boil.BeforeUpsertHook:
		routeBeforeUpsertHooks = append(routeBeforeUpsertHooks, routeHook)
	case boil.AfterInsertHook:
		routeAfterInsertHooks = append(routeAfterInsertHooks, routeHook)
	case boil.AfterSelectHook:
		routeAfterSelectHooks = append(routeAfterSelectHooks, routeHook)
	case boil.AfterUpdateHook:
		routeAfterUpdateHooks = append(routeAfterUpdateHooks, routeHook)
	case boil.AfterDeleteHook:
		routeAfterDeleteHooks = append(routeAfterDeleteHooks, routeHook)
	case boil.AfterUpsertHook:
		routeAfterUpsertHooks = append(routeAfterUpsertHooks, routeHook)
	}
}

// One returns a single route record from the query.
func (q routeQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Route, error) {
	o := &Route{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for routes")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Route records from the query.
func (q routeQuery) All(ctx context.Context, exec boil.ContextExecutor) (RouteSlice, error) {
	var o []*Route

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Route slice")
	}

	if len(routeAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Route records in the query.
func (q routeQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count routes rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q routeQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if routes exists")
	}

	return count > 0, nil
}

// Destination pointed to by the foreign key.
func (o *Route) Destination(mods ...qm.QueryMod) airportQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.DestinationID),
	}

	queryMods = append(queryMods, mods...)

	query := Airports(queryMods...)
	queries.SetFrom(query.Query, "\"airports\"")

	return query
}

// Source pointed to by the foreign key.
func (o *Route) Source(mods ...qm.QueryMod) airportQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.SourceID),
	}

	queryMods = append(queryMods, mods...)

	query := Airports(queryMods...)
	queries.SetFrom(query.Query, "\"airports\"")

	return query
}

// LoadDestination allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (routeL) LoadDestination(ctx context.Context, e boil.ContextExecutor, singular bool, maybeRoute interface{}, mods queries.Applicator) error {
	var slice []*Route
	var object *Route

	if singular {
		object = maybeRoute.(*Route)
	} else {
		slice = *maybeRoute.(*[]*Route)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &routeR{}
		}
		args = append(args, object.DestinationID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &routeR{}
			}

			for _, a := range args {
				if a == obj.DestinationID {
					continue Outer
				}
			}

			args = append(args, obj.DestinationID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`airports`),
		qm.WhereIn(`airports.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Airport")
	}

	var resultSlice []*Airport
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Airport")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for airports")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for airports")
	}

	if len(routeAfterSelectHooks) != 0 {
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
		object.R.Destination = foreign
		if foreign.R == nil {
			foreign.R = &airportR{}
		}
		foreign.R.DestinationRoutes = append(foreign.R.DestinationRoutes, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.DestinationID == foreign.ID {
				local.R.Destination = foreign
				if foreign.R == nil {
					foreign.R = &airportR{}
				}
				foreign.R.DestinationRoutes = append(foreign.R.DestinationRoutes, local)
				break
			}
		}
	}

	return nil
}

// LoadSource allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (routeL) LoadSource(ctx context.Context, e boil.ContextExecutor, singular bool, maybeRoute interface{}, mods queries.Applicator) error {
	var slice []*Route
	var object *Route

	if singular {
		object = maybeRoute.(*Route)
	} else {
		slice = *maybeRoute.(*[]*Route)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &routeR{}
		}
		args = append(args, object.SourceID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &routeR{}
			}

			for _, a := range args {
				if a == obj.SourceID {
					continue Outer
				}
			}

			args = append(args, obj.SourceID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`airports`),
		qm.WhereIn(`airports.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Airport")
	}

	var resultSlice []*Airport
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Airport")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for airports")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for airports")
	}

	if len(routeAfterSelectHooks) != 0 {
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
		object.R.Source = foreign
		if foreign.R == nil {
			foreign.R = &airportR{}
		}
		foreign.R.SourceRoutes = append(foreign.R.SourceRoutes, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.SourceID == foreign.ID {
				local.R.Source = foreign
				if foreign.R == nil {
					foreign.R = &airportR{}
				}
				foreign.R.SourceRoutes = append(foreign.R.SourceRoutes, local)
				break
			}
		}
	}

	return nil
}

// SetDestination of the route to the related item.
// Sets o.R.Destination to related.
// Adds o to related.R.DestinationRoutes.
func (o *Route) SetDestination(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Airport) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"routes\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 0, []string{"destination_id"}),
		strmangle.WhereClause("\"", "\"", 0, routePrimaryKeyColumns),
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

	o.DestinationID = related.ID
	if o.R == nil {
		o.R = &routeR{
			Destination: related,
		}
	} else {
		o.R.Destination = related
	}

	if related.R == nil {
		related.R = &airportR{
			DestinationRoutes: RouteSlice{o},
		}
	} else {
		related.R.DestinationRoutes = append(related.R.DestinationRoutes, o)
	}

	return nil
}

// SetSource of the route to the related item.
// Sets o.R.Source to related.
// Adds o to related.R.SourceRoutes.
func (o *Route) SetSource(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Airport) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"routes\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 0, []string{"source_id"}),
		strmangle.WhereClause("\"", "\"", 0, routePrimaryKeyColumns),
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

	o.SourceID = related.ID
	if o.R == nil {
		o.R = &routeR{
			Source: related,
		}
	} else {
		o.R.Source = related
	}

	if related.R == nil {
		related.R = &airportR{
			SourceRoutes: RouteSlice{o},
		}
	} else {
		related.R.SourceRoutes = append(related.R.SourceRoutes, o)
	}

	return nil
}

// Routes retrieves all the records using an executor.
func Routes(mods ...qm.QueryMod) routeQuery {
	mods = append(mods, qm.From("\"routes\""))
	return routeQuery{NewQuery(mods...)}
}

// FindRoute retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindRoute(ctx context.Context, exec boil.ContextExecutor, iD int64, selectCols ...string) (*Route, error) {
	routeObj := &Route{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"routes\" where \"id\"=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, routeObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from routes")
	}

	return routeObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Route) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no routes provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(routeColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	routeInsertCacheMut.RLock()
	cache, cached := routeInsertCache[key]
	routeInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			routeAllColumns,
			routeColumnsWithDefault,
			routeColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(routeType, routeMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(routeType, routeMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"routes\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"routes\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT \"%s\" FROM \"routes\" WHERE %s", strings.Join(returnColumns, "\",\""), strmangle.WhereClause("\"", "\"", 0, routePrimaryKeyColumns))
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
		return errors.Wrap(err, "models: unable to insert into routes")
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
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == routeMapping["id"] {
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
		return errors.Wrap(err, "models: unable to populate default values for routes")
	}

CacheNoHooks:
	if !cached {
		routeInsertCacheMut.Lock()
		routeInsertCache[key] = cache
		routeInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Route.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Route) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	routeUpdateCacheMut.RLock()
	cache, cached := routeUpdateCache[key]
	routeUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			routeAllColumns,
			routePrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update routes, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"routes\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 0, wl),
			strmangle.WhereClause("\"", "\"", 0, routePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(routeType, routeMapping, append(wl, routePrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update routes row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for routes")
	}

	if !cached {
		routeUpdateCacheMut.Lock()
		routeUpdateCache[key] = cache
		routeUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q routeQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for routes")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for routes")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o RouteSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), routePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"routes\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, routePrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in route slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all route")
	}
	return rowsAff, nil
}

// Delete deletes a single Route record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Route) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Route provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), routePrimaryKeyMapping)
	sql := "DELETE FROM \"routes\" WHERE \"id\"=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from routes")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for routes")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q routeQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no routeQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from routes")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for routes")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o RouteSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(routeBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), routePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"routes\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, routePrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from route slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for routes")
	}

	if len(routeAfterDeleteHooks) != 0 {
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
func (o *Route) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindRoute(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *RouteSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := RouteSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), routePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"routes\".* FROM \"routes\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, routePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in RouteSlice")
	}

	*o = slice

	return nil
}

// RouteExists checks if the Route row exists.
func RouteExists(ctx context.Context, exec boil.ContextExecutor, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"routes\" where \"id\"=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if routes exists")
	}

	return exists, nil
}
