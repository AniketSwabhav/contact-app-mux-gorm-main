package repository

import (
	"contact_app_mux_gorm_main/components/apperror"

	"github.com/jinzhu/gorm"
)

type Repository interface {
	Add(uow *UnitOfWork, out interface{}) error
}

type GormRepository struct{}

func NewGormRepository() *GormRepository {
	return &GormRepository{}
}

type UnitOfWork struct {
	DB        *gorm.DB
	Committed bool
	Readonly  bool
}

func NewUnitOfWork(db *gorm.DB, readonly bool) *UnitOfWork {
	commit := false
	if readonly {
		return &UnitOfWork{
			DB:        db.New(),
			Committed: commit,
			Readonly:  readonly,
		}
	}

	return &UnitOfWork{
		DB:        db.New().Begin(),
		Committed: commit,
		Readonly:  readonly,
	}
}

func executeQueryProcessors(db *gorm.DB, out interface{}, queryProcessors ...QueryProcessor) (*gorm.DB, error) {
	var err error
	for _, query := range queryProcessors {
		if query != nil {
			db, err = query(db, out)
			if err != nil {
				return db, err
			}
		}
	}
	return db, nil
}

func (repository *GormRepository) Add(uow *UnitOfWork, out interface{}) error {
	return uow.DB.Create(out).Error
}

func (uow *UnitOfWork) RollBack() {
	// This condition can be used if Rollback() is defered as soon as UOW is created.
	// So we only rollback if it's not committed.
	if !uow.Committed && !uow.Readonly {
		uow.DB.Rollback()
	}
}

func (uow *UnitOfWork) Commit() {
	if !uow.Readonly && !uow.Committed {
		uow.Committed = true
		uow.DB.Commit()
	}
}

func Filter(condition string, args ...interface{}) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		db = db.Debug().Where(condition, args...)
		return db, nil
	}
}

func DoesEmailExist(db *gorm.DB, email string, out interface{}, queryProcessors ...QueryProcessor) (bool, error) {
	if email == "" {
		return false, apperror.NewNotFoundError("email not present")
	}
	count := 0
	// Below comment would make the tenant check before all query processor (Uncomment only if needed in future)
	// queryProcessors = append([]QueryProcessor{Filter("tenant_id = ?", tenantID)},queryProcessors... )
	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return false, err
	}
	if err := db.Debug().Model(out).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

// func MatchPassword(db *gorm.DB, password string, out interface{}, queryProcessors ...QueryProcessor) (bool, error) {
// 	db, err := executeQueryProcessors(db, out, queryProcessors...)
// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }
