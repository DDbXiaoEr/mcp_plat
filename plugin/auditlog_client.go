// Copyright (C) 2026 Zhaoquan Wang
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package plugin

import (
	"context"
)

func LogAccess(ctx context.Context, lb *GrpcLB, req *LogAccessRequest) (*LogAccessResponse, error) {
	conn, release := lb.Pick()
	if conn == nil {
		return nil, ctx.Err()
	}
	defer release()

	client := NewAuditLogServiceClient(conn)
	return client.LogAccess(ctx, req)
}
