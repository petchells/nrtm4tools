package persist

type Repository interface {
	InitializeConnectionPool(dbUrl string)
	GetState(string) (NRTMState, error)
	SaveState(NRTMState) error
}
