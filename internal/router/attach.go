// GoToSocial
// Copyright (C) GoToSocial Authors admin@gotosocial.org
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package router

import (
	"slices"
	"unsafe"

	"github.com/gin-gonic/gin"
)

func (r *Router) AttachGlobalMiddleware(handlers ...gin.HandlerFunc) gin.IRoutes {
	handlers = slices.DeleteFunc(handlers, func(h gin.HandlerFunc) bool { return h == nil })
	return r.engine.Use(handlers...)
}

func (r *Router) AttachNoRouteHandler(handler gin.HandlerFunc) {
	r.engine.NoRoute(handler)
}

func (r *Router) AttachGroup(relativePath string, handlers ...gin.HandlerFunc) *RouterGroup {
	handlers = slices.DeleteFunc(handlers, func(h gin.HandlerFunc) bool { return h == nil })
	return (*RouterGroup)(unsafe.Pointer(r.engine.Group(relativePath, handlers...)))
}

func (r *Router) AttachHandler(method string, path string, handler gin.HandlerFunc) {
	r.engine.Handle(method, path, handler)
}

// RouterGroup wraps a gin.RouterGroup to
// check for and skip appending nil handlers.
type RouterGroup struct{ gin.RouterGroup } // nolint:revive

func (g *RouterGroup) Handle(httpMethod string, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	handlers = slices.DeleteFunc(handlers, func(h gin.HandlerFunc) bool { return h == nil })
	return g.RouterGroup.Handle(httpMethod, relativePath, handlers...)
}

func (g *RouterGroup) Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	middleware = slices.DeleteFunc(middleware, func(h gin.HandlerFunc) bool { return h == nil })
	return g.RouterGroup.Use(middleware...)
}
