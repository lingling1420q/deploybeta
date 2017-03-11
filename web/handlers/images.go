package handlers

import (
	"net/http"

	"github.com/labstack/echo"

	appModel "github.com/jysperm/deploying/lib/models/app"
	versionModel "github.com/jysperm/deploying/lib/models/version"
	. "github.com/jysperm/deploying/web/handlers/helpers"
)

func CreateImage(ctx echo.Context) error {
	app, err := appModel.FindByName(ctx.Param("name"))
	if err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	version, err := versionModel.CreateVersion(app)
	if err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, NewVersionResponse(&version))
}

// TODO:
func DeleteImage(ctx echo.Context) error {
	return nil
}
