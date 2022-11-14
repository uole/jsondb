package jsondb

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path"
	"sync"
)

var (
	ErrColumnNotExists     = errors.New("column not exists")
	ErrPrimaryKeyNotExists = errors.New("primary key not exists")
	ErrRecordNotFound      = errors.New("record not found")
)

type (
	Drive struct {
		ctx     context.Context
		mutex   sync.Mutex
		baseDir string
	}
)

func (drive *Drive) Open(ctx context.Context, directory string) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	drive.ctx = ctx
	drive.baseDir = directory
	if _, err = os.Stat(drive.baseDir); err != nil {
		err = os.MkdirAll(drive.baseDir, 0755)
	}
	return
}

func (drive *Drive) loadFromFile(ctx context.Context, model Model) (documents []Document, err error) {
	var (
		fp *os.File
	)
	if fp, err = os.Open(path.Join(drive.baseDir, model.TableName()+".json")); err != nil {
		return
	}
	defer func() {
		_ = fp.Close()
	}()
	documents = make([]Document, 0, 20)
	err = json.NewDecoder(fp).Decode(&documents)
	return
}

func (drive *Drive) flushToFile(ctx context.Context, tableName string, documents []Document) (err error) {
	var (
		fp *os.File
	)
	if len(documents) <= 0 {
		return
	}
	if fp, err = os.OpenFile(path.Join(drive.baseDir, tableName+".json"), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644); err != nil {
		return
	}
	err = json.NewEncoder(fp).Encode(documents)
	return
}

//Query create record query object
func (drive *Drive) Query(ctx context.Context, model Model) *Query {
	var (
		err       error
		schema    *Schema
		documents []Document
	)
	if documents, err = drive.loadFromFile(ctx, model); err != nil {
		documents = make([]Document, 0)
	}
	if schema, err = parseModelSchema(model); err != nil {
		panic(err)
	}
	return newQuery(model, schema, documents)
}

//Insert add record to database
func (drive *Drive) Insert(ctx context.Context, model Model) (err error) {
	var (
		documents []Document
	)
	drive.mutex.Lock()
	defer drive.mutex.Unlock()
	if documents, err = drive.loadFromFile(ctx, model); err != nil {
		documents = make([]Document, 0)
	}
	documents = append(documents, toDocument(model))
	err = drive.flushToFile(ctx, model.TableName(), documents)
	return
}

//Find get record by primary key
func (drive *Drive) Find(ctx context.Context, model Model) (err error) {
	var (
		n         int
		schema    *Schema
		documents []Document
		pk        interface{}
		comparePK interface{}
	)
	drive.mutex.Lock()
	defer drive.mutex.Unlock()
	if schema, err = parseModelSchema(model); err != nil {
		return
	}
	if pk, err = schema.GetPrimaryValue(model); err != nil {
		return
	}
	if documents, err = drive.loadFromFile(ctx, model); err != nil {
		return ErrRecordNotFound
	}
	for _, document := range documents {
		if comparePK = document.Getter(schema.PrimaryColumn().Name, schema.PrimaryColumn().JsonName); err == nil {
			if n, err = compare(pk, comparePK); err == nil && n == 0 {
				err = document.Decode(model)
				return
			}
		}
	}
	err = ErrRecordNotFound
	return
}

//Update update record by primary key
func (drive *Drive) Update(ctx context.Context, model Model) (err error) {
	var (
		n         int
		ok        bool
		schema    *Schema
		documents []Document
		pk        interface{}
		comparePK interface{}
	)
	drive.mutex.Lock()
	defer drive.mutex.Unlock()
	if schema, err = parseModelSchema(model); err != nil {
		return
	}
	if pk, err = schema.GetPrimaryValue(model); err != nil {
		return
	}
	if documents, err = drive.loadFromFile(ctx, model); err != nil {
		return ErrRecordNotFound
	}
	for i, v := range documents {
		if comparePK = v.Getter(schema.PrimaryColumn().Name, schema.PrimaryColumn().JsonName); err == nil {
			if n, err = compare(pk, comparePK); err == nil && n == 0 {
				documents[i] = toDocument(model)
				ok = true
				break
			}
		}
	}
	if ok {
		err = drive.flushToFile(ctx, model.TableName(), documents)
	} else {
		err = ErrRecordNotFound
	}
	return
}

//ReplaceInto add or update record by primary key
func (drive *Drive) ReplaceInto(ctx context.Context, model Model) (err error) {
	if err = drive.Update(ctx, model); err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			err = drive.Insert(ctx, model)
		}
	}
	return
}

//Delete remove record by primary key
func (drive *Drive) Delete(ctx context.Context, model Model) (err error) {
	var (
		n         int
		ok        bool
		schema    *Schema
		documents []Document
		pk        interface{}
		comparePK interface{}
	)
	drive.mutex.Lock()
	defer drive.mutex.Unlock()
	if schema, err = parseModelSchema(model); err != nil {
		return
	}
	if pk, err = schema.GetPrimaryValue(model); err != nil {
		return
	}
	if documents, err = drive.loadFromFile(ctx, model); err != nil {
		return ErrRecordNotFound
	}
	for i, v := range documents {
		if comparePK = v.Getter(schema.PrimaryColumn().Name, schema.PrimaryColumn().JsonName); comparePK != nil {
			if n, err = compare(pk, comparePK); err == nil && n == 0 {
				documents = append(documents[:i], documents[i+1:]...)
				ok = true
				break
			}
		}
	}
	if ok {
		err = drive.flushToFile(ctx, model.TableName(), documents)
	} else {
		err = ErrRecordNotFound
	}
	return
}

func (drive *Drive) Close() (err error) {

	return
}

func New() *Drive {
	d := &Drive{}
	return d
}
