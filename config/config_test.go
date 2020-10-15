package config

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func Test_parseConfig(t *testing.T) {
	type args struct {
		path string
	}

	content := "[Server]\n" +
		"ListenAddress=\"0.0.0.0\"\n" +
		"ListenPort=8080\n" +
		"\n" +
		"[Endpoints]\n" +
		"create_users={ Enabled=true, Auth=false, AdminOnly=false }\n" +
		"manage_users={ Enabled=true, Auth=false, AdminOnly=false }\n" +
		"go={ Enabled=true, Auth=false, AdminOnly=false }\n" +
		"manage_paths={ Enabled=true, Auth=false, AdminOnly=false }\n" +
		"\n" +
		"[Cache]\n" +
		"Enabled=false\n" +
		"Type=\"\"\n" +
		"Address=\"\"\n" +
		"Port=0\n" +
		"\n" +
		"[Database]\n" +
		"Type=\"mysql\"\n" +
		"Address=\"localhost\"\n" +
		"Port=3306\n" +
		"SslMode=false\n" +
		"Protocol=\"tcp\"\n" +
		"Name=\"go_there_db\"\n" +
		"User=\"my_user\"\n" +
		"Password=\"superpassword\""

	tmpf, err := ioutil.TempFile(os.TempDir(), "go-there.conf")

	if err != nil {
		assert.Fail(t, err.Error())
	}

	defer func() {
		_ = os.Remove(tmpf.Name())
	}()

	if _, err := tmpf.Write([]byte(content)); err != nil {
		assert.Fail(t, err.Error())
	}

	if err := tmpf.Close(); err != nil {
		assert.Fail(t, err.Error())
	}

	tests := []struct {
		name    string
		args    args
		want    *Configuration
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				path: tmpf.Name(),
			},
			want: &Configuration{
				Server: Server{
					ListenAddress: "0.0.0.0",
					ListenPort:    8080,
				},
				Cache: Cache{
					Enabled: false,
					Type:    "",
					Address: "",
					Port:    0,
				},
				Database: Database{
					Type:     "mysql",
					Address:  "localhost",
					Port:     3306,
					SslMode:  false,
					Protocol: "tcp",
					Name:     "go_there_db",
					User:     "my_user",
					Password: "superpassword",
				},
				Endpoints: map[string]Endpoint{
					"create_users": {Enabled: true, Auth: false, AdminOnly: false},
					"manage_users": {Enabled: true, Auth: false, AdminOnly: false},
					"go":           {Enabled: true, Auth: false, AdminOnly: false},
					"manage_paths": {Enabled: true, Auth: false, AdminOnly: false},
				},
			},
			wantErr: false,
		},
		{
			name: "no_file",
			args: args{
				path: "/bad/path",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseConfig(tt.args.path)

			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInit(t *testing.T) {
	type args struct {
		path string
	}

	content := "[Server]\n" +
		"ListenAddress=\"0.0.0.0\"\n" +
		"ListenPort=8080\n" +
		"\n" +
		"[Endpoints]\n" +
		"create_users={ Enabled=true, Auth=false, AdminOnly=false }\n" +
		"manage_users={ Enabled=true, Auth=false, AdminOnly=false }\n" +
		"go={ Enabled=true, Auth=false, AdminOnly=false }\n" +
		"manage_paths={ Enabled=true, Auth=false, AdminOnly=false }\n" +
		"\n" +
		"[Cache]\n" +
		"Enabled=false\n" +
		"Type=\"\"\n" +
		"Address=\"\"\n" +
		"Port=0\n" +
		"\n" +
		"[Database]\n" +
		"Type=\"mysql\"\n" +
		"Address=\"localhost\"\n" +
		"Port=3306\n" +
		"SslMode=false\n" +
		"Protocol=\"tcp\"\n" +
		"Name=\"go_there_db\"\n" +
		"User=\"my_user\"\n" +
		"Password=\"superpassword\""

	tmpf, err := ioutil.TempFile(os.TempDir(), "go-there.conf")

	if err != nil {
		assert.Fail(t, err.Error())
	}

	defer func() {
		_ = os.Remove(tmpf.Name())
	}()

	if _, err := tmpf.Write([]byte(content)); err != nil {
		assert.Fail(t, err.Error())
	}

	if err := tmpf.Close(); err != nil {
		assert.Fail(t, err.Error())
	}

	tests := []struct {
		name    string
		args    args
		want    *Configuration
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				path: tmpf.Name(),
			},
			want: &Configuration{
				Server: Server{
					ListenAddress: "0.0.0.0",
					ListenPort:    8080,
				},
				Cache: Cache{
					Enabled: false,
					Type:    "",
					Address: "",
					Port:    0,
				},
				Database: Database{
					Type:     "mysql",
					Address:  "localhost",
					Port:     3306,
					SslMode:  false,
					Protocol: "tcp",
					Name:     "go_there_db",
					User:     "my_user",
					Password: "superpassword",
				},
				Endpoints: map[string]Endpoint{
					"create_users": {Enabled: true, Auth: false, AdminOnly: false},
					"manage_users": {Enabled: true, Auth: false, AdminOnly: false},
					"go":           {Enabled: true, Auth: false, AdminOnly: false},
					"manage_paths": {Enabled: true, Auth: false, AdminOnly: false},
				},
			},
			wantErr: false,
		},
		{
			name: "no_file",
			args: args{
				path: "/bad/path",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Init(tt.args.path)

			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
