package mgo

import (
	"github.com/jucardi/go-db"
	"github.com/jucardi/go-db/common"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type IQuery interface {
	dbx.IQuery

	// The default batch size is defined by the Database itself.  As of this writing, MongoDB will use an initial size of min(100 docs, 4MB) on the first batch, and 4MB on remaining ones.
	Batch(n int) IQuery

	// Prefetch sets the point at which the next batch of results will be requested. When there are p*batch_size remaining documents cached in an Iter, the next batch will be
	// requested in background. For instance, when using this:
	Prefetch(p float64) IQuery

	// Explain returns a number of details about how the MongoDB server would execute the requested query, such as the number of objects examined, the number of times the read lock
	// was yielded to allow writes to go in, and so on.
	Explain(result interface{}) error

	// Hint will include an explicit "hint" in the query to force the server to use a specified index, potentially improving performance in some situations. The provided parameters
	// are the fields that compose the key of the index to be used. For details on how the indexKey may be built, see the EnsureIndex method.
	Hint(indexKey ...string) IQuery

	// SetMaxScan constrains the query to stop after scanning the specified number of documents.
	SetMaxScan(n int) IQuery

	// SetMaxTime constrains the query to stop after running for the specified time. When the time limit is reached MongoDB automatically cancels the query.
	SetMaxTime(d time.Duration) IQuery

	// Snapshot will force the performed query to make use of an available index on the _id field to prevent the same document from being returned more than once in a single
	// iteration. This might happen without this setting in situations when the document changes in size and thus has to be moved while the iteration is running.
	Snapshot() IQuery

	// Comment adds a comment to the query to identify it in the Database profiler output.
	Comment(comment string) IQuery

	// LogReplay enables an option that optimizes queries that are typically made on the MongoDB oplog for replaying it. This is an internal implementation aspect and most likely
	// uninteresting for other uses. It has seen at least one use case, though, so it's exposed via the API.
	LogReplay() IQuery

	// Iter executes the query and returns an iterator capable of going over all the results. Results will be returned in batches of configurable size (see the Batch method) and more
	// documents will be requested when a configurable number of documents is iterated over (see the Prefetch method).
	Iter() IIter

	// Tail returns a tailable iterator. Unlike a normal iterator, a tailable iterator may wait for new values to be inserted in the Collection once the end of the current result set
	// is reached, A tailable iterator may only be used with capped collections.
	//     - See the Tail documentation in `gopkg.in/mgo.v2` for more information.
	Tail(timeout time.Duration) IIter

	// MapReduce executes a map/reduce job for documents covered by the query. That kind of job is suitable for very flexible bulk aggregation of data performed at the server side
	// via Javascript functions.
	//     - See the MapReduce documentation in `gopkg.in/mgo.v2` for more information.
	MapReduce(job *MapReduce, result interface{}) (info *MapReduceInfo, err error)

	// Apply runs the findAndModify MongoDB command, which allows updating, upserting or removing a document matching a query and atomically returning either the old version (the
	// default) or the new version of the document (when ReturnNew is true). If no objects are found Apply returns ErrNotFound.
	//     - See the Apply documentation in `gopkg.in/mgo.v2` for more information.
	Apply(change Change, result interface{}) (info *ChangeInfo, err error)

	// Returns the internal mgo.query used by this implementation.
	Q() *mgo.Query
}

// query is the default implementation of IQuery
type query struct {
	*common.AbstractQuery
	qry       *mgo.Query
	col       *mgo.Collection
	batch     *int
	prefetch  *float64
	maxScan   *int
	maxTime   *time.Duration
	snapshot  bool
	logReplay bool
	hints     [][]string
	comments  []string
}

func (q *query) Count() (n int, err error) {
	return q.prepare().Count()
}

func (q *query) First(result interface{}) error {
	return q.One(result)
}

func (q *query) One(result interface{}) error {
	return q.prepare().One(result)
}

func (q *query) Last(result interface{}) error {
	count, err := q.Count()
	if err != nil {
		return err
	}
	return q.prepare().Skip(count - 1).One(result)
}

func (q *query) All(result interface{}) error {
	return q.prepare().All(result)
}

func (q *query) Distinct(key string, result interface{}) error {
	return q.prepare().Distinct(key, result)
}

func (q *query) Update(update interface{}) error {
	_, err := q.prepare().Apply(mgo.Change{
		Update: update,
	}, nil)
	return err
}

func (q *query) Delete() error {
	return q.Remove()
}

func (q *query) Remove() error {
	_, err := q.prepare().Apply(mgo.Change{
		Remove: true,
	}, nil)
	return err
}

func (q *query) Explain(result interface{}) error {
	return q.prepare().Explain(result)
}

func (q *query) Iter() IIter {
	return q.prepare().Iter()
}

func (q *query) Tail(timeout time.Duration) IIter {
	return q.prepare().Tail(timeout)
}

func (q *query) MapReduce(job *MapReduce, result interface{}) (*MapReduceInfo, error) {
	info, err := q.prepare().MapReduce(makeMapReduce(job), result)
	return makeMapReduceInfo(info), err
}

func (q *query) Apply(change Change, result interface{}) (*ChangeInfo, error) {
	info, err := q.prepare().Apply(mgo.Change(change), result)
	return makeChangeInfo(info), err
}

func (q *query) Batch(n int) IQuery {
	q.batch = &n
	return q
}

func (q *query) Prefetch(p float64) IQuery {
	q.prefetch = &p
	return q
}

func (q *query) Hint(indexKey ...string) IQuery {
	q.hints = append(q.hints, indexKey)
	return q
}

func (q *query) SetMaxScan(n int) IQuery {
	q.maxScan = &n
	return q
}

func (q *query) SetMaxTime(d time.Duration) IQuery {
	q.maxTime = &d
	return q
}

func (q *query) Snapshot() IQuery {
	q.snapshot = true
	return q
}

func (q *query) Comment(comment string) IQuery {
	q.comments = append(q.comments, comment)
	return q
}

func (q *query) LogReplay() IQuery {
	q.logReplay = true
	return q
}

func (q *query) Q() *mgo.Query {
	return q.qry
}

func (q *query) prepare() *mgo.Query {
	q.qry = q.col.Find(q.makeQuery())
	if q.LimitVal != nil {
		q.qry = q.qry.Limit(*q.LimitVal)
	}
	if q.SkipVal != nil {
		q.qry = q.qry.Skip(*q.SkipVal)
	}
	if len(q.SortFields) > 0 {
		q.qry = q.qry.Sort(q.SortFields...)
	}
	for _, cond := range q.Selects {
		q.qry = q.qry.Select(cond.Query)
	}
	if q.batch != nil {
		q.qry = q.qry.Batch(*q.batch)
	}
	if q.prefetch != nil {
		q.qry = q.qry.Prefetch(*q.prefetch)
	}
	if q.maxScan != nil {
		q.qry = q.qry.SetMaxScan(*q.maxScan)
	}
	if q.maxTime != nil {
		q.qry = q.qry.SetMaxTime(*q.maxTime)
	}
	if q.snapshot {
		q.qry = q.qry.Snapshot()
	}
	if q.logReplay {
		q.qry = q.qry.LogReplay()
	}
	for _, h := range q.hints {
		q.qry = q.qry.Hint(h...)
	}
	for _, c := range q.comments {
		q.qry = q.qry.Comment(c)
	}

	return q.qry
}

func (q *query) makeQuery() interface{} {
	var blocks []interface{}

	for _, block := range q.Queries {
		var current []interface{}
		for _, cond := range block {
			if cond.Negation {
				current = append(current, bson.M{"$not": cond.Query})
			} else {
				current = append(current, cond.Query)
			}
		}

		switch len(current) {
		case 0:
			continue
		case 1:
			blocks = append(blocks, current[0])
		default:
			blocks = append(blocks, bson.M{"$and": current})
		}
	}
	switch len(blocks) {
	case 0:
		return bson.M{}
	case 1:
		return blocks[0]
	default:
		return bson.M{"or": blocks}
	}
}

func newQuery(col *mgo.Collection, qry interface{}, negated bool) IQuery {
	ret := &query{
		col: col,
	}
	ret.AbstractQuery = &common.AbstractQuery{
		Q: ret,
	}
	if negated {
		ret.Not(qry)
	} else {
		ret.Where(qry)
	}

	return ret
}
