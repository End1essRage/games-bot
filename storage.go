package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type Game struct {
	Id    int
	Owner string
	Name  string
}

type Storage struct {
	basePath string
}

func NewStorage(basePath string) *Storage {
	return &Storage{basePath: basePath}
}

const defaultPerm = 0774 //read and write

func (s *Storage) Add(chatUid string, game Game) error {
	fPath := filepath.Join(s.basePath, chatUid+".json")

	if err := os.MkdirAll(s.basePath, defaultPerm); err != nil {
		return fmt.Errorf("cant save: %w", err)
	}

	file, err := os.Create(fPath)
	if err != nil {
		return fmt.Errorf("cant create file: %w", err)
	}

	defer func() { _ = file.Close() }()

	newId := 0
	data, err := os.ReadFile(fPath)
	if err != nil {
		return err
	}

	games := make([]Game, 0)
	if err := json.Unmarshal(data, games); err != nil {
		logrus.Info("no data in file")
	}

	if len(games) > 0 {
		newId = games[len(games)-1].Id + 1
	}

	game.Id = newId
	games = append(games, game)

	toWrite, err := json.Marshal(games)
	if err != nil {
		return err
	}

	_, err = file.Write(toWrite)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Get(chatUid string) ([]Game, error) {
	fPath := filepath.Join(s.basePath, chatUid+".json")

	data, err := os.ReadFile(fPath)
	if err != nil {
		return nil, err
	}

	games := make([]Game, 0)
	if err := json.Unmarshal(data, games); err != nil {
		logrus.Info("no data in file")
	}

	return games, nil
}
