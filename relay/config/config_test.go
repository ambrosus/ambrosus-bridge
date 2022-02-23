package config

import (
	"os"
	"reflect"
	"testing"
)

func TestConfig_Init(t *testing.T) {
	type env struct {
		configPath    string
		ambPrivateKey string
		ethPrivateKey string
	}

	type args struct{ env env }

	setEnv := func(env env) {
		os.Setenv("CONFIG_PATH", env.configPath)
		os.Setenv("AMB_PRIVATE_KEY", env.ambPrivateKey)
		os.Setenv("ETH_PRIVATE_KEY", env.ambPrivateKey)
	}

	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "OK",
			args: args{env: env{
				configPath:    "fixtures/main",
				ambPrivateKey: "34d8e83fca265e9ab5bcc1094fa64e98692375bf8980d066a9edcf4953f0f2f5",
				ethPrivateKey: "34d8e83fca265e9ab5bcc1094fa64e98692375bf8980d066a9edcf4953f0f2f5",
			}},
			want: &Config{
				AMB: AMBConfig{
					Network: Network{
						URL:          "wss://network.ambrosus-dev.io",
						ContractAddr: "",
					},
					VSContractAddr: "",
				},
				ETH: ETHConfig{
					Network: Network{
						URL:          "wss://rinkeby.infura.io/ws/v3/01117e8ede8e4f36801a6a838b24f36c",
						ContractAddr: "",
					},
					EthashPath: "./",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setEnv(tt.args.env)

			got, err := Init()
			if (err != nil) != tt.wantErr {
				t.Errorf("error initialize config: %s", err.Error())
			}

			tt.want.AMB.PrivateKey = got.AMB.PrivateKey
			tt.want.ETH.PrivateKey = got.ETH.PrivateKey

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("error config are not similar: %s", err.Error())
			}
		})
	}
}
