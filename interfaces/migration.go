package interfaces

type IMigrations interface {
	Run()
	HasChanges() (bool, []string)
}
