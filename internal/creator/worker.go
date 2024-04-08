/*******************************************************************************
 * Copyright 2024 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package creator

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	SdkConfig "github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
	"github.com/project-alvarium/ones-demo-2024/internal/config"
	"github.com/project-alvarium/ones-demo-2024/internal/db"
	"github.com/project-alvarium/ones-demo-2024/internal/models"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type CreateWorker struct {
	cfg    SdkConfig.SdkInfo
	db     *db.MongoProvider
	logger interfaces.Logger
	sdk    interfaces.Sdk
	svc    config.ServiceInfo
}

func NewCreateWorker(sdk interfaces.Sdk, cfg SdkConfig.SdkInfo, mutate config.ServiceInfo, db *db.MongoProvider, logger interfaces.Logger) CreateWorker {
	return CreateWorker{
		cfg:    cfg,
		db:     db,
		logger: logger,
		sdk:    sdk,
		svc:    mutate,
	}
}

func (c *CreateWorker) BootstrapHandler(ctx context.Context, wg *sync.WaitGroup) bool {
	cancelled := false
	wg.Add(1)
	go func() {
		defer wg.Done()

		for !cancelled {
			data, err := models.NewSampleData(c.cfg.Signature.PrivateKey)
			if err != nil {
				c.logger.Error(err.Error())
				continue
			}
			//Save the data
			err = c.db.Save(ctx, data)
			if err != nil {
				c.logger.Error(err.Error())
				continue
			}

			b, _ := json.Marshal(data)
			//Annotate the data
			c.sdk.Create(ctx, b)

			//Send data to the next service
			tr := &http.Transport{
				DisableKeepAlives: true,
				TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			}
			client := &http.Client{Transport: tr}
			resp, err := client.Post(c.svc.Uri()+"/data", config.HeaderValueJson, bytes.NewBuffer(b))
			if err != nil {
				c.logger.Error(err.Error())
			} else {
				//c.logger.Write(logging.DebugLevel,fmt.Sprintf("posted to mutator (%s/data): %s", c.svc.Uri(), string(b)))
				resp.Body.Close()
			}
			time.Sleep(1 * time.Second)
			//cancelled = true /* Uncomment if you just want to run through a single pass */
		}
		c.logger.Write(slog.LevelDebug, "cancel received")
	}()

	wg.Add(1)
	go func() { // Graceful shutdown
		defer wg.Done()

		<-ctx.Done()
		c.logger.Write(slog.LevelDebug, "shutdown received")
		cancelled = true
	}()
	return true
}
