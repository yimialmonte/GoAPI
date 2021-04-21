package main

import (
	"net/http"
	"strings"

	"github.com/asdine/storm"
	"github.com/labstack/echo"
	"github.com/yimialmonte/GoAPI/cache"
	"github.com/yimialmonte/GoAPI/user"
	"gopkg.in/mgo.v2/bson"
)

type jsonResponse map[string]interface{}

func usersOptios(c echo.Context) error {
	methods := []string{http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodHead}
	c.Response().Header().Set("Allow", strings.Join(methods, ","))
	return c.NoContent(http.StatusOK)
}

func userOptios(c echo.Context) error {
	methods := []string{http.MethodGet, http.MethodPut, http.MethodOptions, http.MethodHead, http.MethodDelete, http.MethodPatch}
	c.Response().Header().Set("Allow", strings.Join(methods, ","))
	return c.NoContent(http.StatusOK)
}

func usersGetAll(c echo.Context) error {
	if cache.Serve(c.Response(), c.Request()) {
		return nil
	}
	users, err := user.All()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if c.Request().Method == http.MethodHead {
		return c.NoContent(http.StatusOK)
	}
	c.Response().Writer = cache.NewWriter(c.Response(), c.Request())
	return c.JSON(http.StatusOK, jsonResponse{"users": users})
}

func usersGetOne(c echo.Context) error {
	if cache.Serve(c.Response(), c.Request()) {
		return nil
	}

	if !bson.IsObjectIdHex(c.Param("id")) {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	id := bson.ObjectIdHex(c.Param("id"))

	u, err := user.One(id)
	if err != nil {
		if err == storm.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}

		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if c.Request().Method == http.MethodHead {
		return c.NoContent(http.StatusOK)
	}

	c.Response().Writer = cache.NewWriter(c.Response(), c.Request())
	return c.JSON(http.StatusOK, jsonResponse{"user": u})
}

func usersPostOne(c echo.Context) error {
	u := new(user.User)
	err := c.Bind(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	u.ID = bson.NewObjectId()
	err = u.Save()
	if err != nil {
		if err == user.ErrRecordInvalid {
			return echo.NewHTTPError(http.StatusBadRequest)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	cache.Drop("/users")
	c.Response().Header().Set("Location", "/users/"+u.ID.Hex())
	return c.NoContent(http.StatusCreated)
}

func usersPutOne(c echo.Context) error {
	u := new(user.User)
	err := c.Bind(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !bson.IsObjectIdHex(c.Param("id")) {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	id := bson.ObjectIdHex(c.Param("id"))

	u.ID = id
	err = u.Save()
	if err != nil {
		if err == user.ErrRecordInvalid {
			return echo.NewHTTPError(http.StatusBadRequest)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	cache.Drop("/users")
	c.Response().Writer = cache.NewWriter(c.Response(), c.Request())
	return c.JSON(http.StatusOK, jsonResponse{"user": u})
}

func usersPatchOne(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	id := bson.ObjectIdHex(c.Param("id"))

	u, err := user.One(id)
	if err != nil {
		if err == storm.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	err = c.Bind(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	u.ID = id
	err = u.Save()
	if err != nil {
		if err == user.ErrRecordInvalid {
			return echo.NewHTTPError(http.StatusBadRequest)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	cache.Drop("/users")
	c.Response().Writer = cache.NewWriter(c.Response(), c.Request())
	return c.JSON(http.StatusOK, jsonResponse{"user": u})
}

func usersDeleteOne(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	id := bson.ObjectIdHex(c.Param("id"))

	err := user.Delete(id)
	if err != nil {
		if err == storm.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	cache.Drop("/users")
	cache.Drop(cache.MakeResource(c.Request()))
	return c.NoContent(http.StatusOK)
}

func root(c echo.Context) error {
	return c.String(http.StatusOK, "Running API V1")
}

func main() {
	e := echo.New()

	e.GET("/", root)

	u := e.Group("/users")

	u.OPTIONS("", usersOptios)
	u.HEAD("", usersGetAll)
	u.GET("", usersGetAll)
	u.POST("", usersPostOne)

	uid := u.Group("/:id")

	uid.OPTIONS("", userOptios)
	uid.GET("", usersGetOne)
	uid.PUT("", usersPutOne)
	uid.PATCH("", usersPatchOne)
	uid.DELETE("", usersDeleteOne)

	e.Start(":12345")
}
