package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/eclipse-xfsc/ssi-vdr-core/types"
	"github.com/gin-gonic/gin"
)

type IpfsResponse struct {
	Identifier *types.DataIdentifier `json:"identifier"`
	Data       []byte                `json:"data"`
}

type IpfsListResponse struct {
	Identifiers []*types.DataIdentifier `json:"identifiers"`
}

type IpfsErrorResponse struct {
	Error string `json:"error"`
}

func WrapHandler(f func(ctx *gin.Context, env Env) ([]byte, *types.DataIdentifier, error), env Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ipfsError IpfsError
		data, iden, err := f(c, env)
		if err == nil {
			c.JSON(http.StatusOK, IpfsResponse{
				Identifier: iden,
				Data:       data,
			})
		} else if errors.As(err, &ipfsError) {
			logger.Error(ipfsError, "")
			c.JSON(ipfsError.Code, IpfsErrorResponse{
				Error: ipfsError.Error(),
			})
		} else {
			err = fmt.Errorf("api call failed with error: %t", err)
			logger.Error(err, "")
			c.JSON(http.StatusInternalServerError, IpfsErrorResponse{
				Error: err.Error(),
			})
		}
	}
}

// GetDocument godoc
//
//	@Summary		Get document from ipfs
//	@Description	get document
//	@Tags			docs
//	@Produce		json
//	@Param			tenantId	path		string	true	"id of host tenant"
//	@Param			id			path		string	true	"CID of stored document"
//	@Success		200			{object}	IpfsResponse
//	@Failure		400			{object}	IpfsErrorResponse
//	@Failure		404			{object}	IpfsErrorResponse
//	@Failure		500			{object}	IpfsErrorResponse
//	@Router			/{id} [get]
func GetDocument(ctx *gin.Context, env Env) ([]byte, *types.DataIdentifier, error) {
	identifier, err := getId(ctx)
	if err != nil {
		return nil, nil, err
	}
	res, err := env.Ipfs().Get(identifier)
	if errors.Is(err, types.DataIdentifierNotFound) {
		return nil, nil, IpfsError{
			Err:  err,
			Msg:  fmt.Sprintf("ipfs could not find a document %s", identifier.Value),
			Code: http.StatusNotFound,
		}
	} else if err != nil {
		return nil, nil, err
	} else {
		return res.Data, identifier, nil
	}
}

// GetDocuments godoc
//
//	@Summary		Get documents' identifiers from ipfs
//	@Description	list documents
//	@Tags			docs
//	@Produce		json
//	@Param			tenantId	path		string	true	"id of host tenant"
//	@Success		200			{object}	IpfsListResponse
//	@Failure		500			{object}	IpfsErrorResponse
//	@Router			/list [get]
func GetDocuments(env Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := env.Ipfs().List()
		if err != nil {
			logger.Error(err, "")
			c.JSON(http.StatusInternalServerError, IpfsErrorResponse{
				Error: err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, IpfsListResponse{Identifiers: res})
		}
	}
}

// CreateDocument godoc
//
//	@Summary		Create a new document in ipfs
//	@Description	create document
//	@Tags			docs
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			tenantId	path		string	true	"id of host tenant"
//	@Param			document	formData	file	true	"Document data to be created"
//	@Success		200			{object}	IpfsResponse
//	@Failure		500			{object}	IpfsErrorResponse
//	@Router			/create [post]
func CreateDocument(ctx *gin.Context, env Env) ([]byte, *types.DataIdentifier, error) {
	identifier := types.DataIdentifier{}
	data := ctx.Request.Body
	res, err := env.Ipfs().Put(&identifier, data)
	if err != nil {
		return nil, nil, err
	} else {
		return nil, res, nil
	}
}

// UpdateDocument godoc
//
//	@Summary		Update a document in ipfs
//	@Description	update document
//	@Tags			docs
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			tenantId	path		string	true	"id of host tenant"
//	@Param			id			path		string	true	"CID of stored document"
//	@Param			document	formData	file	true	"Document data to be updated with"
//	@Success		200			{object}	IpfsResponse
//	@Failure		400			{object}	IpfsErrorResponse
//	@Failure		500			{object}	IpfsErrorResponse
//	@Router			/{id}/update [put]
func UpdateDocument(ctx *gin.Context, env Env) ([]byte, *types.DataIdentifier, error) {
	identifier, err := getId(ctx)
	if err != nil {
		return nil, nil, err
	}
	data := ctx.Request.Body
	res, err := env.Ipfs().Update(identifier, data)
	if err != nil {
		return nil, nil, err
	} else {
		return nil, res, nil
	}
}

// DeleteDocument godoc
//
//	@Summary		Delete a document from ipfs
//	@Description	delete document
//	@Tags			docs
//	@Produce		json
//	@Param			tenantId	path		string	true	"id of host tenant"
//	@Param			id			path		string	true	"CID of document to be deleted"
//	@Success		200			{object}	IpfsResponse
//	@Failure		400			{object}	IpfsErrorResponse
//	@Failure		500			{object}	IpfsErrorResponse
//	@Router			/{id} [delete]
func DeleteDocument(ctx *gin.Context, env Env) ([]byte, *types.DataIdentifier, error) {
	identifier, err := getId(ctx)
	if err != nil {
		return nil, nil, err
	}
	err = env.Ipfs().Delete(identifier)
	if err != nil {
		return nil, nil, err
	} else {
		return nil, identifier, nil
	}
}

func getId(ctx *gin.Context) (*types.DataIdentifier, error) {
	id := ctx.Param("id")
	if id == "" {
		return nil, IpfsError{
			Err:  fmt.Errorf("no id provided"),
			Msg:  "id is a required url parameter",
			Code: http.StatusBadRequest,
		}
	}
	identifier := types.DataIdentifier{
		Format: "cid",
		Value:  id,
	}
	return &identifier, nil
}
