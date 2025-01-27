// +build go1.10

package instrumentedsql

import (
	"context"
	"database/sql/driver"
)

var _ driver.SessionResetter = wrappedConn{}

func (c wrappedConn) ResetSession(ctx context.Context) error {
	ctx, cancel := (&c.opts).setTimeout(ctx)
	defer cancel()

	conn, ok := c.parent.(driver.SessionResetter)
	if !ok {
		return nil
	}

	return conn.ResetSession(ctx)
}
