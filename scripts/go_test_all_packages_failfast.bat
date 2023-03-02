cd /d "%~dp0.."

go test -failfast -v github.com/kuoss/lethe                           || exit /b
go test -failfast -v github.com/kuoss/lethe/cli                       || exit /b
go test -failfast -v github.com/kuoss/lethe/cli/cmd                   || exit /b
go test -failfast -v github.com/kuoss/lethe/cli/cmd/list              || exit /b
go test -failfast -v github.com/kuoss/lethe/cli/cmd/task              || exit /b
go test -failfast -v github.com/kuoss/lethe/cli/util                  || exit /b
go test -failfast -v github.com/kuoss/lethe/clock                     || exit /b
go test -failfast -v github.com/kuoss/lethe/config                    || exit /b
go test -failfast -v github.com/kuoss/lethe/handlers                  || exit /b
go test -failfast -v github.com/kuoss/lethe/letheql                   || exit /b
go test -failfast -v github.com/kuoss/lethe/logs                      || exit /b
go test -failfast -v github.com/kuoss/lethe/storage                   || exit /b
go test -failfast -v github.com/kuoss/lethe/storage/driver            || exit /b
go test -failfast -v github.com/kuoss/lethe/storage/driver/factory    || exit /b
go test -failfast -v github.com/kuoss/lethe/storage/driver/filesystem || exit /b
go test -failfast -v github.com/kuoss/lethe/testutil                  || exit /b
go test -failfast -v github.com/kuoss/lethe/util                      || exit /b
