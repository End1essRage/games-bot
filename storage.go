package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type Data struct {
	Games []Game
}

type Game struct {
	Id    int
	Owner string
	Name  string
}

type Storage struct {
	basePath string
}

const defaultPerm = 0774 //read and write

func NewStorage(basePath string) *Storage {
	//Создаем папку
	if err := os.MkdirAll(basePath, defaultPerm); err != nil {
		logrus.Fatal(err)
	}

	return &Storage{basePath: basePath}
}

func (s *Storage) Add(chatUid string, game Game) error {
	//путь до файла
	fPath := filepath.Join(s.basePath, chatUid+".json")

	data := s.getData(fPath)

	//задаем id новой записи
	newId := 1

	if len(data.Games) > 0 {
		newId = data.Games[len(data.Games)-1].Id + 1
	}

	game.Id = newId
	//Добавляем новую запись
	data.Games = append(data.Games, game)

	err := s.writeData(data, fPath)
	if err != nil {
		logrus.Error(err)
	}

	return nil
}

func (s *Storage) Get(chatUid string) []Game {
	fPath := filepath.Join(s.basePath, chatUid+".json")

	data := s.getData(fPath)

	return data.Games
}

func (s *Storage) getData(path string) Data {
	file := s.createOpenFile(path)

	defer func() { _ = file.Close() }()

	//считываем данные с файла
	buffer, err := os.ReadFile(path)
	if err != nil {
		logrus.Error(err)
	}

	//десериализуем данные
	var data Data

	if err := json.Unmarshal(buffer, &data); err != nil {
		logrus.Info("no data in file")
	}

	return data
}

func (s *Storage) writeData(data Data, path string) error {
	file := s.createOpenFile(path)

	defer func() { _ = file.Close() }()
	//Сериализуем
	toWrite, err := json.Marshal(data)
	if err != nil {
		return err
	}

	//Записываем

	_, err = file.Write(toWrite)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) createOpenFile(path string) *os.File {
	file, err := os.OpenFile(path, os.O_RDWR, defaultPerm)
	if err != nil {
		f, err := os.Create(path)
		if err != nil {
			logrus.Fatal(err)
		}

		return f
	}

	return file
}
