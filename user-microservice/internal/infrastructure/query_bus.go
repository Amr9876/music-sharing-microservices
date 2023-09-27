package infrastructure

import "music-sharing/user-microservice/internal/interfaces"

type QueryBus struct{}

var queryBus = &QueryBus{}

func (q *QueryBus) Send(query interfaces.QueryHandler) (interface{}, error) {
	return query.Handle()
}

func GetQueryBus() *QueryBus {
	return queryBus
}
