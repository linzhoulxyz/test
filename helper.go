package test

import "github.com/jinzhu/gorm"

// get gorm model table name
func getModelsTablesName(db *gorm.DB, models []interface{}) []string {
	if db == nil || len(models) == 0 {
		return nil
	}

	var modelTableNames = make([]string, len(models))
	for i, model := range models {
		modelTableNames[i] = db.NewScope(model).TableName()
	}

	return modelTableNames
}
