package definitions

import (
	"context"

	"github.com/sarulabs/di"
)

func New(stopCtx context.Context) (di.Container, error) {
	builder, err := di.NewBuilder()
	if err != nil {
		return nil, err
	}

	if err := builder.Add([]di.Def{
		getConfigDef(),

		// database
		getDatabaseDef(stopCtx),

		// repos
		getServersRepoDef(),
		getUsersRepoDef(),

		// server
		getHTTPSrvDef(),

		// services
		getErrorsServiceDef(),
		getServersServiceDef(),
		getUsersServiceDef(),
		getOutlineServiceDef(),
		getParserServiceDef(),
	}...); err != nil {
		return nil, err
	}

	return builder.Build(), nil
}
