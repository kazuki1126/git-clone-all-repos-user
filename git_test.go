package main

import (
	"os"
	"testing"
)

func TestCreateRepoInLocal(t *testing.T) {
	type args struct {
		url      string
		repoName string
	}
	tt := []struct {
		name        string
		args        args
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Fail pattern",
			args: args{
				url:      "",
				repoName: "",
			},
			wantErr:     true,
			expectedErr: unexpectedErr,
		},
		{
			name: "Fail pattern2",
			args: args{
				url:      "https://api.github.com/repos/kazuki1126/black-jack-game/contents/",
				repoName: "",
			},
			wantErr:     true,
			expectedErr: unexpectedErr,
		},
		{
			name: "Success",
			args: args{
				url:      "https://api.github.com/repos/kazuki1126/black-jack-game/contents/",
				repoName: "black-jack-game",
			},
			wantErr: false,
		},
	}
	os.Mkdir("black-jack-game", dirPerm)
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := createRepoInLocal(tc.args.url, tc.args.repoName)
			if err != nil && err != tc.expectedErr {
				t.Errorf("Expected error: %v, Got: %v", tc.expectedErr, err)
			} else if (err != nil) != tc.wantErr {
				t.Errorf("Error occured %v WantErr: %v", err, tc.wantErr)
			}
		})
	}
	os.RemoveAll("black-jack-game")
}

func TestProcessFile(t *testing.T) {
	type args struct {
		repoName    string
		filepath    string
		downloadURL string
	}
	tt := []struct {
		name        string
		args        args
		wantErr     bool
		expectedErr error
	}{
		{
			name: "fail pattern",
			args: args{
				repoName:    "",
				filepath:    "main.go",
				downloadURL: "https://raw.githubusercontent.com/kazuki1126/black-jack-game/master/main.go",
			},
			wantErr:     true,
			expectedErr: unexpectedErr,
		},
		{
			name: "Success",
			args: args{
				repoName:    "black-jack-game",
				filepath:    "main.go",
				downloadURL: "https://raw.githubusercontent.com/kazuki1126/black-jack-game/master/main.go",
			},
			wantErr: false,
		},
	}
	os.Mkdir("black-jack-game", dirPerm)
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := processFile(tc.args.repoName, tc.args.filepath, tc.args.downloadURL)
			if err != nil && err != tc.expectedErr {
				t.Errorf("Expected error: %v, Got: %v", tc.expectedErr, err)
			} else if (err != nil) != tc.wantErr {
				t.Errorf("Error occurred: %v WantErr: %v", err, tc.wantErr)
			}
		})
	}
	os.RemoveAll("black-jack-game")
}

func TestGetAllRepos(t *testing.T) {
	type args struct {
		userName string
	}

	tt := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Fail pattern",
			args: args{
				userName: "",
			},
			wantErr: true,
		},
		{
			name: "Fail pattern2",
			args: args{
				// username on github cannot start with hyphen
				userName: "-notexitsting",
			},
			wantErr: true,
		},
		{
			name: "Success",
			args: args{
				userName: "kazuki1126",
			},
			wantErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := getAllRepos(tc.args.userName); (err != nil) != tc.wantErr {
				t.Errorf("Error occurred: %v WantErr: %v", err, tc.wantErr)
			}
		})
	}
}
