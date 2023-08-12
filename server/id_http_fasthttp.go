package main

import (
	"encoding/json"
	"strconv"

	"github.com/threeq/xid"
	"github.com/valyala/fasthttp"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func newIDFastHttp(basePath string) fasthttp.RequestHandler {
	urlHandlers := make(map[string]func(ctx *fasthttp.RequestCtx))
	for _, t := range xid.GetIdTypes() {
		switch t {
		case xid.ID_Snake:
			urlHandlers[basePath+"/snake"] = func(ctx *fasthttp.RequestCtx) {
				idHandlerFunc(ctx, xid.ID_Snake, false)
			}
			urlHandlers[basePath+"/snake/batch"] = func(ctx *fasthttp.RequestCtx) {
				idHandlerFunc(ctx, xid.ID_Snake, true)
			}
		case xid.ID_14:
			urlHandlers[basePath+"/14"] = func(ctx *fasthttp.RequestCtx) {
				idHandlerFunc(ctx, xid.ID_14, false)
			}
			urlHandlers[basePath+"/14/batch"] = func(ctx *fasthttp.RequestCtx) {
				idHandlerFunc(ctx, xid.ID_14, true)
			}
		}
	}

	m := func(ctx *fasthttp.RequestCtx) {
		if h, ok := urlHandlers[string(ctx.Path())]; ok {
			h(ctx)
		} else {
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}
	return m
}

func idHandlerFunc(ctx *fasthttp.RequestCtx, idtype xid.IdType, batch bool) {
	gen := string(ctx.QueryArgs().Peek("gen"))
	if !batch {
		id := xid.GetID(idtype, gen)
		ctx.SetBodyString(strconv.FormatInt(id, 10))
		return
	}

	numStr := string(ctx.QueryArgs().Peek("num"))

	var res Response

	num, err := strconv.Atoi(numStr)
	if err != nil {
		res.Code = 422
		res.Message = err.Error()
		data, _ := json.Marshal(res)
		ctx.SetBody(data)
		return
	}
	if num < 1 {
		res.Code = 422
		res.Message = "num最小值应为1"
		data, _ := json.Marshal(res)
		ctx.SetBody(data)
		return
	}

	ids := xid.GetIDS(idtype, gen, num)

	res.Code = 200
	res.Message = "success"
	res.Data = map[string]interface{}{
		"ids": ids,
	}
	data, _ := json.Marshal(res)
	ctx.SetBody(data)
}
