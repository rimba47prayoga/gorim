package migrations

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/labstack/gommon/color"
	"github.com/rimba47prayoga/gorim.git/conf"
	"github.com/rimba47prayoga/gorim.git/utils"
)

var MIGRATE_MODELS string = "migrate_models"

type Operation struct {
	Name	string
	Func	func(name string) error
}

type Migrations struct {
	Models		[]interface{}
	Operations	[]Operation
	// hasChanges	bool
}

func (m *Migrations) CreateMigrationTable() {
	var migrationTable GorimMigrations
	if !m.HasTableMigrations() {
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

func (m *Migrations) AddOperation(operation ...Operation) {
	m.Operations = append(m.Operations, operation...)
}

func (m *Migrations) MigrateModels() {
	err := conf.DB.AutoMigrate(
		m.Models...
	)
	if err != nil {
		panic(err)
	}
}

func (m *Migrations) GetOperationNames() []string {
	var names []string
	for _, operation := range m.Operations {
		names = append(names, operation.Name)
	}
	return names
}

func (m *Migrations) GetUnAppliedMigrations() []string {
	var appliedMigrations []string
	err := conf.DB.
		Model(&GorimMigrations{}).
		Pluck("name", &appliedMigrations).
		Error
	if err != nil {
		log.Fatal(err)
	}

	unAppliedMigrations := []string{}
	for _, name := range m.GetOperationNames() {
		// Check if the migration name is not in the applied migrations list
		found := false
		for _, applied := range appliedMigrations {
			if name == applied {
				found = true
				break
			}
		}
		if !found {
			unAppliedMigrations = append(unAppliedMigrations, name)
		}
	}
	return unAppliedMigrations
}

func (m *Migrations) HasTableMigrations() bool {
    return conf.DB.Migrator().HasTable(&GorimMigrations{}) 
}

func (m *Migrations) HasChanges() (bool, []string) {
	if !m.HasTableMigrations() {
		return true, m.GetOperationNames()
	}
	var migration GorimMigrations
	var unAppliedMigrations []string
	err := conf.DB.Where("name = ?", MIGRATE_MODELS).First(&migration).Error
	if err != nil {
		unAppliedMigrations = append(unAppliedMigrations, MIGRATE_MODELS)
	} else {
		hashVersion := m.GenerateHash()
		if hashVersion != migration.Version {
			unAppliedMigrations = append(unAppliedMigrations, MIGRATE_MODELS)
		}
	}
	unAppliedMigrations = append(unAppliedMigrations, m.GetUnAppliedMigrations()...)
	return len(unAppliedMigrations) > 0, unAppliedMigrations
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
			return nil
		}
		if migration.Version != hashVersion {
			m.MigrateModels()
			migration.Version = hashVersion
			now := time.Now()
			migration.UpdatedAt = &now
			conf.DB.Save(&migration)
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
			return err
		}
		migration = GorimMigrations{
			Name: name,
			CreatedAt: time.Now(),
		}
		conf.DB.Create(&migration)
		return nil
	}
}

// Apply method executes a migration operation and prints status
func (m *Migrations) Apply(operation Operation) {
	// Call the operation function
	err := operation.Func(operation.Name)
	if err != nil {
		// If there is an error, print it (failed migration) in red
		fmt.Printf("  Applying %s... %s\n", operation.Name, color.Red("FAILED", color.B))
		log.Printf("Error: %v\n", err) // Log the error with details
	} else {
		// If no error, print success (mimicking Django's "OK") in green
		fmt.Printf("  Applying %s... %s\n", operation.Name, color.Green("OK", color.B))
	}
}

// Applies method iterates through the operations and applies each
func (m *Migrations) Applies(unApplied []string) {
	// Print the 'Running migrations:' header
	fmt.Println(color.Cyan("Running migrations:", color.B))

	// Apply each migration operation
	for _, operation := range m.Operations {
		if utils.Contains(unApplied, operation.Name) {
			m.Apply(operation)
		}
	}
}

func (m *Migrations) Run() {
	m.CreateMigrationTable()
	isChanged, unApplied := m.HasChanges()
	if isChanged {
		m.Applies(unApplied)
	} else {
		fmt.Println("No changes detected.")
	}
}
