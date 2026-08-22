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
	"math"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

// grpcBackend 表示一个后端 gRPC 地址及其常驻连接与在途请求数。
type grpcBackend struct {
	addr     string
	conn     *grpc.ClientConn
	inflight atomic.Int64
}

// GrpcLB 基于「最少连接数」在多个后端 gRPC 地址间做负载均衡。
// 优先选择已就绪（READY）的连接，再按在途请求数最少进行分配。
type GrpcLB struct {
	backends []*grpcBackend
}

var (
	lbMu    sync.Mutex
	lbCache = map[string]*GrpcLB{}
)

// GetGrpcLB 获取（或创建）覆盖给定地址集合的负载均衡器。
// 地址顺序无关且自动去重，相同的地址集合复用同一实例与底层连接。
func GetGrpcLB(addrs []string) *GrpcLB {
	normalized := normalizeAddrs(addrs)
	key := strings.Join(normalized, ",")

	lbMu.Lock()
	defer lbMu.Unlock()
	if lb, ok := lbCache[key]; ok {
		return lb
	}

	lb := &GrpcLB{}
	for _, addr := range normalized {
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			continue
		}
		lb.backends = append(lb.backends, &grpcBackend{
			addr: addr,
			conn: conn,
		})
	}
	lbCache[key] = lb
	return lb
}

func normalizeAddrs(addrs []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, a := range addrs {
		a = strings.TrimSpace(a)
		if a == "" || seen[a] {
			continue
		}
		seen[a] = true
		out = append(out, a)
	}
	sort.Strings(out)
	return out
}

// Pick 返回当前负载最小的后端连接，并返回请求结束后的释放函数。
// 无可用后端时返回 nil 连接与空操作释放函数。
func (lb *GrpcLB) Pick() (*grpc.ClientConn, func()) {
	if lb == nil || len(lb.backends) == 0 {
		return nil, func() {}
	}

	var picked *grpcBackend
	best := int64(math.MaxInt64)

	for _, b := range lb.backends {
		if b.conn.GetState() != connectivity.Ready {
			continue
		}
		if n := b.inflight.Load(); n < best {
			best = n
			picked = b
		}
	}
	if picked == nil {
		best = int64(math.MaxInt64)
		for _, b := range lb.backends {
			if n := b.inflight.Load(); n < best {
				best = n
				picked = b
			}
		}
	}

	picked.inflight.Add(1)
	return picked.conn, func() { picked.inflight.Add(-1) }
}

// Len 返回后端地址数量。
func (lb *GrpcLB) Len() int {
	if lb == nil {
		return 0
	}
	return len(lb.backends)
}
