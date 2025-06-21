package metricsscrapper

type Store struct {
	metrics map[string][]byte
}

func NewStore() *Store {
	return &Store{
		metrics: map[string][]byte{},
	}
}

func (s *Store) Get(host string) ([]byte, bool) {
	content, ok := s.metrics[host]
	return content, ok
}

func (s *Store) Put(host string, content []byte) {
	s.metrics[host] = content
}
