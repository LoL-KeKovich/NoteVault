package mongodb

import "github.com/LoL-KeKovich/NoteVault/internal/model"

func (mc MongoClient) CreateNote(model.Note) (string, error) {
	return "", nil
}

func (mc MongoClient) GetNoteByID(id string) (model.Note, error) {
	return model.Note{}, nil
}

func (mc MongoClient) GetNotes() ([]model.Note, error) {
	return []model.Note{}, nil
}

func (mc MongoClient) UpdateNote(id, name, text, color, media string, order int) (int, error) {
	return 0, nil
}

func (mc MongoClient) DeleteNote(id string) (int, error) {
	return 0, nil
}
