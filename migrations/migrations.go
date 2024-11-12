package migrations

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/rimba47prayoga/gorim.git/conf"
)

type Operation struct {
	Name	string
	Func	func(name string) error
}

type Migrations struct {
	Models		[]interface{}
	Operations	[]Operation
	hasChanges	bool
}

func (m *Migrations) CreateMigrationTable() {
	var migrationTable GorimMigrations
	if !conf.DB.Migrator().HasTable(&migrationTable) {
		// Create the table if it doesn't exist
		err := conf.DB.Migrator().CreateTable(&migrationTable);
		if err != nil {
			panic(err)
		}
	}
}

// Serialize a model's structure into a consistent string representation, excluding non-database fields
func (m *Migrations) SerializeModel(model interface{}) string {
    t := reflect.TypeOf(model).Elem() // Get the type of the model
    var sb strings.Builder

    // Iterate over the fields of the struct
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)

        // Skip fields that are structs (e.g., associations) or have no tags
        if field.Type.Kind() == reflect.Struct || field.Tag == "" {
            continue
        }

        // Include field name, type, and tag in the string representation
        sb.WriteString(field.Name)
        sb.WriteString(field.Type.String())
        sb.WriteString(fmt.Sprintf("%v", field.Tag))
    }

    return sb.String()
}

// Generate a hash from the serialized model structures
func (m *Migrations) GenerateHash() string {
    var combinedString string
    for _, model := range m.Models {
        combinedString += m.SerializeModel(model)
    }

    // Create a SHA-256 hash
    hash := sha256.Sum256([]byte(combinedString))
    return hex.EncodeToString(hash[:])
}

func (m *Migrations) AddOperation(operation Operation) {
	m.Operations = append(m.Operations, operation)
}

func (m *Migrations) MigrateModels() {
	err := conf.DB.AutoMigrate(
		m.Models...
	)
	if err != nil {
		panic(err)
	}
}

func (m *Migrations) RunMigrationModels() func(string) error {

	// wrap to function, cause it called from commandline
	return func(name string) error {
		hashVersion := m.GenerateHash()
		var migration GorimMigrations
		err := conf.DB.Where("name = ?", name).First(&migration).Error
		if err != nil {
			m.MigrateModels()
			migration := GorimMigrations{
				Name: name,
				Version: hashVersion,
				CreatedAt: time.Now(),
			}
			conf.DB.Create(&migration)
			m.hasChanges = true
			fmt.Println("Models successfully migrated.")
			return nil
		}
		if migration.Version != hashVersion {
			fmt.Println("Applying new migration..")
			m.MigrateModels()
			migration.Version = hashVersion
			now := time.Now()
			migration.UpdatedAt = &now
			conf.DB.Save(&migration)
			m.hasChanges = true
			fmt.Println("Models successfully migrated.")
			fmt.Printf("Current version: %s\n", hashVersion)
			return nil
		}
		return nil
	}
}

func (m *Migrations) RunGo(callable func() error) func(string) error {
	// wrap to function, cause it called from commandline
	return func(name string) error {
		var migration GorimMigrations
		err := conf.DB.Where("name = ?", name).First(&migration).Error
		if err == nil {
			// skip migration that already applied.
			return nil
		}
		err = callable()
		if err != nil {
			fmt.Printf("Failed to run migration: %s\n", name)
			log.Fatal(err.Error())
		}
		migration = GorimMigrations{
			Name: name,
			CreatedAt: time.Now(),
		}
		conf.DB.Create(&migration)
		m.hasChanges = true
		fmt.Printf("%s migration applied\n", name)
		return nil
	}
}

func (m *Migrations) Run() {
	m.CreateMigrationTable()
	for _, operation := range m.Operations {
		operation.Func(operation.Name)
	}
	if !m.hasChanges {
		fmt.Println("No changes detected.")
	}
}
