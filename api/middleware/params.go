package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/coinfinitygroup/cfapp/api/shared"
	"github.com/coinfinitygroup/cfapp/pkg/errdef"
	"github.com/coinfinitygroup/cfapp/pkg/httpio"

	"github.com/coinfinitygroup/cfapp/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

// ParamUint parses uint from route
// řeší parametry předané v URL objektu
func ParamUint(paramName string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// vytáhne context z requestu
			ctx := req.Context()
			// vytáhne z URL parametry podle předaného názvu parametru z argumentu
			code := chi.URLParam(req, paramName)
			// převede parametry do uint64 aby bylo možné je použít ve funkci with value k 
			param, err := strconv.ParseUint(code, 10, 64)
			if err != nil {
				errSet := errdef.NewErrData(paramName, "is not valid")
				// pokud mají parametry špatně formát odpoví ze servu error 
				httpio.WriteJSONErrSet(w, errSet)
				return
			}
			// obohatí kontext o loger
			ctx = logger.CtxFields(ctx, logrus.Fields{
				paramName: param,
			})
			// umožnuje sdílet data, vytvoří nový kontext v zavislosti na poskytnutém rodiči 
			// přidává do takto vytvořeného kontextu novou hodnotu na poskytnutý key, jsou tam data uložená jako klíč a hodnota 
			// které lze později vytváhnout a pracovat s nimi
			// ukládá tam klíč předaný do funkce jako argument prohnaný funkcí ParamCTx ze složky shared, hodnotou jsou parametry vytažené z
			// url převedené do formátu uint
			ctx = context.WithValue(ctx, shared.ParamCtx(paramName), uint(param))
			// píše headrs a value do odpovědi 
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}
