// Copyright 2021 The Kubeswitch authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package store

import (
	"context"
	"fmt"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/danielfoehrkn/kubeswitch/types"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func NewEKSStore(store types.KubeconfigStore) (*EKSStore, error) {
	eksStoreConfig := &types.StoreConfigEKS{}
	if store.Config != nil {
		buf, err := yaml.Marshal(store.Config)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(buf, eksStoreConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal eks config: %w", err)
		}
	}

	return &EKSStore{
		Logger:          logrus.New().WithField("store", types.StoreKindEKS),
		KubeconfigStore: store,
		Config:          eksStoreConfig,
	}, nil
}

func (s *EKSStore) InitializeEKSStore() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	optFns := []func(*awsconfig.LoadOptions) error{}

	if s.Config != nil {
		if s.Config.Region != nil {
			optFns = append(optFns, awsconfig.WithRegion(*s.Config.Region))
		}
		if s.Config.Profile != nil {
			optFns = append(optFns, awsconfig.WithSharedConfigProfile(*s.Config.Profile))
		}
	}

	cfg, err := awsconfig.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		return err
	}

	s.Client = awseks.NewFromConfig(cfg)

	return nil
}

func (s *EKSStore) IsInitialized() bool {
	return s.Client != nil && s.Config != nil
}

func (s *EKSStore) GetID() string {
	id := "default"
	if s.KubeconfigStore.ID != nil {
		id = *s.KubeconfigStore.ID
	}
	return fmt.Sprintf("%s.%s", types.StoreKindEKS, id)
}

func (s *EKSStore) GetKind() types.StoreKind {
	return types.StoreKindEKS
}

func (s *EKSStore) GetContextPrefix(path string) string {
	if s.GetStoreConfig().ShowPrefix != nil && !*s.GetStoreConfig().ShowPrefix {
		return ""
	}

	return fmt.Sprintf("%s/%s/%s", s.GetID(), *s.Config.Profile, *s.Config.Region)
}

func (s *EKSStore) VerifyKubeconfigPaths() error {
	return nil
}

func (s *EKSStore) StartSearch(channel chan SearchResult) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.InitializeEKSStore(); err != nil {
		err := fmt.Errorf("failed to initialize store. This is most likely a problem with your provided kubeconfig: %v", err)
		channel <- SearchResult{
			Error: err,
		}
		return
	}

	opts := &awseks.ListClustersInput{}
	for {
		clusters, err := s.Client.ListClusters(ctx, opts)
		if err != nil {
			channel <- SearchResult{
				Error: err,
			}
			return
		}

		for _, clusterName := range clusters.Clusters {
			channel <- SearchResult{
				KubeconfigPath: clusterName,
				Error:          nil,
			}
		}

		if clusters.NextToken != nil {
			opts.NextToken = clusters.NextToken
		} else {
			break
		}
	}
}

func (s *EKSStore) GetKubeconfigForPath(path string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if !s.IsInitialized() {
		if err := s.InitializeEKSStore(); err != nil {
			return nil, fmt.Errorf("failed to initialize EKS store: %w", err)
		}
	}

	cluster, err := s.Client.DescribeCluster(ctx, &awseks.DescribeClusterInput{Name: &path})
	if err != nil {
		return nil, err
	}

	kubeconfig := &types.KubeConfig{
		TypeMeta: types.TypeMeta{
			APIVersion: "v1",
			Kind:       "Config",
		},
		Clusters: []types.KubeCluster{{
			Name: path,
			Cluster: types.Cluster{
				CertificateAuthorityData: *cluster.Cluster.CertificateAuthority.Data,
				Server:                   *cluster.Cluster.Endpoint,
			},
		}},
		CurrentContext: path,
		Contexts: []types.KubeContext{
			{
				Name: path,
				Context: types.Context{
					Cluster: path,
					User:    path,
				},
			},
		},
	}

	if cluster.Cluster.Identity != nil {
		kubeconfig.Users = []types.KubeUser{
			{
				Name: path,
				User: types.User{
					ExecProvider: &types.ExecProvider{
						APIVersion: "client.authentication.k8s.io/v1alpha1",
						Command:    "aws",
						Args:       []string{"--region", *s.Config.Region, "eks", "get-token", "--cluster-name", path},
						Env: []types.EnvMap{
							{Name: "AWS_PROFILE", Value: *s.Config.Profile},
						},
					},
				},
			},
		}
	}

	bytes, err := yaml.Marshal(kubeconfig)

	return bytes, err
}

func (s *EKSStore) GetLogger() *logrus.Entry {
	return s.Logger
}

func (s *EKSStore) GetStoreConfig() types.KubeconfigStore {
	return s.KubeconfigStore
}
