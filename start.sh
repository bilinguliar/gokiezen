#!/bin/sh

/opt/gokiezen/gokiezen \
	--port $PORT \
	--token $TOKEN \
	--redis_host $REDIS_HOST \
	--redis_port $REDIS_PORT \
	--redis_pool_size $REDIS_POOL_SIZE \
	--redis_conn_type $REDIS_CONNECTION_TYPE	
