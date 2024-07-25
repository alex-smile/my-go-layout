package util

import (
	"encoding/hex"

	"github.com/gofrs/uuid"
)

// GenUUID4 ...
func GenUUID4() string {
	return hex.EncodeToString(uuid.Must(uuid.NewV4()).Bytes())
}
