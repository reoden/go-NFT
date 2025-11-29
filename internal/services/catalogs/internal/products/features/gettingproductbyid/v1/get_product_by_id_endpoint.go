package v1

import (
    "net/http"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/reoden/go-NFT/catalogs/internal/products/dtos/v1/fxparams"
    "github.com/reoden/go-NFT/catalogs/internal/products/features/gettingproductbyid/v1/dtos"
    "github.com/reoden/go-NFT/pkg/core/web/route"
    customErrors "github.com/reoden/go-NFT/pkg/http/httperrors/customerrors"

    "emperror.dev/errors"
    "github.com/labstack/echo/v4"
    "github.com/mehdihadeli/go-mediatr"
)

type getProductByIdEndpoint struct {
    fxparams.ProductRouteParams
}

func NewGetProductByIdEndpoint(
    params fxparams.ProductRouteParams,
) route.Endpoint {
    return &getProductByIdEndpoint{ProductRouteParams: params}
}

func (ep *getProductByIdEndpoint) MapEndpoint() {
    ep.ProductsGroup.GET("/:id", ep.handler())
}

// GetProductByID
// @Tags Products
// @Summary Get product by id
// @Description Get product by id
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dtos.GetProductByIdResponseDto
// @Security BearerAuth
// @Router /api/v1/products/{id} [get]
func (ep *getProductByIdEndpoint) handler() echo.HandlerFunc {
    return func(c echo.Context) error {
        // 安全地检查 user 上下文
        if user := c.Get("user"); user != nil {
            if token, ok := user.(*jwt.Token); ok {
                ep.Logger.Printf("GetProductByIdEndpoint authtoken = %+v", token)
                claims, ok := token.Claims.(jwt.MapClaims)
                if !ok {
                    return errors.New("failed to cast claims as jwt.MapClaims")
                }

                ep.Logger.Printf("[DEBUG] claims = %+v", claims)
            } else {
                ep.Logger.Printf("GetProductByIdEndpoint authtoken = %+v", user)
                claims, ok := token.Claims.(jwt.MapClaims)
                if !ok {
                    return errors.New("failed to cast claims as jwt.MapClaims")
                }

                ep.Logger.Printf("[DEBUG] claims = %+v", claims)
            }
        } else {
            ep.Logger.Printf("No user context found (likely skipped by authSkipper)")
        }

        ctx := c.Request().Context()

        request := &dtos.GetProductByIdRequestDto{}
        if err := c.Bind(request); err != nil {
            badRequestErr := customErrors.NewBadRequestErrorWrap(
                err,
                "error in the binding request",
            )

            return badRequestErr
        }

        query, err := NewGetProductByIdWithValidation(request.ProductId)
        if err != nil {
            return err
        }

        queryResult, err := mediatr.Send[*GetProductById, *dtos.GetProductByIdResponseDto](
            ctx,
            query,
        )
        if err != nil {
            return errors.WithMessage(
                err,
                "error in sending GetProductById",
            )
        }

        userID := "123"

        // 创建 token
        claims := jwt.MapClaims{
            "sub":  userID,
            "exp":  time.Now().Add(time.Hour * 24).Unix(), // 过期时间 24h
            "iat":  time.Now().Unix(),
            "role": "admin", // 可以加自定义字段
        }

        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

        // 签名生成完整 token 字符串
        t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
        if err != nil {
            return c.JSON(http.StatusInternalServerError, map[string]string{
                "error": "failed to generate token",
            })
        }

        ep.Logger.Printf("[DEBUG] token = %+v", t)

        return c.JSON(http.StatusOK, queryResult)
    }
}
