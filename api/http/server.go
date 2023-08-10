// Copyright 2023 Democratized Data Foundation
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package http

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sourcenetwork/defradb/api"
	"github.com/sourcenetwork/defradb/api/core"
	"github.com/sourcenetwork/defradb/client"
	"github.com/sourcenetwork/defradb/net"
)

type Server struct {
	api api.API
}

func NewServer(db client.DB, node *net.Node) *Server {
	server := &Server{core.New(db, node)}

	router := gin.Default()
	api := router.Group("/api/v0")

	backup := api.Group("/backup")
	backup.POST("/export", server.ExportBackup)
	backup.POST("/import", server.ImportBackup)

	debug := api.Group("/debug")
	debug.GET("/dump", server.Dump)
	debug.GET("/block/:cid", server.GetBlock)

	index := api.Group("/index")
	index.POST("/", server.CreateIndex)
	index.GET("/", server.ListIndexes)
	index.GET("/:name", server.ListIndexesForCollection)

	schema := api.Group("/schema")
	schema.GET("/", server.ListSchemas)
	schema.POST("/", server.LoadSchema)
	schema.PATCH("/", server.PatchSchema)

	migration := schema.Group("/migration")
	migration.GET("/", server.GetMigration)
	migration.POST("/", server.SetMigration)

	peer := api.Group("/peer")
	peer.GET("/info", server.PeerInfo)

	replicators := peer.Group("/replicators")
	replicators.GET("/replicators", server.ListReplicators)
	replicators.POST("/replicators", server.SetReplicator)
	replicators.DELETE("/replicators", server.DeleteReplicator)

	collections := peer.Group("/collections")
	collections.GET("/collections", server.ListPeerCollections)
	collections.POST("/collections/:id", server.AddPeerCollection)
	collections.DELETE("/collections/:id", server.RemovePeerCollection)

	return server
}

func (s *Server) ExportBackup(c *gin.Context) {
	var config client.BackupConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.api.ExportBackup(c.Request.Context(), config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (s *Server) ImportBackup(c *gin.Context) {
	var config client.BackupConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.api.ImportBackup(c.Request.Context(), config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (s *Server) Dump(c *gin.Context) {
	if err := s.api.Dump(c.Request.Context()); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (s *Server) GetBlock(c *gin.Context) {
	block, err := s.api.GetBlock(c.Request.Context(), c.Param("cid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, block)
}

func (s *Server) CreateIndex(c *gin.Context) {
	var req api.CreateIndexRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	index, err := s.api.CreateIndex(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, index)
}

func (s *Server) DropIndex(c *gin.Context) {
	var req api.DropIndexRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.api.DropIndex(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (s *Server) ListIndexes(c *gin.Context) {
	indexes, err := s.api.ListIndexes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, indexes)
}

func (s *Server) ListIndexesForCollection(c *gin.Context) {
	indexes, err := s.api.ListIndexesForCollection(c.Request.Context(), c.Param("name"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, indexes)
}

func (s *Server) PeerInfo(c *gin.Context) {
	info, err := s.api.PeerInfo(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

func (s *Server) SetReplicator(c *gin.Context) {
	var req api.ReplicatorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.api.SetReplicator(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (s *Server) DeleteReplicator(c *gin.Context) {
	var req api.ReplicatorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.api.DeleteReplicator(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (s *Server) ListReplicators(c *gin.Context) {
	reps, err := s.api.ListReplicators(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reps)
}

func (s *Server) AddPeerCollection(c *gin.Context) {
	if err := s.api.AddPeerCollection(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (s *Server) RemovePeerCollection(c *gin.Context) {
	if err := s.api.RemovePeerCollection(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (s *Server) ListPeerCollections(c *gin.Context) {
	cols, err := s.api.ListPeerCollections(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cols)
}

func (s *Server) LoadSchema(c *gin.Context) {
	schema, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cols, err := s.api.LoadSchema(c.Request.Context(), string(schema))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cols)
}

func (s *Server) PatchSchema(c *gin.Context) {
	patch, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cols, err := s.api.PatchSchema(c.Request.Context(), string(patch))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cols)
}

func (s *Server) ListSchemas(c *gin.Context) {
	cols, err := s.api.ListSchemas(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cols)
}

func (s *Server) SetMigration(c *gin.Context) {
	var req client.LensConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.api.SetMigration(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (s *Server) GetMigration(c *gin.Context) {
	cfgs, err := s.api.GetMigration(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cfgs)
}
