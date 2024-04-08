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

package mutator

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
)

type MutateWorker struct {
	cfg         SdkConfig.SdkInfo
	chSubscribe chan []byte
	db          *db.MongoProvider
	logger      interfaces.Logger
	sdk         interfaces.Sdk
	svc         config.ServiceInfo
}

func NewMutateWorker(sdk interfaces.Sdk, chSub chan []byte, cfg SdkConfig.SdkInfo, transit config.ServiceInfo, db *db.MongoProvider, logger interfaces.Logger) MutateWorker {
	return MutateWorker{
		cfg:         cfg,
		chSubscribe: chSub,
		db:          db,
		logger:      logger,
		sdk:         sdk,
		svc:         transit,
	}
}

func (m *MutateWorker) BootstrapHandler(ctx context.Context, wg *sync.WaitGroup) bool {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			msg, ok := <-m.chSubscribe
			if ok {
				//Create new/transformed data
				data, err := models.NewSampleData(m.cfg.Signature.PrivateKey)
				if err != nil {
					m.logger.Error(err.Error())
					continue
				}

				//Save the data
				err = m.db.Save(ctx, data)
				if err != nil {
					m.logger.Error(err.Error())
					continue
				}

				b, _ := json.Marshal(data)
				//Annotate the data, linking new to old
				m.sdk.Mutate(ctx, msg, b)
				//Send data to the next service
				tr := &http.Transport{
					DisableKeepAlives: true,
					TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
				}
				client := &http.Client{Transport: tr}
				resp, err := client.Post(m.svc.Uri()+"/data", config.HeaderValueJson, bytes.NewBuffer(b))
				if err != nil {
					m.logger.Error(err.Error())
				} else {
					resp.Body.Close()
				}
			} else { //channel has been closed. End goroutine.
				m.logger.Write(slog.LevelInfo, "mutator::chSubscribe closed, exiting")
				return
			}
		}
	}()

	wg.Add(1)
	go func() { // Graceful shutdown
		defer wg.Done()

		<-ctx.Done()
		m.logger.Write(slog.LevelInfo, "shutdown received")
	}()
	return true
}
