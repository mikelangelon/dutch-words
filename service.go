package main

type Service struct {
	store store
}

func NewService(store store) Service {
	return Service{store: store}
}

func (s Service) InsertWord(word Word) error {
	return s.store.Insert(word)
}

func (s Service) DeleteWord(id string) error {
	return s.store.Delete(id)
}

func (s Service) FindWordByID(id string) (*Word, error) {
	return s.store.FindByID(id)
}

func (s Service) FindWordByDutch(dutch string) (*Word, error) {
	return s.store.FindByDutch(dutch)
}
func (s Service) FindAllWords() ([]*Word, error) {
	return s.store.FindAll()
}
