package swag

import "sync"

var (
	global Swagger
	gm     = &sync.RWMutex{}
)

func init() {
	ReplaceGlobal(New("swagger"))
}

func G() Swagger {
	gm.RLock()
	defer gm.RUnlock()
	return global
}

func ReplaceGlobal(s Swagger) {
	gm.Lock()
	defer gm.Unlock()
	global = s
}
