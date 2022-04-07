package broker

import (
	"github.com/uesleicarvalhoo/go-auth-service/internal/infra/config"
	"github.com/uesleicarvalhoo/go-auth-service/internal/schemas"
)

type Config = config.BrokerConfig

type Streamer interface {
	Start(eventChannel <-chan schemas.Event)
	End()
}
