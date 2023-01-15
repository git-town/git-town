package storage

import (
	"fmt"
	"sync"

	"github.com/cucumber/messages-go/v16"
	"github.com/hashicorp/go-memdb"

	"github.com/cucumber/godog/internal/models"
)

const (
	writeMode bool = true
	readMode  bool = false

	tableFeature         string = "feature"
	tableFeatureIndexURI string = "id"

	tablePickle         string = "pickle"
	tablePickleIndexID  string = "id"
	tablePickleIndexURI string = "uri"

	tablePickleStep        string = "pickle_step"
	tablePickleStepIndexID string = "id"

	tablePickleResult              string = "pickle_result"
	tablePickleResultIndexPickleID string = "id"

	tablePickleStepResult                  string = "pickle_step_result"
	tablePickleStepResultIndexPickleStepID string = "id"
	tablePickleStepResultIndexPickleID     string = "pickle_id"
	tablePickleStepResultIndexStatus       string = "status"

	tableStepDefintionMatch            string = "step_defintion_match"
	tableStepDefintionMatchIndexStepID string = "id"
)

// Storage is a thread safe in-mem storage
type Storage struct {
	db *memdb.MemDB

	testRunStarted     models.TestRunStarted
	testRunStartedLock *sync.Mutex
}

// NewStorage will create an in-mem storage that
// is used across concurrent runners and formatters
func NewStorage() *Storage {
	schema := memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			tableFeature: {
				Name: tableFeature,
				Indexes: map[string]*memdb.IndexSchema{
					tableFeatureIndexURI: {
						Name:    tableFeatureIndexURI,
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Uri"},
					},
				},
			},
			tablePickle: {
				Name: tablePickle,
				Indexes: map[string]*memdb.IndexSchema{
					tablePickleIndexID: {
						Name:    tablePickleIndexID,
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Id"},
					},
					tablePickleIndexURI: {
						Name:    tablePickleIndexURI,
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Uri"},
					},
				},
			},
			tablePickleStep: {
				Name: tablePickleStep,
				Indexes: map[string]*memdb.IndexSchema{
					tablePickleStepIndexID: {
						Name:    tablePickleStepIndexID,
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Id"},
					},
				},
			},
			tablePickleResult: {
				Name: tablePickleResult,
				Indexes: map[string]*memdb.IndexSchema{
					tablePickleResultIndexPickleID: {
						Name:    tablePickleResultIndexPickleID,
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "PickleID"},
					},
				},
			},
			tablePickleStepResult: {
				Name: tablePickleStepResult,
				Indexes: map[string]*memdb.IndexSchema{
					tablePickleStepResultIndexPickleStepID: {
						Name:    tablePickleStepResultIndexPickleStepID,
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "PickleStepID"},
					},
					tablePickleStepResultIndexPickleID: {
						Name:    tablePickleStepResultIndexPickleID,
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "PickleID"},
					},
					tablePickleStepResultIndexStatus: {
						Name:    tablePickleStepResultIndexStatus,
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "Status"},
					},
				},
			},
			tableStepDefintionMatch: {
				Name: tableStepDefintionMatch,
				Indexes: map[string]*memdb.IndexSchema{
					tableStepDefintionMatchIndexStepID: {
						Name:    tableStepDefintionMatchIndexStepID,
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "StepID"},
					},
				},
			},
		},
	}

	db, err := memdb.NewMemDB(&schema)
	if err != nil {
		panic(err)
	}

	return &Storage{db: db, testRunStartedLock: new(sync.Mutex)}
}

// MustInsertPickle will insert a pickle and it's steps,
// will panic on error.
func (s *Storage) MustInsertPickle(p *messages.Pickle) {
	txn := s.db.Txn(writeMode)

	if err := txn.Insert(tablePickle, p); err != nil {
		panic(err)
	}

	for _, step := range p.Steps {
		if err := txn.Insert(tablePickleStep, step); err != nil {
			panic(err)
		}
	}

	txn.Commit()
}

// MustGetPickle will retrieve a pickle by id and panic on error.
func (s *Storage) MustGetPickle(id string) *messages.Pickle {
	v := s.mustFirst(tablePickle, tablePickleIndexID, id)
	return v.(*messages.Pickle)
}

// MustGetPickles will retrieve pickles by URI and panic on error.
func (s *Storage) MustGetPickles(uri string) (ps []*messages.Pickle) {
	it := s.mustGet(tablePickle, tablePickleIndexURI, uri)
	for v := it.Next(); v != nil; v = it.Next() {
		ps = append(ps, v.(*messages.Pickle))
	}

	return
}

// MustGetPickleStep will retrieve a pickle step and panic on error.
func (s *Storage) MustGetPickleStep(id string) *messages.PickleStep {
	v := s.mustFirst(tablePickleStep, tablePickleStepIndexID, id)
	return v.(*messages.PickleStep)
}

// MustInsertTestRunStarted will set the test run started event and panic on error.
func (s *Storage) MustInsertTestRunStarted(trs models.TestRunStarted) {
	s.testRunStartedLock.Lock()
	defer s.testRunStartedLock.Unlock()

	s.testRunStarted = trs
}

// MustGetTestRunStarted will retrieve the test run started event and panic on error.
func (s *Storage) MustGetTestRunStarted() models.TestRunStarted {
	s.testRunStartedLock.Lock()
	defer s.testRunStartedLock.Unlock()

	return s.testRunStarted
}

// MustInsertPickleResult will instert a pickle result and panic on error.
func (s *Storage) MustInsertPickleResult(pr models.PickleResult) {
	s.mustInsert(tablePickleResult, pr)
}

// MustInsertPickleStepResult will insert a pickle step result and panic on error.
func (s *Storage) MustInsertPickleStepResult(psr models.PickleStepResult) {
	s.mustInsert(tablePickleStepResult, psr)
}

// MustGetPickleResult will retrieve a pickle result by id and panic on error.
func (s *Storage) MustGetPickleResult(id string) models.PickleResult {
	v := s.mustFirst(tablePickleResult, tablePickleResultIndexPickleID, id)
	return v.(models.PickleResult)
}

// MustGetPickleResults will retrieve all pickle results and panic on error.
func (s *Storage) MustGetPickleResults() (prs []models.PickleResult) {
	it := s.mustGet(tablePickleResult, tablePickleResultIndexPickleID)
	for v := it.Next(); v != nil; v = it.Next() {
		prs = append(prs, v.(models.PickleResult))
	}

	return prs
}

// MustGetPickleStepResult will retrieve a pickle strep result by id and panic on error.
func (s *Storage) MustGetPickleStepResult(id string) models.PickleStepResult {
	v := s.mustFirst(tablePickleStepResult, tablePickleStepResultIndexPickleStepID, id)
	return v.(models.PickleStepResult)
}

// MustGetPickleStepResultsByPickleID will retrieve pickle strep results by pickle id and panic on error.
func (s *Storage) MustGetPickleStepResultsByPickleID(pickleID string) (psrs []models.PickleStepResult) {
	it := s.mustGet(tablePickleStepResult, tablePickleStepResultIndexPickleID, pickleID)
	for v := it.Next(); v != nil; v = it.Next() {
		psrs = append(psrs, v.(models.PickleStepResult))
	}

	return psrs
}

// MustGetPickleStepResultsByStatus will retrieve pickle strep results by status and panic on error.
func (s *Storage) MustGetPickleStepResultsByStatus(status models.StepResultStatus) (psrs []models.PickleStepResult) {
	it := s.mustGet(tablePickleStepResult, tablePickleStepResultIndexStatus, status)
	for v := it.Next(); v != nil; v = it.Next() {
		psrs = append(psrs, v.(models.PickleStepResult))
	}

	return psrs
}

// MustInsertFeature will insert a feature and panic on error.
func (s *Storage) MustInsertFeature(f *models.Feature) {
	s.mustInsert(tableFeature, f)
}

// MustGetFeature will retrieve a feature by URI and panic on error.
func (s *Storage) MustGetFeature(uri string) *models.Feature {
	v := s.mustFirst(tableFeature, tableFeatureIndexURI, uri)
	return v.(*models.Feature)
}

// MustGetFeatures will retrieve all features by and panic on error.
func (s *Storage) MustGetFeatures() (fs []*models.Feature) {
	it := s.mustGet(tableFeature, tableFeatureIndexURI)
	for v := it.Next(); v != nil; v = it.Next() {
		fs = append(fs, v.(*models.Feature))
	}

	return
}

type stepDefinitionMatch struct {
	StepID         string
	StepDefinition *models.StepDefinition
}

// MustInsertStepDefintionMatch will insert the matched StepDefintion for the step ID and panic on error.
func (s *Storage) MustInsertStepDefintionMatch(stepID string, match *models.StepDefinition) {
	d := stepDefinitionMatch{
		StepID:         stepID,
		StepDefinition: match,
	}

	s.mustInsert(tableStepDefintionMatch, d)
}

// MustGetStepDefintionMatch will retrieve the matched StepDefintion for the step ID and panic on error.
func (s *Storage) MustGetStepDefintionMatch(stepID string) *models.StepDefinition {
	v := s.mustFirst(tableStepDefintionMatch, tableStepDefintionMatchIndexStepID, stepID)
	return v.(stepDefinitionMatch).StepDefinition
}

func (s *Storage) mustInsert(table string, obj interface{}) {
	txn := s.db.Txn(writeMode)

	if err := txn.Insert(table, obj); err != nil {
		panic(err)
	}

	txn.Commit()
}

func (s *Storage) mustFirst(table, index string, args ...interface{}) interface{} {
	txn := s.db.Txn(readMode)
	defer txn.Abort()

	v, err := txn.First(table, index, args...)
	if err != nil {
		panic(err)
	} else if v == nil {
		err = fmt.Errorf("Couldn't find index: %q in table: %q with args: %+v", index, table, args)
		panic(err)
	}

	return v
}

func (s *Storage) mustGet(table, index string, args ...interface{}) memdb.ResultIterator {
	txn := s.db.Txn(readMode)
	defer txn.Abort()

	it, err := txn.Get(table, index, args...)
	if err != nil {
		panic(err)
	}

	return it
}
