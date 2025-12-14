package authcertification

import (
    "context"
    "fmt"

    "github.com/go-resty/resty/v2"
    "github.com/reoden/go-NFT/pkg/logger"
)

type AuthCertificationService interface {
    Auth(ctx context.Context, realName, idCardNo string) (bool, error)
}

type MockAuthCertificationServiceImpl struct {
}

func (impl *MockAuthCertificationServiceImpl) Auth(_ context.Context, realName, idCardNo string) (bool, error) {
    return true, nil
}

type AuthCertificationServiceImpl struct {
    host    string
    path    string
    appCode string
    client  *resty.Client
    logger  logger.Logger
}

func (impl *AuthCertificationServiceImpl) Auth(ctx context.Context, realName, idCardNo string) (bool, error) {
    headers := map[string]string{
        "Authorization": fmt.Sprintf("APPCODE %s", impl.appCode),
        "Content-Type":  "application/x-www-form-urlencoded; charset=UTF-8",
    }

    body := map[string]string{
        "id_number": idCardNo,
        "name":      realName,
    }

    resp, err := impl.client.R().
        SetHeaders(headers).
        SetFormData(body).
        SetContext(ctx).
        Post(fmt.Sprintf("%s%s", impl.host, impl.path))

    if err != nil {
        impl.logger.Error("auth certification request error", err)
        return false, err
    }

    var result map[string]interface{}
    if err := impl.client.JSONUnmarshal(resp.Body(), &result); err != nil {
        impl.logger.Error("failed to unmarshal auth response", err)
        return false, err
    }

    impl.logger.Infow("auth result", logger.Fields{"result": result})

    if state, ok := result["state"]; ok {
        if stateVal, ok := state.(int); ok {
            if stateVal == 1 {
                return true, nil
            }
        }
    }

    return false, nil
}

// NewAuthCertificationService create new auth certification service
func NewAuthCertificationService(cfg *AuthCertificationOptions, client *resty.Client, logger logger.Logger) AuthCertificationService {
    if cfg.Host == "" ||
        cfg.Path == "" ||
        cfg.AppCode == "" {
        return &MockAuthCertificationServiceImpl{}
    }

    return &AuthCertificationServiceImpl{
        host:    cfg.Host,
        path:    cfg.Path,
        appCode: cfg.AppCode,
        client:  client,
        logger:  logger,
    }
}
