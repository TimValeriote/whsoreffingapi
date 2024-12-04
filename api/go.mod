module whosreffing-api/api

go 1.19

replace whosreffing-api/utils => ../utils

replace whosreffing-api/config => ../config

replace whosreffing-api/core => ../core

replace whosreffing-api/apibuilder => ../apibuilder

replace whosreffing-api/controllers => ../controllers

require (
	github.com/Tomasen/realip v0.0.0-20180522021738-f0c99a92ddce
	github.com/go-sql-driver/mysql v1.7.0
	github.com/gorilla/context v1.1.1
	github.com/jasonlvhit/gocron v0.0.1
	github.com/julienschmidt/httprouter v1.3.0
	github.com/justinas/alice v1.2.0
	github.com/rs/cors v1.8.3
	github.com/sirupsen/logrus v1.9.0
	gopkg.in/guregu/null.v3 v3.5.0
	whosreffing-api/apibuilder v0.0.0-00010101000000-000000000000
	whosreffing-api/config v0.0.0-00010101000000-000000000000
	whosreffing-api/controllers v0.0.0-00010101000000-000000000000
	whosreffing-api/core v0.0.0-00010101000000-000000000000
	whosreffing-api/utils v0.0.0-00010101000000-000000000000
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)
