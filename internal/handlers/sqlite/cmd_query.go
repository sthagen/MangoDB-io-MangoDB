// Copyright 2021 FerretDB Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sqlite

import (
	"context"
	"fmt"
	"time"

	"github.com/FerretDB/FerretDB/internal/handlers/common"
	"github.com/FerretDB/FerretDB/internal/handlers/commonerrors"
	"github.com/FerretDB/FerretDB/internal/types"
	"github.com/FerretDB/FerretDB/internal/util/must"
	"github.com/FerretDB/FerretDB/internal/wire"
)

// CmdQuery implements HandlerInterface.
func (h *Handler) CmdQuery(ctx context.Context, query *wire.OpQuery) (*wire.OpReply, error) {
	cmd := query.Query.Command()
	collection := query.FullCollectionName

	// both are valid
	if (cmd == "ismaster" || cmd == "isMaster") && collection == "admin.$cmd" {
		return &wire.OpReply{
			NumberReturned: 1,
			Documents: []*types.Document{must.NotFail(types.NewDocument(
				"ismaster", true, // only lowercase
				"maxBsonObjectSize", int32(types.MaxDocumentLen),
				"maxMessageSizeBytes", int32(wire.MaxMsgLen),
				"maxWriteBatchSize", int32(100000),
				"localTime", time.Now(),
				"connectionId", int32(42),
				"minWireVersion", common.MinWireVersion,
				"maxWireVersion", common.MaxWireVersion,
				"readOnly", false,
				"ok", float64(1),
			))},
		}, nil
	}

	return nil, commonerrors.NewCommandErrorMsgWithArgument(
		commonerrors.ErrNotImplemented,
		fmt.Sprintf("CmdQuery: unhandled command %q for collection %q", cmd, collection),
		"OpQuery: "+cmd,
	)
}